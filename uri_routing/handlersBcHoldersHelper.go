package uri_routing

import (
	p5 "../client"
	tkn "../tokens"
	//"../p4"
	"../pow"
	//s "../p5security"
	s "../identity"
	"../resource"
	//t "../tokens"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"strings"

	//"io/ioutil"
	//"log"
	"net/http"
	"os"
)

func ServeClient(w http.ResponseWriter, r *http.Request) {
	//cwd, _ := os.Getwd()
	//filePath := cwd + "/resource/html/cover.html"
	//w.Header().Set("Content-type", "text/html")
	//http.ServeFile(w, r, filePath)
	coverHtml(w, r)

}

func BcHolders(w http.ResponseWriter, r *http.Request) {

	bcHolderPids := Peers.CopyPids()
	bcHolderPids[SELF_ADDR] = ID.GetMyPublicIdentity()

	peersPidJson, err := json.Marshal(bcHolderPids)
	if err != nil {
		peersPidJson = []byte("{}")
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, string(peersPidJson))

}

func ClientSignUp(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error in parsing the signup form in ClientSignUp : err - ", err)
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "Bad Request, PLease try again")
	}

	username := r.FormValue("phrase")
	fmt.Println("Signup request from user - -- ", username)

	id := p5.NewClientId(username)

	GiveDefaultTokens(id) //giveDefaultTokens

	jsonBytesId := id.ClientIdToJsonByteArray()
	w.WriteHeader(200)
	_, _ = fmt.Fprintf(w, "Whole Key, please save this :\n"+string(jsonBytesId))

}

func ClientLogin(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "login req receieved")

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error in parsing the login form : err - ", err)
	}
	phrase := r.FormValue("phrase")
	cidJson := r.FormValue("key")
	cid, err := p5.JsonToClientId(cidJson)
	if err == nil {
		sig := cid.GenSignature([]byte(phrase))
		verified := s.VerifySingature(cid.PublicKey, []byte(phrase), sig)
		//_, _ = fmt.Fprintf(w, "Verified ?? :\n" + strconv.FormatBool(verified))

		if verified {

			//obj := resource.UserLandingPage{}
			//obj.Pid = cid.GetMyPublicIdentity()
			//obj.FromPid = obj.Pid.PublicIdentityToJson()
			//
			//clientLandingHtml(w,r, obj)
			pid := cid.GetMyPublicIdentity()
			clientLandingHtml(w, r, pid)

		} else {

			coverHtml(w, r)
		}

	} else {

		coverHtml(w, r)
	}

}

//func ClientLogin(w http.ResponseWriter, r *http.Request){
//	//fmt.Fprintf(w, "login req receieved")
//
//	err := r.ParseForm()
//	if err != nil {
//		fmt.Println("Error in parsing the login form : err - ", err)
//	}
//	phrase := r.FormValue("phrase")
//	cidJson := r.FormValue("key")
//	cid, err := p5.JsonToClientId(cidJson)
//	if err == nil {
//		sig := cid.GenSignature([]byte(phrase))
//		verified := p5.VerifySingature(cid.PublicKey, []byte(phrase), sig)
//		//_, _ = fmt.Fprintf(w, "Verified ?? :\n" + strconv.FormatBool(verified))
//
//		if verified {
//
//			obj := resource.UserLandingPage{}
//			obj.FromPid = p5.PublicIdentityToJson(cid.GetMyPublicIdentity())
//
//			clientLandingHtml(w,r, obj)
//
//		} else {
//
//			coverHtml(w,r)
//		}
//
//	} else {
//
//		coverHtml(w, r)
//	}
//
//}

func coverHtml(w http.ResponseWriter, r *http.Request) {
	cwd, _ := os.Getwd()
	filePath := cwd + "/resource/html/cover.html"
	http.ServeFile(w, r, filePath)
}

func clientLandingHtml(w http.ResponseWriter, r *http.Request, pid s.PublicIdentity /*cid p5.ClientId/*obj resource.UserLandingPage*/) {

	obj := resource.UserLandingPage{}
	obj.Pid = pid                            //cid.GetMyPublicIdentity()
	obj.FromPid = pid.PublicIdentityToJson() //obj.Pid.PublicIdentityToJson()
	chains := pow.GetCanonicalChains(&SBC)
	obj.BTxs = tkn.BuildBorrowingTransactions(chains)
	bb := tkn.NewBalanceBook()
	bb.BuildBalanceBook(chains[0], 2)
	obj.BB = bb
	obj.PromisedInString = bb.ShowPromised()
	obj.Purse = tkn.NewWallet()
	obj.Purse.Balance = bb.GetBalanceFromPublicKey(pid.PublicKey)

	cwd, _ := os.Getwd()
	tmpl, err := template.ParseFiles(cwd + "/resource/html/user_landing_page.html")
	if err == nil {

		w.WriteHeader(200)
		w.Header().Set("Content-type", "application/html")
		tmpl.Execute(w, obj)
	}
}

