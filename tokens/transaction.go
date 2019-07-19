package tokens

import (
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"log"
	"strconv"
	"time"
	//s "../p5security"
	"../client"
	s "../identity"
)

type Transaction struct {
	Id        string           `json:"id"`
	From      s.PublicIdentity `json:"from"`
	To        s.PublicIdentity `json:"to"` //if To is empty then its a borrowing tx
	ToTxId    string           `json:"toTxId"`
	Tokens    float64          `json:"tokens"`
	Fees      float64          `json:"fees"`
	TxType    string           `json:"txtype"`
	Timestamp time.Time        `json:"timestamp"`
}

func NewTransaction(from s.PublicIdentity, to s.PublicIdentity, toTxId string, tokens float64, fees float64, txType string) Transaction {
	tx := Transaction{
		From:      from,
		To:        to,
		ToTxId:    toTxId,
		Tokens:    tokens,
		Fees:      fees,
		TxType:    txType,
		Timestamp: time.Now(),
	}

	tx.Id = tx.genId()

	return tx
}

func (tx *Transaction) genId() string {
	str := tx.From.PublicIdentityToJson() +
		tx.To.PublicIdentityToJson() +
		tx.ToTxId +
		strconv.FormatFloat(float64(tx.Tokens), 'f', -1, 64) +
		strconv.FormatFloat(float64(tx.Fees), 'f', -1, 64) +
		tx.TxType +
		tx.Timestamp.String()
	sum := sha3.Sum256([]byte(str))
	return hex.EncodeToString(sum[:])
}

func (tx *Transaction) Show() string {
	str := "\ntx id :" + tx.Id +
		"\ntx From :" + tx.From.PublicIdentityToJson() +
		"\ntx To :" + tx.To.PublicIdentityToJson() +
		"\ntx ToTxId :" + tx.ToTxId +
		"\ntx Tokens :" + strconv.FormatFloat(float64(tx.Tokens), 'f', -1, 64) +
		"\ntx Fees :" + strconv.FormatFloat(float64(tx.Fees), 'f', -1, 64) +
		"\ntx Type :" + tx.TxType +
		"\ntx Time :" + tx.Timestamp.String() + "\n"
	return str
}

func (tx *Transaction) CreateTxSig(fromCid client.ClientId) []byte {
	return fromCid.GenSignature(tx.TransactionToJsonByteArray())
}

func (tx *Transaction) CreateTxSigForMiner(fromId s.Identity) []byte {
	return fromId.GenSignature(tx.TransactionToJsonByteArray())
}

func VerifyTxSig(fromPid s.PublicIdentity, tx Transaction, txSig []byte) bool {

	return s.VerifySingature(fromPid.PublicKey, tx.TransactionToJsonByteArray(), txSig)
}

func (tx *Transaction) TransactionToJsonByteArray() []byte {
	txJson, err := json.Marshal(tx)
	if err != nil {
		log.Println("in TransactionToJsonByteArray : Error in marshalling Tx : ", err)
	}

	return txJson
}

func (tx *Transaction) TransactionToJson() string {
	txJson, err := json.Marshal(tx)
	if err != nil {
		log.Println("in TransactionToJsonByteArray : Error in marshalling Tx : ", err)
	}

	return string(txJson)
}

func JsonToTransaction(txJson string) Transaction {
	tx := Transaction{}
	err := json.Unmarshal([]byte(txJson), &tx)
	if err != nil {
		log.Println("Error in unmarshalling Transaction, err - ", err)
		log.Println("String given to unmarshall Transaction, ================> \n ", txJson, "\nxxxxxxxxxxxxxxxxxxxxxxxxxxx\n")
	}

	return tx
}

//func IsTransactionValid(tx Transaction, balanceBook BalanceBook) bool {
//
//	//balanceBook.Book <hash of PublicKey, Balance Amount>
//	//getting hash of public key of tx.From - to get key for balance.Book
//	//hash :=sha3.Sum256(tx.From.PublicKey.N.Bytes())
//	//hashKey := hex.EncodeToString(hash[:])
//	//using hashKey to get the Balance amount
//	//balanceStr, err := balanceBook.Book.Get(hashKey)
//	//balance, err := strconv.ParseFloat(balanceStr, 64) // todo ?? if ERR then should i make balance zero ???? !!!
//	//if err != nil {
//	//	return false
//	//}
//
//	//if  balance > tx.Tokens {
//	//	return true
//	//}
//	//return false
//
//	return balanceBook.IsBalanceEnough(balanceBook.GetKey(tx.From.PublicKey), tx.Tokens+tx.Fees)
//}
