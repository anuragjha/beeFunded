package gossip_protocol

import (
	"bytes"
	"container/ring"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"

	//"strconv"
	//"../../p5"
	"sync"
	//s "../../p5security"
	s "../identity"
)

//PeerList contains selfId, peerMap, max length, and a mutex
type PeerList struct {
	selfId     int32
	secureId   s.Identity //p5
	peerMap    map[string]int32
	peerMapPid map[string]s.PublicIdentity
	maxLength  int32
	mux        sync.Mutex
}

///////////
// Pair - data structure to hold a key/value pair - addr/id.
type Pair struct {
	addr string
	id   int32
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].id < p[j].id }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]int32) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		//fmt.Println("in sortMapByValue : k, v :", k, v)
		p[i] = Pair{
			addr: k,
			id:   int32(v),
		}
		//fmt.Println("in sortMapByValue : p[i] :", p[i])
		i++
	}
	//fmt.Println("in sortMapByValue : p :", p)
	sort.Sort(p)
	//fmt.Println("in sortMapByValue : sorted p :", p)
	return p
}

///////////

//NewPeerList func creates a New PeerList for a id and maxLength
func NewPeerList(id int32, sid s.Identity, maxLength int32) PeerList {

	return PeerList{
		selfId:     id,
		secureId:   sid,
		peerMap:    make(map[string]int32),
		peerMapPid: make(map[string]s.PublicIdentity),
		maxLength:  maxLength,
	}
}

// ONLY FOR TEST PURPOSES
func TestNewPeerList(id int32 /*sid s.Identity,*/, maxLength int32) PeerList {

	return PeerList{
		selfId: id,
		/*secureId:   sid,*/
		peerMap:    make(map[string]int32),
		peerMapPid: make(map[string]s.PublicIdentity),
		maxLength:  maxLength,
	}
}

//Add func adds a peer with addr and id to peerMap
func (peers *PeerList) Add(addr string, id int32) {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	peers.peerMap[addr] = id
}

//Add func adds a peer with addr and id to peerMap
func (peers *PeerList) AddPid(addr string, id s.PublicIdentity) {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	peers.peerMapPid[addr] = id
}

//Delete func deletes a peer with specific addr
func (peers *PeerList) Delete(addr string) {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	delete(peers.peerMap, addr)
}

////Delete func deletes a peer with specific addr
//func (peers *PeerList) DeletePid(addr string) {
//	peers.mux.Lock()
//	defer peers.mux.Unlock()
//
//	delete(peers.peerMapPid, addr)
//}

//Rebalance func changes the PeerMap to contain take maxLength(32) closest peers (by Id)
func (peers *PeerList) Rebalance() {

	if int32(len(peers.peerMap)) > peers.maxLength {

		peers.mux.Lock()
		defer peers.mux.Unlock()

		//fmt.Println("in Rebalance")
		//fmt.Println("in Rebalance : original peerMap length : ", len(peers.peerMap))
		peers.peerMap["selfAddr"] = peers.selfId //adding self id to peerMap
		sortedAddrIDList := sortMapByValue(peers.peerMap)
		//fmt.Println("in Rebalance : sortedAddrIDList : ", sortedAddrIDList)
		sortedAddrIDListLength := len(sortedAddrIDList)
		//fmt.Println("in Rebalance : sortedAddrIDListLength : ", sortedAddrIDListLength)

		peers.peerMap = peers.getBalancedPeerMap(sortedAddrIDListLength, sortedAddrIDList)

	}
}

func (peers *PeerList) getBalancedPeerMap(sortedAddrIDListLength int, sortedAddrIDList PairList) map[string]int32 {
	r := ring.New(sortedAddrIDListLength) // new ring
	useRingPtr := r

	//initialize ring with sortedAddrIDList values
	for i := 0; i < sortedAddrIDListLength; i++ {
		r.Value = sortedAddrIDList[i]
		//fmt.Println("in Rebalance : r.Value : ", r.Value)
		if sortedAddrIDList[i].id == peers.selfId {
			useRingPtr = r
			//fmt.Println("in Rebalance : useRingPtr : ", useRingPtr)
		}
		r = r.Next()
	}
	newPeerMap := make(map[string]int32)
	r = useRingPtr
	//fmt.Println("in Rebalance : useRingPtr : ", useRingPtr)
	for i := 1; i <= int(peers.maxLength/2); i++ {
		r = r.Prev()
		pair := r.Value.(Pair)
		newPeerMap[pair.addr] = pair.id
	}
	r = useRingPtr
	for i := 1; i <= int(peers.maxLength/2); i++ {
		r = r.Next()
		pair := r.Value.(Pair)
		newPeerMap[pair.addr] = pair.id
	}

	return newPeerMap
}

//Show func returns PeerMap string
func (peers *PeerList) Show() string {
	var buffer bytes.Buffer

	buffer.WriteString("This is PeerMap:\n")
	for k := range peers.peerMap {
		buffer.WriteString("Addr:" + k + " Id:" + strconv.Itoa(int(peers.peerMap[k])) + "\n")
	}
	return buffer.String()
}

