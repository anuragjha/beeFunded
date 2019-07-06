package p5

import (
	//"../p1"
	"../data_structure/mpt"
	//"../p2"
	"../blockchain"
	//b "../p2/block"
	b "../block"
	//s "../p5security"
	s "../identity"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"

	//"go/types"
	"log"
	"strconv"
	"sync"
	//"sync"
)

type BalanceBook struct {
	Book mpt.MerklePatriciaTrie //key - hashOfPubKey and Value - balance

	// key - Requirement transaction id -||- value - BorrowingTransaction
	Promised map[string]BorrowingTransaction
	mux      sync.Mutex
}

func NewBalanceBook() BalanceBook {
	book := mpt.MerklePatriciaTrie{}
	book.Initial()
	//promised := mpt.MerklePatriciaTrie{}
	//promised.Initial()

	promised := make(map[string]BorrowingTransaction)

	return BalanceBook{
		Book:     book,
		Promised: promised,
	}
}

func (bb *BalanceBook) BuildBalanceBook(chain blockchain.Blockchain, fromHeight int32) { // not using fromHeight for now

	log.Println(">>>>>>>>>>>>>>> In BuildBalanceBook  <<<<<<<<<<<<<<<<")
	//loop over the blockchain[0] of chains
	var i int32

	for i = fromHeight; i <= chain.Length; i++ {
		blks, found := chain.Get(i)
		if found && len(blks) > 0 {
			blk := b.Block(blks[0])
			//mpt := mpt.MerklePatriciaTrie(blk.Value)
			mpt := mpt.MerklePatriciaTrie{}
			mpt.Initial()
			mpt = blk.Value
			keyValuePairs := mpt.GetAllKeyValuePairs() //key - txid value - txJson
			//loop over all key valye pairs and collect borrowing txs
			for keyTxid, valueTxJson := range keyValuePairs {

				//log.Println("\nIn BuildBalanceBook - txJson to be converted to Tx ---> \n", txjson)
				tx := JsonToTransaction(valueTxJson)
				log.Println("\nIn BuildBalanceBook - txJson to be converted to Tx ---> keyTxid\n", keyTxid)
				log.Println("\nIn BuildBalanceBook - txJson to be converted to Tx ---> valueTxJson\n", valueTxJson)
				log.Println("\nIn BuildBalanceBook - txJson to be converted to Tx ---> tx ID\n", tx.Id)
				log.Println("\nIn BuildBalanceBook - txJson to be converted to Tx ---> tx Tokens\n", tx.Tokens)

				bb.UpdateABalanceBookForTx(tx, blk.Header.Miner) // updating BalanceBook for transaction

				//blk.Header.//todo  blk here  // add a mined by in the block // for now do without tx fess
				//minerKey :=  bb.GetKey(blk.Header.miner)  // assuming miner - public key of miner
			}
		}
	}

}

//func (bb *BalanceBook) UpdateABalanceBookForBlock() { // update bb based on transaction here // todo
//	toKey 	:=  bb.GetKey(tx.To.PublicKey)
//	fromKey := 	bb.GetKey(tx.From.PublicKey)
//
//	bb.UpdateABalanceInBook(toKey, -tx.Tokens)
//	bb.UpdateABalanceInBook(toKey, -tx.Fees)
//	bb.UpdateABalanceInBook(fromKey, tx.Tokens)
//}

