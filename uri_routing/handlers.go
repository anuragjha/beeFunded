package uri_routing

import (
	//"bytes"
	//"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	//"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	//"../p1"
	mptp "../data_structure/mpt"
	//b "../p2/block"
	b "../block"
	//"../p2"
	//"../p4"
	"../p5"
	"../pow"
	//s "../p5security"
	s "../identity"
	//"./data"
	gp "../gossip_protocol"
	sbc "../sync_blockchain"
)

//var TA_SERVER = "http://localhost:6688"
var INIT_SERVER = "http://localhost:6686"

//var REGISTER_SERVER = TA_SERVER + "/peer"
//var REGISTER_SERVER = INIT_SERVER + "/peer"

//SELF_ADDR var BC_DOWNLOAD_SERVER = TA_SERVER + "/upload"
var BC_DOWNLOAD_SERVER = INIT_SERVER + "/upload"
var PID_DOWNLOAD_SERVER = INIT_SERVER + "/uploadpids"

//changes in init for arg of port provided
var SELF_ADDR = "http://localhost:6686"
var SELF_ADDR_PREFIX = "http://localhost:"

// SBC is safe for distributed use
var SBC sbc.SyncBlockChain

//Peers is the Peer List which is for each node
var Peers gp.PeerList

var tryingForHeight int32
var GetNewParent bool

const Difficulty = 5

var ifStarted bool

var ID s.Identity

//var BalanceBook p5.BalanceBook
//var Wallet p5.Wallet

var TxPool p5.TransactionPool

func init() {
	// This function will be executed before everything else.

	//init coz node not removed from peerlist and receieve heartbeat even before it start()s
	//id := Register()
	//Peers = gp.NewPeerList(id, SID,32)

	SELF_ADDR = SELF_ADDR_PREFIX + os.Args[1]
	fmt.Println("Node : ", SELF_ADDR)

	//init BalanceBook

	//init wallet

}

// Start handler - does Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {

	if ifStarted == false {
		ifStarted = true

		id := Register() //register ID\

		ID = s.NewIdentity(strconv.Itoa(int(id)))

		Peers = gp.NewPeerList(id, ID, 32) //initialize PeerList // 32 sunnit

		Peers.AddPid(SELF_ADDR, ID.GetMyPublicIdentity()) //p5

		SBC = sbc.NewBlockChain() //create new Block chain //apr4

		//currency //p5 //todo p5 todo p5 CURRENCY
		TxPool = p5.NewTransactionPool()
		if SELF_ADDR != INIT_SERVER {
			log.Println("Asking for Transaction Pool !!!")
			GetTransactionPool() //initialize txPool with existing txs in Pool
		}

		//InitBalanceBook() // in handlerCurrencyHelper
		//InitWallet()      // in handlerCurrencyHelper

		if strings.Compare(SELF_ADDR, INIT_SERVER) == 0 {
			fmt.Println("Generating Genesis block")
			mpt := mptp.MerklePatriciaTrie{}
			mpt.Insert("beeFunded", "{}")
			nonce := pow.FindNonce("genesis", &mpt, Difficulty)
			b1 := SBC.GenBlock(1, "genesis", mpt, nonce, ID.GetMyPublicIdentity())
			SBC.Insert(b1)

			//GiveGenesisTokens(ID)

		}
		GiveMinerTokens(ID)

		//if Peers.GetSelfId() != 6686 { //download if not 6686
		if SELF_ADDR != INIT_SERVER {

			Peers.Add(INIT_SERVER, int32(6686)) // add Init server to peer list of node
			Download()                          //download BlockChain
			DownloadPeerMapPid()
		}

		//start HearBeat
		go StartHeartBeat()
		go StartTryingNonces() //pow //p5 - also creates mpt

	}
	w.WriteHeader(200)
	_, err := w.Write([]byte("started : " + SELF_ADDR + "\n Pid : \n" + Peers.ShowPids()))
	if err != nil {
		log.Println("Err -  in start - during writing to client")
	}

}

