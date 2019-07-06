package gossip_protocol

import (
	//"encoding/hex"
	"encoding/json"
	"fmt"
	//"golang.org/x/crypto/openpgp/packet"
	"log"
	//"golang.org/x/crypto/sha3"
	//"math/rand"
	//"os"
	//"strconv"
	//"strings"
	//"time"
	//s "../../p5security"
	s "../identity"
	sbc_p "../sync_blockchain"
	//"../../p2/block"
)

//HeartBeatData struct defines the data to be sent between peers perodically
type HeartBeatData struct {
	IfNewBlock       bool             `json:"ifNewBlock"`
	Id               int32            `json:"id"`
	Pid              s.PublicIdentity `json:"pid"`
	BlockJson        string           `json:"blockJson"`
	SignForBlockJson []byte           `json:"signForBlockJson"`
	PeerMapJson      string           `json:"peerMapJson"`
	PeerMapPidJson   string           `json:"peerMapPidJson"`
	Addr             string           `json:"addr"`
	Hops             int32            `json:"hops"`
}

//NewHeartBeatData creates new HeartBeatData
func NewHeartBeatData(ifNewBlock bool, id int32, pid s.PublicIdentity, blockJson string, signForBlockJson []byte, peerMapJson string, peerMapPidJson string, addr string) HeartBeatData {

	return HeartBeatData{
		IfNewBlock:       ifNewBlock,
		Id:               id,
		Pid:              pid,
		BlockJson:        blockJson,
		SignForBlockJson: signForBlockJson,
		PeerMapJson:      peerMapJson,
		PeerMapPidJson:   peerMapPidJson,
		Addr:             addr,
		Hops:             2, //todo change to 3
	}
}

//PrepareHeartBeatData func prepares  and returns heartbeat
func PrepareHeartBeatData(sbc *sbc_p.SyncBlockChain, selfId int32, selfpid s.PublicIdentity, peerMapBase64 string, peerMapPidBase64 string, addr string, makingNew bool, newBlockJson string, signForBlockJson []byte) HeartBeatData {

	//makeNew := rand.Int() % 2
	//makeNew := 1
	var ifNewBlock bool
	blockJSON := "{}"

	if makingNew == true {
		ifNewBlock = true
		blockJSON = newBlockJson //newBlockJson(sbc) // creating a new block's json

	} else {
		ifNewBlock = false
	}

	return NewHeartBeatData(
		ifNewBlock,
		selfId,
		selfpid,
		blockJSON,
		signForBlockJson,
		peerMapBase64,
		peerMapPidBase64,
		addr,
	)
}

////newBlockJson func create a new block, inserts it into blockchain and returns json string
//func newBlockJson(sbc *SyncBlockChain) string {
//	mpt := p1.MerklePatriciaTrie{}
//	mpt.Initial()
//
//	for i := 0; i <= 5; i++ {
//		random := strconv.Itoa(rand.Int() % 10)
//		mpt.Insert("Time Now "+random, strconv.Itoa(rand.Int())+"is humbly yours "+os.Args[1])
//	}
//	b1 := sbc.GenBlock(mpt, "0")
//	sbc.Insert(b1)
//	return block.EncodeToJSON(&b1)
//}

//EncodeToJson func encodes HeartBeatData to json byte array
func (data *HeartBeatData) EncodeToJsonByteArray() []byte {

	dataEncodedBytearray, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Here in err condition of EncodeToJsonByteArray of heartbeat.go")
		return nil
	}
	return dataEncodedBytearray
}

//EncodeToJson func encodes HeartBeatData to json string
func (data *HeartBeatData) EncodeToJson() string {

	dataEncodedBytearray, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Here in err condition of EncodeToJson of heartbeat.go")
		return "{}" //empty heartbeat
	}
	return string(dataEncodedBytearray)
}

//DecodeToHeartBeatData func decodes json string to HeartBeatData
func DecodeToHeartBeatData(heartBeatDatajson string) HeartBeatData {
	hbd := HeartBeatData{}
	err := json.Unmarshal([]byte(heartBeatDatajson), &hbd)
	if err != nil {
		log.Println("Err in DecodeToHeartBeatData in heartbeat.go - err : ", err)
		log.Println("Error heartBeatDatajson : ", heartBeatDatajson)
	} else {
		//log.Println("DecodeToHeartBeatData in heartbeat.go - SUCCESSFUL : heartBeatDatajson : ", heartBeatDatajson)
	}
	return hbd
}

func TestHeartBeat() {

}