// TransactionBeatRecv func takes http request POST /txbeat/receive - verify txSig - verify tx valid - add to TxPool - forward to Peers
func TransactionBeatRecv(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Err : in TransactionBeatRecv - reached err of ioutil.ReadAll -")
		log.Println(err)
	}
	defer r.Body.Close()

	txBeat := tkn.DecodeToTransactionBeat(body)

	if tkn.VerifyTxSig(txBeat.Tx.From, txBeat.Tx, txBeat.TxSig) {
		//check if transaction is valid todo verification in - when getting tx for mpt
		if txBeat.Tx.Tokens != 0 {
			//if p5.IsTransactionValid(txBeat.Tx, BalanceBook) { //checks both book and amt promised //todo
			//put transaction in Txpool
			log.Println("In TransactionBeatRecv - Signature verified !!!!!!!!!!!!!!!!!!!")
			TxPool.AddToTransactionPool(txBeat.Tx)

			go forwardTxBeat(txBeat)

		}
	}

	pid := txBeat.FromPid
	clientLandingHtml(w, r, pid)
}

func forwardTxBeat(txBeat tkn.TransactionBeat) {
	//forward TransactionBeat
	log.Println("In fowrardTxBeat =========>")
	txBeat.Hops--
	if txBeat.Hops >= 0 {
		Peers.Rebalance()
		peerMap := Peers.Copy()
		//list over peers and send them heartBeat
		if len(peerMap) > 0 {
			for peerAddr := range peerMap {
				log.Println("Sending tx ======================> : ", txBeat.Tx.Id)
				log.Println("Sending to ======================> : ", peerAddr)
				_, _ = http.Post(peerAddr+"/txbeat/receive", "application/json; charset=UTF-8",
					strings.NewReader(txBeat.EncodeToJson()))
			}
		}
	}
}

func GetTransactionPool() { // called in start func of bcHolder

	var respBodyBytes []byte

	resp, err := http.Get(INIT_SERVER + "/txbeat/allprev")
	if err != nil {
		log.Println("Error in asking for TxPool, err - ", err)
	} else {
		respBodyBytes, err = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		log.Println("Receved TransactionPoolJson json is --- ", string(respBodyBytes), "!!!")
		recvTxPool := tkn.DecodeJsonToTransactionPoolJson(string(respBodyBytes))

		for _, tx := range recvTxPool.Pool {
			TxPool.AddToTransactionPool(tx)
		}
	}

}

func TransactionPoolRecv(w http.ResponseWriter, r *http.Request) {
	log.Println("Asked for Transaction Pool !!!")
	txpj := TxPool.GetTransactionPoolJsonObj()
	jsonstr := txpj.EncodeToJsonTransactionPoolJson()

	log.Println("Sending TransactionPoolJson json is --- ", jsonstr, "!!!")

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, jsonstr)
}

func GiveDefaultTokens(cid p5.ClientId) {
	tx := tkn.NewTransaction(ID.GetMyPublicIdentity(), cid.GetMyPublicIdentity(), "", 1000, 0, "default")
	TxPool.AddToTransactionPool(tx)

	//txBeat := p5.NewTransactionBeat(tx, ID.GetMyPublicIdentity(), tx.CreateTxSigForMiner(ID))
	//forwardTxBeat(txBeat)

}

func GiveMinerTokens(cid s.Identity) {
	tx := tkn.NewTransaction(ID.GetMyPublicIdentity(), cid.GetMyPublicIdentity(), "", 10000, 0, "start")
	TxPool.AddToTransactionPool(tx)

	//txBeat := p5.NewTransactionBeat(tx, ID.GetMyPublicIdentity(), tx.CreateTxSigForMiner(ID))
	//forwardTxBeat(txBeat)

}

//func GiveGenesisTokens() {
//tx := p5.NewTransaction(ID.GetMyPublicIdentity(), cid.GetMyPublicIdentity(), "", 10000, 0, "start")
//TxPool.AddToTransactionPool(tx)
//
////txBeat := p5.NewTransactionBeat(tx, ID.GetMyPublicIdentity(), tx.CreateTxSigForMiner(ID))
////forwardTxBeat(txBeat)
//}