//Show func -  Display peerList and sbc
func Show(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "%s\n%s\n%s", Peers.Show(), Peers.ShowPids(), SBC.Show())
	if err != nil {
		log.Println("Err in show func while writing response")
	}
}

// Upload blockchain to whoever called this method, return jsonStr
func Upload(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		//data.PrintError(err, "Upload") // todo
		log.Println("Err - in Upload func")
	}
	fmt.Fprint(w, blockChainJson)

	//remove comments above after testing
	//UploadGenesis(w, r)
}

// Upload blockchain to whoever called this method, return jsonStr
func UploadPids(w http.ResponseWriter, r *http.Request) {

	copyPids := Peers.CopyPids()
	peerMapPidJson, err := gp.PeerMapPidToJson(copyPids) //SBC.BlockChainToJson()
	if err != nil {
		//data.PrintError(err, "Upload") // todo
		log.Println("Err - in Upload func")
	}
	fmt.Fprint(w, peerMapPidJson)

}

// UploadBlock func - Upload a block to whoever called this method, return jsonStr
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ubHeight, err := strconv.Atoi(vars["height"])
	if err != nil {
		returnCode500(w, r)
	} else {
		ubHash := vars["hash"]
		//fmt.Println("\nuploading block of -\nubHeight : ", ubHeight)
		//fmt.Println("ubHash : ", ubHash, "\n\n")

		uBlock, found := SBC.GetBlock(int32(ubHeight), ubHash)
		if found == false {
			fmt.Println("Err : in Handlers - UploadBlock - found = false - 204")
			returnCode204(w, r)
		} else {
			fmt.Println("in Handlers - UploadBlock - found = true")
			blockJson := b.EncodeToJSON(&uBlock)
			_, err = fmt.Fprint(w, blockJson)
			if err != nil {
				log.Println("Err : in Handlers - UploadBlock - during writing response")
			}
		}
	}

}

// HeartBeatReceive func - Received a heartbeat in request body
func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Err : in HeartBeatReceive - reached err of ioutil.ReadAll -")
		log.Println(err)
	} else {
		defer r.Body.Close()

		//decrypt with privateKey //p5
		//deciphered := p5.DecryptMessageWithPrivateKey(SID.GetMyPrivateKey(), body)
		heartBeat := gp.DecodeToHeartBeatData(string( /*deciphered*/ body)) // heartBeat struct

		if isSignVerified(heartBeat) { //p5
			fmt.Println("Sign verified :-)")

			go processHeartBeat(gp.DecodeToHeartBeatData(string(body))) // process for the receieved heartbeat

			go forwardHeartBeat(gp.DecodeToHeartBeatData(string(body))) // forward the heartBeat // here sunnit

		} else {
			fmt.Println("Sign NOT verified :-(")
			fmt.Println("sig recv : ", heartBeat.SignForBlockJson)
		}

	}

}

//Canonical func -  Display canonical chain
func Canonical(w http.ResponseWriter, r *http.Request) {

	canonicalChains := pow.GetCanonicalChains(&SBC)

	_, _ = fmt.Fprint(w, "Canonical Chain(s) : \n")
	for i, chain := range canonicalChains {
		_, _ = fmt.Fprint(w, "\nChain #"+strconv.Itoa(i+1))
		_, err := fmt.Fprint(w, "\n", chain.ShowCanonical())
		if err != nil {
			_, _ = fmt.Fprint(w, "ERROr in Canonical")
		}
	}
	//
}

// Register to TA's server, get an ID
func Register() int32 {

	body := os.Args[1]

	id, err := strconv.Atoi(string(body))
	if err != nil {
		log.Fatal(err)
		return 0
	}

	return int32(id)
}

// Download blockchain from TA server
func Download() {
	resp, err := http.Get(BC_DOWNLOAD_SERVER)
	//resp, err := http.Get("http://localhost:6686/upload/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) //blockChainJson
	if err != nil {
		log.Fatal(err)
	}

	SBC.UpdateEntireBlockChain(string(body))
}

