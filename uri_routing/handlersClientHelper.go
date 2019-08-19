package uri_routing

import (
	p5 "../client"
	tkn "../tokens"
	//s "../p5security"
	s "../identity"
	"../wallet"
	"bytes"
	"os"
	"strconv"

	//"crypto/rsa"

	//"./data"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	//"sync"
	"time"
)

var CID p5.ClientId

//can use here --> var Peers data.PeerList
var BCH p5.BlockChainHolders

//  it 	has SELF_ADDR   -> from Init func
// 		has INIT_SERVER -> localhost/6686

//var Wallet p5.Wallet
var Wallet wallet.Wallet

func NewClient(w http.ResponseWriter, r *http.Request) {

	client := &http.Client{}

	//Wallet = p5.NewWallet()

	resp, err := client.Get(INIT_SERVER + "/client")
	if err != nil {
		w.WriteHeader(404)
		_, _ = fmt.Fprintf(w, "<h1>Page not found</h1>")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error in reading response body in NewClient, err - ", err)
	}

	w.WriteHeader(200)
	_, _ = fmt.Fprintf(w, string(body))

	go startBcHolderUpdate()
}

func startBcHolderUpdate() {

	log.Println("In startBcHolderUpdate vvvv >>>> VVVV")
	BCH = p5.NewBlockChainHolders()

	for true {
		resp, err := http.Get(INIT_SERVER + "/bcholders")
		if err != nil {
			log.Println("Error in fetching BC holders at client, err - ", err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error in reading response body in BcHolderUpdate, err - ", err)
		}

		bch := p5.BlockChainHolders{}
		err = json.Unmarshal(body, &bch.Holders)
		if err != nil {
			log.Println("Error in unmarshalling blockchain holders, err - ", err)
		}

		for holderaddr, holderPid := range bch.Holders {
			BCH.AddBlockChainHolder(holderaddr, holderPid)
		}

		time.Sleep(60 * time.Second)
	}

}

func ShowBcHolders(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	fmt.Fprintf(w, BCH.Show())
}

func SignUp(w http.ResponseWriter, r *http.Request) {

	//err := r.ParseForm()
	//if err != nil {
	//	fmt.Println("Error in parsing the signup form : err - ", err)
	//}
	//username := r.FormValue("phrase")

	var body []byte
	for holderaddr := range BCH.Holders {
		resp, err := http.Post(holderaddr+"/clientsignup", "application/x-www-form-urlencoded", r.Body)
		if err != nil {
			log.Println("Error in SignUp, err - ", err)
			continue
		} else {
			body, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Println("Error in SignUp, err - ", err)
				continue
			}
			break
		}
	}

	_, _ = fmt.Fprintf(w, string(body))

}

func Login(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error in GetBody ----------> err: ,", err)
	}

	var respBody []byte
	var loggedIn = false

	for holderaddr := range BCH.Holders {
		resp, err := http.Post(holderaddr+"/clientlogin", "application/x-www-form-urlencoded", ioutil.NopCloser(bytes.NewBuffer(reqBody)))
		if err != nil {
			log.Println("Error in Login, err - ", err)
			continue
		} else {
			respBody, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Println("Error in reading response, in Login, err - ", err)
				continue
			} else {
				loggedIn = true
			}
			break
		}
	}

	if loggedIn == true {
		log.Println("Setting CID from login route")
		_, _ = http.Post(SELF_ADDR+"/cidset", "application/x-www-form-urlencoded", ioutil.NopCloser(bytes.NewBuffer(reqBody)))

	}

	_, _ = fmt.Fprintf(w, string(respBody))
}

func CIDSet(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error in parsing the login form : err - ", err)
	}

	cidJson := r.FormValue("key")
	cid, err := p5.JsonToClientId(cidJson)
	if err == nil {
		CID = cid
	}

	pid := CID.GetMyPublicIdentity()
	pidJsonString := pid.PublicIdentityToJson()
	str := "CID :\n" + string(CID.ClientIdToJsonByteArray()) + "\n\nPID :\n" + pidJsonString
	log.Println("\n\n\n\n", str, "\n\n\n")

}

func TransactionForm(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error in parsing the Transaction Form : err - ", err)
	}

	fromVal := CID.GetMyPublicIdentity()                                             //r.FormValue("from")
	log.Println("in TransactionForm - FromValue : ", fromVal.PublicIdentityToJson()) //empth here
	//fromVal := p5.JsonToPublicIdentity(fromValue)
	toFormValue := r.FormValue("to")
	log.Println()
	toVal := s.JsonToPublicIdentity(toFormValue)
	fmt.Println("HEHAAH : ", toFormValue)
	fmt.Println("HEHAAH : ", toVal.Label)
	toTxId := r.FormValue("txid")

	amountVal, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		log.Println("Error in Amount conversion : err - ", err)
	}

	feesVal, err := strconv.ParseFloat(r.FormValue("fees"), 64)
	if err != nil {
		log.Println("Error in Amount conversion : err - ", err)
	}

	var tx tkn.Transaction

	log.Println("================================")

	if toVal.Label == "" && toTxId == "" && amountVal >= 0 && feesVal >= 0 {
		tx = tkn.NewTransaction(fromVal, toVal, toTxId, amountVal, feesVal, "req")
		log.Println("================================ Requirement Tx")
	} else if toTxId != "" && toVal.Label != "" /*&& amountVal >= 0*/ && feesVal >= 0 {

		tx = tkn.NewTransaction(fromVal, toVal, toTxId, amountVal, feesVal, "promise")
		log.Println("================================ Promise Tx")
	} else if amountVal >= 0 && feesVal >= 0 {
		tx = tkn.NewTransaction(fromVal, toVal, toTxId, amountVal, feesVal, "")
		log.Println("================================ Normal Tx")
	}

	log.Println("================================ toVal.Label : ", toVal.Label)
	log.Println("================================ toTxId : ", toTxId, " = ", tx.ToTxId)
	log.Println("================================ toTxId : ", tx.TxType)
	log.Println("================================")

	txBeat := tkn.NewTransactionBeat(tx, fromVal, tx.CreateTxSig(CID))
	txBeatJson := txBeat.EncodeToJsonByteArray()

	var resp *http.Response
	for holderaddr := range BCH.Holders {
		resp, err = http.Post(holderaddr+"/txbeat/receive", "application/json", ioutil.NopCloser(bytes.NewBuffer(txBeatJson)))
		if err != nil {
			log.Println("Error in sending txBeat from client to bcHolders, err -", err)
		}
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("in TransactionForm - Error in reading response from bcHolder")
	}
	fmt.Fprint(w, string(respBody))
}

func CIDPage(w http.ResponseWriter, r *http.Request) {
	cwd, _ := os.Getwd()
	filePath := cwd + "/resource/html/toSetCid.html"
	http.ServeFile(w, r, filePath)

}

func SetCID(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error in parsing the login form : err - ", err)
	}
	cidJson := r.FormValue("key")
	cid, err := p5.JsonToClientId(cidJson)
	if err == nil {
		CID = cid
	}

	pid := CID.GetMyPublicIdentity()
	pidJsonString := pid.PublicIdentityToJson()
	str := "CID :\n" + string(CID.ClientIdToJsonByteArray()) + "\n\nPID :\n" + pidJsonString
	_, _ = fmt.Fprint(w, str)
	//GetMyId(w,r)
}

func GetMyId(w http.ResponseWriter, r *http.Request) {
	pid := CID.GetMyPublicIdentity()
	pidJsonString := pid.PublicIdentityToJson()
	fmt.Fprint(w, pidJsonString)

}