func (bb *BalanceBook) UpdateABalanceBookForTx(tx Transaction, miner s.PublicIdentity) { // update bb based on transaction here

	log.Println(">>>>>>>>>>>>>>> In UpdateABalanceBookForTx  Overall start <<<<<<<<<<<<<<<<")
	minerKey := bb.GetKey(miner.PublicKey)

	if tx.TxType == "default" {
		log.Println("<<<<<<<<<<<<<<<<<<<<<< In UpdateABalanceBookForTx - default !!!!!!!! !!!!!! - tx id - ", tx.Id,
			">>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		fromKey := bb.GetKey(tx.From.PublicKey)
		bb.UpdateABalanceInBook(fromKey, -tx.Tokens)

		toKey := bb.GetKey(tx.To.PublicKey)
		bb.UpdateABalanceInBook(toKey, tx.Tokens)

		// start of fees
		bb.UpdateABalanceInBook(minerKey, tx.Fees)
		bb.UpdateABalanceInBook(fromKey, -tx.Fees)
		// end of fees

	} else if tx.TxType == "start" { //default transaction read
		log.Println("<<<<<<<<<<<<<<<<<<<<<< In UpdateABalanceBookForTx - start !!!!!!!! !!!!!! - tx id - ", tx.Id,
			">>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		toKey := bb.GetKey(tx.To.PublicKey)
		bb.UpdateABalanceInBook(toKey, tx.Tokens)

	} else if tx.TxType == "req" /*&& tx.ToTxId == ""*/ /*tx.To.Label == "" && tx.ToTxId == "" && tx.From.Label != ""*/ { // A 's Req Tx // Requirement

		if _, ok := bb.Promised[tx.ToTxId]; !ok {
			log.Println("<<<<<<<<<<<<<<<<<<<<<< In UpdateABalanceBookForTx - // A 's Req Tx // Requirement !!!!!!!! !!!!!! - tx id - ", tx.Id,
				">>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			bb.PutTxInPromised(tx) //create a new key in promised
		}

		// start of fees
		bb.UpdateABalanceInBook(minerKey, tx.Fees)
		fromKey := bb.GetKey(tx.From.PublicKey)
		bb.UpdateABalanceInBook(fromKey, -tx.Fees)
		// end of fees

	} else if tx.TxType == "promise" /*&& tx.ToTxId != ""*/ /*tx.ToTxId != "" && tx.From.Label != "" && tx.To.Label == ""*/ { // B to Req Tx

		log.Println("<<<<<<<<<<<<<<<<<<<<<< In UpdatePromiseBookForTx - // B +++++> Req Tx !!!!!!!! !!!!!! - tx id - ", tx.Id,
			">>>>>>>>>>>>>>>>>>>>>>>>>>>>", tx.Tokens)
		//bb1 := NewBalanceBook()
		//bb1.Book = bb.Book
		//bb1.Promised = bb.Promised
		vary := bb.Promised[tx.ToTxId]
		bb.UpdateABalanceInPromised(tx, vary)

		// start of fees
		bb.UpdateABalanceInBook(minerKey, tx.Fees)
		fromKey := bb.GetKey(tx.From.PublicKey)
		bb.UpdateABalanceInBook(fromKey, -tx.Fees)
		// end of fees

	} else if tx.TxType == "" && tx.To.Label != "" && tx.From.Label != "" && tx.ToTxId == "" { // A to B token transfer // pay
		log.Println("<<<<<<<<<<<<<<<<<<<<<< In UpdateABalanceBookForTx -A to B token transfer !!!!!!!! !!!!!! - tx id - ", tx.Id,
			">>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		fromKey := bb.GetKey(tx.From.PublicKey)
		bb.UpdateABalanceInBook(fromKey, -tx.Tokens)

		toKey := bb.GetKey(tx.To.PublicKey)
		bb.UpdateABalanceInBook(toKey, tx.Tokens)

		// start of fees
		bb.UpdateABalanceInBook(minerKey, tx.Fees)
		bb.UpdateABalanceInBook(fromKey, -tx.Fees)
		// end of fees
	}
}

func (bb *BalanceBook) UpdateABalanceInBook(PublicKeyHashStr string, updateBalanceBy float64) {
	bb.mux.Lock()
	defer bb.mux.Unlock()

	log.Println(">>>>>>>>>>>>>>> In UpdateABalanceInBook  <<<<<<<<<<<<<<<<")
	//pubIdHashStr := GetHashOfPublicKey(pubId) //hashOfPublicKey

	currBalance := bb.GetBalanceFromKey(PublicKeyHashStr)
	//log.Println(">>>>>>>>>>>>>>> In UpdateABalanceInBook - currBalance :", currBalance," <<<<<<<<<<<<<<<<")

	newBalance := currBalance + updateBalanceBy
	//log.Println(">>>>>>>>>>>>>>> In UpdateABalanceInBook - newBalance :", newBalance," <<<<<<<<<<<<<<<<")

	bb.Book.Insert(PublicKeyHashStr, fmt.Sprintf("%f", newBalance))
}

func (bb *BalanceBook) UpdateABalanceInPromised(tx Transaction, btx BorrowingTransaction) { //todo check coz changed to array - var PromisesMade

	log.Println(">>>>>>>>>>>>>>>  In UpdateABalanceInPromised  <<<<<<<<<<<<<<<<") // todo todo ::: getting correct tx but not able to put it on the map
	log.Println("\nTransaction being processed : \n", tx.TransactionToJson())     // todo todo --- start with testing - to test the changes made
	log.Println("\nAnd Promised dataStructure is >---->>>", bb.Promised)          // []Transaction init when borrowing tx created in NewBorrowingTransaction
	log.Println("\nAnd ShowPromised is >---->>>", bb.ShowPromised())              // in PutTxInPromised

	// check if somebody wants to take back promise - full  or partial
	if tx.Tokens < 0 && btx.PromisedValue < btx.BorrowingTx.Tokens {

		//check if the tx.From has promised a Sum equal to more than that of cancelling amount
		amountPromised := bb.CheckAmountPromisedByOne(tx.From)
		if amountPromised > -1*tx.Tokens {
			//taking promise back
			btx.PromisedValue += tx.Tokens
			btx.PromisesMade = append(btx.PromisesMade, tx)
			bb.Promised[tx.ToTxId] = btx
		}

	} else if tx.Tokens > 0 {
		//default working
		btx.PromisedValue += tx.Tokens
		btx.PromisesMade = append(btx.PromisesMade, tx)
		bb.Promised[tx.ToTxId] = btx
	}

	//enough := btx.CheckForEnoughPromises()
	if btx.PromisedValue >= btx.BorrowingTx.Tokens { //enough promises made
		//transfer token from Promised Tx User -to- Req Tx User
		log.Println("Enough Promises ----------> -----------> ------->", "Achived")

		bb.TransferPromisesMade(btx)

		delete(bb.Promised, tx.ToTxId)
	}

	// Promised --->
	// key - Requirement transaction id -||- value - BorrowingTransaction
	////if _, ok := bb.Promised[tx.ToTxId]; !ok {
	//
	//	btx := bb.Promised[tx.ToTxId]
	//
	//	//btx.PromisesMade[tx.Id] = tx
	//	btx.PromisesMade = append(btx.PromisesMade, tx)
	//
	//	log.Println(">>>>>>>>>>>>>>> In UpdateABalanceInPromised - PromisesMade : ", bb.Promised[tx.ToTxId].PromisesMade)
	//
	//	enough := btx.CheckForEnoughPromises()
	//	if enough {
	//		//transfer token from Promised Tx User -to- Req Tx User
	//		log.Println("Enough Promises ----------> -----------> ------->", enough)
	//
	//		bb.TransferPromisesMade(btx)
	//
	//		delete(bb.Promised, tx.ToTxId)
	//	}
	//
	////}
}

func (bb *BalanceBook) TransferPromisesMade(btx BorrowingTransaction) { //todo check coz changed to array - var PromisesMade

	log.Println(">>>>>>>>>>>>>>> In TransferPromisesMade  <<<<<<<<<<<<<<<<")
	for _, ptx := range btx.PromisesMade {
		// for subtract
		KeyForSub := GetKeyForBook(ptx.From.PublicKey)
		log.Println("KeyForSub (should be present in ShowBlance) ----->", KeyForSub)
		bb.UpdateABalanceInBook(KeyForSub, -ptx.Tokens)

		//for add
		KeyForAdd := GetKeyForBook(btx.BorrowingTx.From.PublicKey)
		log.Println("KeyForAdd (should be present in ShowBlance) ----->", KeyForAdd)
		bb.UpdateABalanceInBook(KeyForAdd, ptx.Tokens)
	}

}

func (btx *BorrowingTransaction) CheckForEnoughPromises() bool {

	log.Println(">>>>>>>>>>>>>>> In CheckForEnoughPromises  <<<<<<<<<<<<<<<<")

	valueNeeded := btx.BorrowingTx.Tokens

	var valuePromised float64
	valuePromised = 0.0

	for _, ptx := range btx.PromisesMade {
		valuePromised += ptx.Tokens
	}

	log.Println("Value Needed - ---- > ", valueNeeded)
	log.Println("Value Promised - ---- > ", valuePromised)
	if valueNeeded <= valuePromised {
		return true
	}
	return false
}

func (bb *BalanceBook) PutTxInPromised(tx Transaction) {
	bb.mux.Lock()
	defer bb.mux.Unlock()

	log.Println(">>>>>>>>>>>>>>> In PutTxInPromised  <<<<<<<<<<<<<<<<")

	// Promised --->
	// key - Requirement transaction in json -||- value - MPT of <To Txs in json, value>
	if _, ok := bb.Promised[tx.Id]; !ok {

		btx := NewBorrowingTransaction(tx)

		log.Println("\nCheck Start !!! Check Start !!! Check Start !!!\n\n")

		log.Println("BTX is - > ", btx.EncodeTojsonString())

		log.Println("\n\nCheck End !!! Check End !!! Check End !!!")
		log.Println("^^^^")

		bb.Promised[tx.Id] = btx

		log.Println(">>>>>>>>>>>>>>> In PutTxInPromised - bb.Promised for btx : ", bb.Promised[tx.Id].BorrowingTx.Tokens)
	}

}

func GetKeyForBook(publicKey *rsa.PublicKey) string {

	hash := sha3.Sum256(publicKey.N.Bytes())
	hashKey := hex.EncodeToString(hash[:])
	return hashKey
}

// GetKey func takes in PublicKey and returns Key for Book
func (bb *BalanceBook) GetKey(publicKey *rsa.PublicKey) string {
	hash := sha3.Sum256(publicKey.N.Bytes())
	hashKey := hex.EncodeToString(hash[:])
	return hashKey
}

func (bb *BalanceBook) GetBalanceFromPublicKey(publicKey *rsa.PublicKey) float64 {

	//pubIdHashStr := GetHashOfPublicKey(pubId) //hashOfPublicKey
	PublicKeyHashStr := bb.GetKey(publicKey)
	//balance, err := bb.Book.Get(PublicKeyHashStr)
	balance, err := bb.Book.Get(PublicKeyHashStr)
	if err != nil {
		log.Println("GetBalanceFromPublicKey - Error In GetBalance returning 0 , err : ", err)
		return float64(0)
	}

	bal, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		log.Println("Error in strconv from string to float returning 0, err :", err)
		return float64(0)
	}
	return bal

}

func (bb *BalanceBook) GetBalanceFromKey(PublicKeyHashStr string) float64 {

	//pubIdHashStr := GetHashOfPublicKey(pubId) //hashOfPublicKey
	//PublicKeyHashStr := bb.GetKey(publicKey)
	//balance, err := bb.Book.Get(PublicKeyHashStr)
	balance, err := bb.Book.Get(PublicKeyHashStr)
	if err != nil {
		log.Println("GetBalanceFromKey - Error In GetBalance returning 0 , err : ", err)
		return float64(0)
	}

	bal, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		log.Println("Error in strconv from string to float returning 0, err :", err)
		return float64(0)
	}
	return bal

}

//func (bb *BalanceBook) IsBalanceEnough(PublicKeyHashStr string, balanceNeeded float64) bool {
//	currentBalance := bb.GetBalanceFromKey(PublicKeyHashStr) - bb.GetBalanceFromKey(PublicKeyHashStr, bb.Promised)
//	if currentBalance >= balanceNeeded {
//		return true
//	}
//	return false
//}

func (bb *BalanceBook) Show() string {
	return bb.Book.String()
}

func (bb *BalanceBook) ShowPromised() string {
	str := ""
	for _, btx := range bb.Promised {
		str += "#  Asked sum : >>---> " + btx.BorrowingTxId + " of amt : " + fmt.Sprintf("%f", btx.BorrowingTx.Tokens) + " " + TOKENUNIT + "\n"
		for _, promise := range btx.PromisesMade {
			str += "\n-->   Promised sum -> " + promise.From.Label + " has promised : " + fmt.Sprintf("%f", promise.Tokens) + " " + TOKENUNIT + " #  \n"
		}
	}
	return str
}

func (bb *BalanceBook) CheckAmountPromisedByOne(pid s.PublicIdentity) float64 {

	amountPromised := float64(0.0)
	for _, btx := range bb.Promised {
		for _, ptx := range btx.PromisesMade {
			if ptx.From.PublicKey.N.String() == pid.PublicKey.N.String() { //ptx.From   and ptx.Tokens
				amountPromised += ptx.Tokens
			}
		}
	}
	return amountPromised

}