func DownloadPeerMapPid() {
	resp, err := http.Get(PID_DOWNLOAD_SERVER)
	//resp, err := http.Get("http://localhost:6686/uploadpids/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) //blockChainJson
	if err != nil {
		log.Fatal(err)
	}

	//SBC.UpdateEntireBlockChain(string(body))
	Peers.InjectPeerMapPidJson(string(body), SELF_ADDR)
}

// Upload genesis blockchain to whoever called this method, return jsonStr
func UploadGenesis(w http.ResponseWriter, r *http.Request) {

	nbc := sbc.NewBlockChain()
	gbl, _ := SBC.Get(1)
	nbc.Insert(gbl[0])

	blockChainJson, err := nbc.BlockChainToJson()
	if err != nil {
		//data.PrintError(err, "Upload") // todo
		log.Println("in Err of Upload Genesis")
	}
	_, err = fmt.Fprint(w, blockChainJson)
	if err != nil {
		log.Println("in Err of Upload Genesis writing response")
	}
}

func isSignVerified(heartBeat gp.HeartBeatData) bool {

	senderPid := s.PublicIdentity(heartBeat.Pid)
	//senderHashForKey := s.GenerateHashForKey(senderPid.Label)
	isMatch := s.VerifySingature(senderPid.PublicKey, []byte(heartBeat.BlockJson), []byte(heartBeat.SignForBlockJson))
	return isMatch
}

// processHeartBeat func updates the peerlist, and IfNewBlock then insert the block in SBC
func processHeartBeat(heartBeat gp.HeartBeatData) {

	//use hearBeatData to update peer list
	updatePeerList(&heartBeat) //updates PeerMap and PeerMapPid //p5

	//and get block if the IfNewBlock is set to true
	if heartBeat.IfNewBlock { //add block in blockchain

		newBlock := b.DecodeFromJSON(heartBeat.BlockJson)
		mptHash := mptp.MerklePatriciaTrie(newBlock.Value).Root

		//y := sha3.Sum256([]byte(newBlock.Header.ParentHash + newBlock.Header.Nonce + mptHash))
		//y1 := hex.EncodeToString(y[:])
		//log.Println("++++++++++++++++++++ receieving MPT ROOT at height : ", mptHash)
		//log.Println("++++++++++++++++++++ receving proof : ", y1)

		if pow.POW(newBlock.Header.ParentHash, newBlock.Header.Nonce, mptHash, Difficulty) {
			//fmt.Println(":::: PROCESSING HeartBeat : in ProcessHeartBeat : newBlock.Value.Root : ", mptHash)

			//apr4
			//hold parent / grandparent / etc blocks to be put once we find the begining block based on nodes local copy
			//var blockHolder []block.Block
			//// apr4

			if SBC.CheckParentHash(newBlock) {

				SBC.Insert(newBlock) // if parentHash exist then directly insert and POW is satisfied

			} else if AskForBlock(newBlock.Header.Height-1, newBlock.Header.ParentHash, make([]b.Block, 0) /*, SBC.GetLength(), newBlock.Header.Height-1*/) {
				//if parent cannot be found then ask for parent blocks and insert all parent then insert newBlock
				SBC.Insert(newBlock)
				//AskForBlock(newBlock.Header.Height, newBlock.Header.ParentHash, make([]block.Block, 0), SBC.GetLength()+1, newBlock.Header.Height+1)

			}
		}

		//fmt.Println("NOT processing HeartBeat : in ProcessHeartBeat : newBlock.Value.Root : ", newBlock.Value.Root)

	}
}