//Show func returns PeerMap string
func (peers *PeerList) ShowPids() string {
	var buffer bytes.Buffer

	buffer.WriteString("This is PeerMapPid:\n")
	for k := range peers.peerMapPid {
		buffer.WriteString("Addr:" + k + " Pid PubK : " + peers.peerMapPid[k].PublicKey.N.String() + "\n" +
			"Pid Label : " + peers.peerMapPid[k].Label + "\n")
	}
	return buffer.String()
}

////Register func assigns a value to selfId
//func (peers *PeerList) Register(id int32) {
//	peers.selfId = id
//	fmt.Printf("SelfId=%v\n", id)
//}

//Copy func returns a copy of the peerMap
func (peers *PeerList) Copy() map[string]int32 {

	peers.mux.Lock()
	defer peers.mux.Unlock()

	copyOfPeerMap := make(map[string]int32)
	for k := range peers.peerMap {
		copyOfPeerMap[k] = peers.peerMap[k]
	}

	return copyOfPeerMap
}

//Copy func returns a copy of the peerMap
func (peers *PeerList) CopyPids() map[string]s.PublicIdentity {

	peers.mux.Lock()
	defer peers.mux.Unlock()

	copyOfPeerMapPid := make(map[string]s.PublicIdentity)
	for k := range peers.peerMapPid {
		copyOfPeerMapPid[k] = peers.peerMapPid[k]
	}

	return copyOfPeerMapPid
}

//GetSelfId func returns selfId of Peer
func (peers *PeerList) GetSelfId() int32 {
	return peers.selfId
}

//PeerMapToJson func returns a json string of PeerMap or an error
func (peers *PeerList) PeerMapToJson() (string, error) {
	peers.mux.Lock()

	jsonOfPeerMap, err := json.Marshal(peers.peerMap)

	peers.mux.Unlock()

	return string(jsonOfPeerMap), err
}

//PeerMapToJson func returns a json string of PeerMap or an error
func PeerMapToJson(peermap map[string]int32) (string, error) {

	jsonOfPeerMap, err := json.Marshal(peermap)

	return string(jsonOfPeerMap), err
}

//PeerMapToJson func returns a json string of PeerMap or an error
func (peers *PeerList) PeerMapPidToJson() (string, error) {
	peers.mux.Lock()

	jsonOfPeerMapPid, err := json.Marshal(peers.peerMapPid)

	peers.mux.Unlock()

	return string(jsonOfPeerMapPid), err
}

//PeerMapSIDToJson func returns a json string of PeerMap or an error
func PeerMapPidToJson(peerMapPid map[string]s.PublicIdentity) (string, error) {

	jsonOfPeerMapPid, err := json.Marshal(peerMapPid)

	return string(jsonOfPeerMapPid), err
}

//InjectPeerMapJson func injects the new PeerMap into existing PeerMap, except for the entry corresponding to self
func (peers *PeerList) InjectPeerMapJson(peerMapJsonStr string, selfAddr string) {

	var newPeerMap map[string]int32
	err := json.Unmarshal([]byte(peerMapJsonStr), &newPeerMap)
	if err == nil {
		peers.mux.Lock()

		for k := range newPeerMap {
			if /*_, ok := peers.peerMap[k]; !ok &&*/ k != selfAddr {
				peers.peerMap[k] = newPeerMap[k]
			}
		}

		peers.mux.Unlock()
	}
}

//InjectPeerMapJson func injects the new PeerMap into existing PeerMap, except for the entry corresponding to self
func (peers *PeerList) InjectPeerMapPidJson(peerMapPidJsonStr string, selfAddr string) {

	var recvPeerMapPid map[string]s.PublicIdentity
	err := json.Unmarshal([]byte(peerMapPidJsonStr), &recvPeerMapPid)
	if err == nil {
		//peerMapPidCopy := peers.CopyPids()
		for addr, pid := range recvPeerMapPid {
			if _, ok := peers.peerMapPid[addr]; !ok {
				peers.mux.Lock()
				peers.peerMapPid[addr] = pid
				peers.mux.Unlock()
			}
		}
	} else {
		log.Println("Error in Inject PeerMapPidJson : err : ", err)
	}
}

func TestPeerListRebalance() {
	peers := NewPeerList(5, s.Identity{}, 4)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected := NewPeerList(5, s.Identity{}, 4)
	expected.Add("1111", 1)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	expected.Add("-1-1", -1)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, s.Identity{}, 2)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected = NewPeerList(5, s.Identity{}, 2)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, s.Identity{}, 4)
	peers.Add("1111", 1)
	peers.Add("7777", 7)
	peers.Add("9999", 9)
	peers.Add("11111111", 11)
	peers.Add("2020", 20)
	peers.Rebalance()
	expected = NewPeerList(5, s.Identity{}, 4)
	expected.Add("1111", 1)
	expected.Add("7777", 7)
	expected.Add("9999", 9)
	expected.Add("2020", 20)
	fmt.Println(reflect.DeepEqual(peers, expected))
}
