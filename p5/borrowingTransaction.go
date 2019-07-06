package p5

import (
	//"../p1"
	"../data_structure/mpt"
	//"../p2"
	"../blockchain"
	//b "../p2/block"
	b "../block"
	"encoding/json"
	"log"
)

type BorrowingTransaction struct {
	BorrowingTxId string        `json:"borrowingtxid"`
	BorrowingTx   Transaction   `json:"borrowingtx"`
	PromisesMade  []Transaction `json:"promisesmade"` // key - transaction id (Lending) // todo todo -- changed from map to array
	PromisedValue float64       `json:"promisedvalue"`
}

func NewBorrowingTransaction(tx Transaction) BorrowingTransaction {
	bt := BorrowingTransaction{}
	bt.BorrowingTxId = tx.Id
	bt.BorrowingTx = tx
	bt.PromisedValue = 0.0
	bt.PromisesMade = make([]Transaction, 0)

	return bt

}

type BorrowingTransactions struct {
	BorrowingTxs map[string]Transaction // key - BorrowingTxId value - txJson
	Borrower     map[string]string      //json of pid of borrower
}

func NewBorrowingTransactions() BorrowingTransactions {
	btxs := BorrowingTransactions{}
	btxs.BorrowingTxs = make(map[string]Transaction)
	btxs.Borrower = make(map[string]string)
	return btxs
}

func BuildBorrowingTransactions(chains []blockchain.Blockchain) BorrowingTransactions {

	btx := NewBorrowingTransactions()
	//loop over the blockchain[00 of chains
	var i int32
	if len(chains) > 0 {
		for i = 1; i <= chains[0].Length; i++ {
			blks, found := chains[0].Get(i)
			if found && len(blks) > 0 {
				blk := b.Block(blks[0])
				mpt := mpt.MerklePatriciaTrie(blk.Value)
				keyValuePairs := mpt.GetAllKeyValuePairs() //key - txid value - txJson
				//loop over all key valye pairs and collect borrowing txs
				for _, txjson := range keyValuePairs {
					tx := JsonToTransaction(txjson)
					if tx.TxType == "req" /*tx.To.Label == "" && tx.ToTxId == "" && tx.Tokens > 0 && tx.TxType != "start" && tx.TxType != "default" && tx.From.Label != ""*/ {
						btx.BorrowingTxs[tx.Id] = tx
						btx.Borrower[tx.Id] = tx.From.PublicIdentityToJson()
					}

				}
			}
		}
	}

	return btx

}

func (btx *BorrowingTransaction) EncodeTojsonString() string {
	jsonBytes, err := json.Marshal(btx)
	if err != nil {
		log.Println("Error in marshalling BorrowingTransaction to json , err - ", err)
	}
	return string(jsonBytes)
}