// ForwardHeartBeat func forwards the receieved heartbeat to all its peers
func forwardHeartBeat(heartBeatData gp.HeartBeatData) {

	Peers.Rebalance()
	peerMap := Peers.Copy()
	hopCount := heartBeatData.Hops //to forward heartbeat
	if hopCount > 0 {
		heartBeatData.Hops--
		heartBeatData.Id = Peers.GetSelfId()
		heartBeatData.Pid = ID.GetMyPublicIdentity()
		heartBeatData.SignForBlockJson = ID.GenSignature([]byte(heartBeatData.BlockJson))
		heartBeatData.Addr = SELF_ADDR

		//json, _ := json.Marshal(peerMap)
		//heartBeatData.PeerMapJson = string(json)

		//list over peers and send them heartBeat
		if len(peerMap) > 0 {
			for peer := range peerMap {
				_, _ = http.Post(peer+"/heartbeat/receive", "application/json; charset=UTF-8",
					strings.NewReader(heartBeatData.EncodeToJson()))
			}
		}
	}

}

// updatePeerList func updates the existing peerlist with data from received peerMap
func updatePeerList(heartBeat *gp.HeartBeatData) {
	Peers.Add(heartBeat.Addr, heartBeat.Id)
	Peers.InjectPeerMapJson(heartBeat.PeerMapJson, SELF_ADDR)

	Peers.AddPid(heartBeat.Addr, heartBeat.Pid) //p5
	Peers.InjectPeerMapPidJson(heartBeat.PeerMapPidJson, SELF_ADDR)
}

// AskForBlock - Ask another server to return a block of certain height and hash
func AskForBlock(height int32, hash string, blockHolder []b.Block) bool {

	//found := false
	Peers.Rebalance()
	peerMap := Peers.Copy()
	//var peersToRemove []string

	//list over peers and send them heartBeat
	//if len(peerMap) > 0 {
	for peer := range peerMap {
		//fmt.Println("\n\nin AskForBlock : req URL : ", peer+"/block/"+strconv.Itoa(int(height))+"/"+hash)
		resp, err := http.Get(peer + "/block/" + strconv.Itoa(int(height)) + "/" + hash)
		if err != nil {
			log.Println("Askblock Err 1 : ", err)
			log.Println("in AskForBlock - deleting peer : ", peer)
			Peers.Delete(peer)
			continue

		} else {
			defer resp.Body.Close() //moved from above err check to here

			body, err := ioutil.ReadAll(resp.Body) //blockJson
			if err != nil {
				log.Println("Askblock Err 2 : ", err)
				continue
			}

			reqBlock := b.DecodeFromJSON(string(body))
			//fmt.Println("\n in AskForBlock : reqBlock", reqBlock, "\n")

			if SBC.CheckParentHash(reqBlock) {
				SBC.Insert(reqBlock)                         // this block
				for i := len(blockHolder) - 1; i >= 0; i-- { // and rest of previous block
					SBC.Insert(blockHolder[i])
				}
				return true
			}

			//if !SBC.CheckParentHash(reqBlock) {  // if parenthash not in local blockchain
			fmt.Println("AskBlock - cannot find parent block for block in height : ", height)
			blockHolder = append(blockHolder, reqBlock) //apr4
			if height <= 1 {
				//continue //apr5
				break
			}
			AskForBlock(height-1, reqBlock.Header.ParentHash, blockHolder) // ask for parents parent block

		} // parsing responsed block
	} // looping peerlist

	return false // if parent block not found

}

//StartHeartBeat func periodically sends heartbeatdata to peers
func StartHeartBeat() {

	for true {
		Peers.Rebalance()
		peerMap := Peers.Copy()
		//PeerMapJson, _ := Peers.PeerMapToJson() //
		PeerMapJson, _ := gp.PeerMapToJson(peerMap) //apr4

		peerMapPid := Peers.CopyPids()
		PeerMapPidJson, err := gp.PeerMapPidToJson(peerMapPid) //p5
		if err != nil {
			fmt.Println("in StartHeartBeat, Error in Converting PeerMapPid to json : err - ", err)
		}

		//selfAddr := "http://localhost:" + os.Args[1] // SELF_ADDR

		//func PrepareHeartBeatData(sbc *SyncBlockChain, selfId int32, selfSid p5.PublicIdentity, peerMapBase64 string, peerMapSIDBase64 string, addr string, makingNew bool, newBlockJson string, signForBlockJson string) HeartBeatData
		heartBeat := gp.PrepareHeartBeatData(&SBC, Peers.GetSelfId(), ID.GetMyPublicIdentity(), PeerMapJson, PeerMapPidJson, SELF_ADDR, false, "{}", ID.GenSignature([]byte("{}")))

		//list over peers and send them heartBeat
		if len(peerMap) > 0 {
			for peerAddr := range peerMap {

				_, err := http.Post(peerAddr+"/heartbeat/receive", "application/json; charset=UTF-8",
					strings.NewReader(heartBeat.EncodeToJson())) //apr4

				//encrypting heartbeat json with peerpublic id // todo en-crypto
				//EncryptMessageWithPublicKey(hash1 hash.Hash, publicKey *rsa.PublicKey, message string, label string)
				//peerPid := p5.PublicIdentity(peerMapPid[peerAddr]).PublicKey
				//ciphered := p5.EncryptMessageWithPublicKey( peerPid, heartBeat.EncodeToJson() )
				//_, err := http.Post(peerAddr+"/heartbeat/receive", "application/octet-stream; charset=UTF-8",
				//	strings.NewReader(string(ciphered))) //apr30 p5
				if err != nil {
					Peers.Delete(peerAddr)
					//fmt.Println("deleting peer : ", peer)
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

////          pow ///

//StartTryingNonces func sends heartbeatdata with new block information to peers
func StartTryingNonces() {

	tryingNonces( /*parentHash ,&mpt, */ Difficulty)

}

// tryingNonces func tries to create a new block
func tryingNonces( /*parentHash string, mpt *mpt.MerklePatriciaTrie, */ difficulty int) {

	// y = SHA3(parentHash + nonce + mptRootHash)

	var parentBlock b.Block
	var parentHash string
	var mpt mptp.MerklePatriciaTrie

	GetNewParent = true
	var nonce string

	for {

		if GetNewParent == true {
			GetNewParent = false

			parentBlock = SBC.GetLatestBlocks()[0] //[rand.Int()%len(SBC.GetLatestBlocks())]//random parent from blocks at latest height
			parentHash = parentBlock.Header.Hash
			tryingForHeight = parentBlock.Header.Height + 1

			//mpt = mpt.GenerateRandomMPT() //pow
			mpt = GenerateTransactionsMPT() // txs in Mpt for p5
			fmt.Println("in tryingNonces : parentHash : ", parentHash)
			keyValPair := mpt.GetAllKeyValuePairs()
			if len(keyValPair) == 0 {
				GetNewParent = true
				time.Sleep(3 * time.Second)
				continue
			}

			nonce = pow.InitializeNonce(8)

		}

		if pow.POW(parentHash, nonce, mpt.Root, difficulty) {
			//generate block send heartbeat (blockBeat)
			SendBlockBeat(tryingForHeight, parentHash, nonce, mpt)

			//MarkTxInTxPoolAsUsed(mpt) // marking transaction true(used in block) //todo

			GetNewParent = true
		}

		nonce = pow.InitializeNonce(8) //NextNonce(Nonce)

		//if tryingForHeight < SBC.GetLength() { //?? todo if the sbc length has increased
		//	GetNewParent = true
		//}

	}

}

// SendBlockBeat func prepares heartbeat data and sends across to peers
func SendBlockBeat(height int32, parentHash string, nonce string, mpt mptp.MerklePatriciaTrie) {

	//log.Println("-------------------- sending MPT ROOT at height : ", mpt.Root)
	//y := sha3.Sum256([]byte(parentHash + nonce + mpt.Root))
	//y1 := hex.EncodeToString(y[:])
	//log.Println("-------------------- sending proof : ", y1)

	Peers.Rebalance()
	peerMap := Peers.Copy()
	PeerMapJson, _ := Peers.PeerMapToJson()

	b1 := SBC.GenBlock(height, parentHash, mpt, nonce, ID.GetMyPublicIdentity())
	SBC.Insert(b1)
	blockJson := b.EncodeToJSON(&b1)

	peerMapPid := Peers.CopyPids()
	PeerMapPidJson, err := gp.PeerMapPidToJson(peerMapPid) //p5
	if err != nil {
		log.Println("in sendBlockBeat, Error in converting peerMapPid to json : err - ", err)
	}

	signForBlockJson := ID.GenSignature([]byte(blockJson))

	heartBeat := gp.PrepareHeartBeatData(&SBC, Peers.GetSelfId(), ID.GetMyPublicIdentity(), PeerMapJson, PeerMapPidJson, SELF_ADDR, true, blockJson, signForBlockJson)

	//list over peers and send them heartBeat
	if len(peerMap) > 0 {
		for peerAddr := range peerMap {

			_, err := http.Post(peerAddr+"/heartbeat/receive", "application/json; charset=UTF-8",
				strings.NewReader(heartBeat.EncodeToJson())) //apr4

			//encrypting heartbeat json with peerpublic id //Todo crypto
			////EncryptMessageWithPublicKey(hash1 hash.Hash, publicKey *rsa.PublicKey, message string, label string)
			//peerPid := p5.PublicIdentity(peerMapPid[peerAddr]).PublicKey
			//ciphered := p5.EncryptMessageWithPublicKey( peerPid, heartBeat.EncodeToJson() )
			//_, err := http.Post(peerAddr+"/heartbeat/receive", "application/octet-stream; charset=UTF-8",
			//	strings.NewReader(string(ciphered))) //apr30 //p5
			if err != nil {
				Peers.Delete(peerAddr)
				//fmt.Println("deleting peer : ", peer)
			}
		}
	}
}

///////////////////////////////

func ShowWallet(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "%s\n", Wallet.Show())
	if err != nil {
		log.Println("Err in ShowWallet func while writing response")
	}
}

func ShowBalanceBook(w http.ResponseWriter, r *http.Request) {

	bb := p5.NewBalanceBook()
	chain := pow.GetCanonicalChains(&SBC)
	bb.BuildBalanceBook(chain[0], 2)

	str := "Balance Book :\n"
	for key, value := range bb.Book.GetAllKeyValuePairs() {
		str += "\n" + key + "\t\t" + value + " " + p5.TOKENUNIT
	}
	str += "\n\n"
	//
	//str += "\n\nPromised :\n"
	//for _, btx := range bb.Promised {
	//
	//	log.Println("In ShowBalanceBook - Borrowing tx - ",btx.BorrowingTxId)
	//	tokensNeeded := fmt.Sprintf("%f", btx.BorrowingTx.Tokens)
	//	str += "Borrowing Tx : "+ btx.BorrowingTxId + " -> "+ "required amount : " + tokensNeeded +"\n"
	//
	//	for _, ptx := range btx.PromisesMade {
	//
	//		log.Println("In ShowBalanceBook - Promise made by "+ ptx.From.Label +"  of -->  ",ptx.Tokens)
	//		str +="\t\t" + "Promise made by "+ ptx.From.Label +" of -->  " + strconv.FormatFloat(ptx.Tokens, 'f', 6, 64) + p5.TOKENUNIT + "\n"
	//	}
	//}
	str += bb.ShowPromised()
	fmt.Printf("\nIn ShowBalanceBook \n %s\n", str)

	_, err := fmt.Fprint(w, str)
	if err != nil {
		log.Println("Err in ShowWallet func while writing response")
	}
}

func ShowTransactionPool(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "%s\n", TxPool.Show())
	if err != nil {
		log.Println("Err in ShowWallet func while writing response")
	}
}

////Transaction func takes http request POST /transaction - checks if the transaction is valid and adds to its TxPool - sendTXBeat
//func Transaction(w http.ResponseWriter, r *http.Request) {
//
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Println("Err : in HeartBeatReceive - reached err of ioutil.ReadAll -")
//		log.Println(err)
//	}
//	defer r.Body.Close()
//
//	tx := p5.DecodeToTransaction(body)
//
//	//check if transaction is valid
//	if p5.IsTransactionValid(tx, BalanceBook) {
//		//put transaction in Txpool
//		TxPool.AddToTransactionPool(tx)
//	}
//
//	go sendTxBeat(tx)
//
//}
//
//// sendTxBeat func takes a transaction and sends to its peer list
//func sendTxBeat(tx p5.Transaction) {
//	//send TransactionBeat
//	txBeat := p5.PrepareTransactionBeat(tx, &ID)
//
//	Peers.Rebalance()
//	peerMap := Peers.Copy()
//	//list over peers and send them heartBeat
//	if len(peerMap) > 0 {
//		for peerAddr := range peerMap {
//			_, _ = http.Post(peerAddr+"/txbeat/receive", "application/json; charset=UTF-8",
//				strings.NewReader(txBeat.EncodeToJson()))
//		}
//	}
//}

// called when /showBlockMpt/{height}/{hash} is called
func ShowBlockMpt(w http.ResponseWriter, r *http.Request) {
	log.Println("In ShowBlockMpt")
	vars := mux.Vars(r)
	ubHeight, err := strconv.Atoi(vars["height"])
	if err != nil {
		returnCode500(w, r)
	} else {
		//ubHash := vars["hash"]
		//fmt.Println("\nuploading block of -\nubHeight : ", ubHeight)
		//fmt.Println("ubHash : ", ubHash, "\n\n")

		//uBlock, found := SBC.GetBlock(int32(ubHeight), ubHash)
		uBlocks, found := SBC.Get(int32(ubHeight))
		if found == false {
			fmt.Println("Err : in Handlers - ShowBlockMpt - found = false - 204")
			returnCode204(w, r)
		} else {
			fmt.Println("in Handlers - ShowBlockMpt - found = true")

			//use uBlock to get all key value pairs
			blk := b.Block(uBlocks[0])
			mpt := mptp.MerklePatriciaTrie(blk.Value)
			keyValuePairs := mpt.GetAllKeyValuePairs()
			var outputBuilder strings.Builder
			outputBuilder.WriteString(blk.Header.Hash + "\n")
			for key, value := range keyValuePairs {
				//fmt.Fprintf(&outputBuilder, "%s :", key)
				//fmt.Fprintf(&outputBuilder, ": %s\n", value)
				outputBuilder.WriteString("\n" + "key=>" + key + "\n value=>" + value + "\n")
			}

			//blockJson := block.EncodeToJSON(&uBlock)
			//_, err = fmt.Fprint(w, blockJson)
			_, err = fmt.Fprint(w, outputBuilder.String())
			if err != nil {
				log.Println("Err : in Handlers - ShowBlockMpt - during writing response")
			}
		}
	}

}

func ShowBlock(w http.ResponseWriter, r *http.Request) {
	log.Println("In ShowBlock")
	vars := mux.Vars(r)
	ubHeight, err := strconv.Atoi(vars["height"])
	if err != nil {
		returnCode500(w, r)
	} else {
		uBlocks, found := SBC.Get(int32(ubHeight))
		if found == false {
			fmt.Println("Err : in Handlers - ShowBlock - found = false - 204")
			returnCode204(w, r)
		} else {
			fmt.Println("in Handlers - ShowBlock - found = true")
			blk := b.Block(uBlocks[0])
			str := b.EncodeToJSON(&blk)

			_, err = fmt.Fprint(w, str)
			if err != nil {
				log.Println("Err : in Handlers - ShowBlockMpt - during writing response")
			}

		}

	}

}
