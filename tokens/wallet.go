package tokens

import (
	//"../p1"
	//"../p2"
	//b "../p2/block"
	//"crypto/rsa"
	"fmt"
	//"golang.org/x/crypto/sha3"
	"sync"
)

//contains funcs for maintaing wallet

//type Currency struct {
//	Value float64
//	Unit  string
//}

const TOKENUNIT = "pingala"

type Wallet struct {
	Balance float64
	Unit    string
	mux     sync.Mutex
}

func NewWallet() Wallet {
	//balance := make(map[string]float64, 1)
	//balance[TOKENUNIT] = 0

	//balance[TokenUNIT] = 1001

	return Wallet{
		Balance: 0,
		Unit:    "pingala",
	}
}

func (w *Wallet) Update(value float64) {
	w.mux.Lock()
	defer w.mux.Unlock()

	w.Balance = w.Balance + value

}

func (w *Wallet) Show() string {
	return "Wallet : \n" + fmt.Sprintf("%f", w.Balance) + w.Unit
	//return showStr
}

//
//func BuildWallet(chains []p2.Blockchain, walletKey string) Wallet {
//
//	//walletKey := sha3.Sum256()
//
//	//btx := NewBorrowingTransactions()
//	//wallet := NewWallet()
//	//loop over the blockchain[00 of chains
//	var i int32
//	if len(chains) > 0 {
//		for i =1; i<= chains[0].Length; i++ {
//			blks, found := chains[0].Get(i)
//			if found && len(blks) > 0 {
//				blk := b.Block(blks[0])
//				mpt := p1.MerklePatriciaTrie(blk.Value)
//				keyValuePairs := mpt.GetAllKeyValuePairs()  //key - txid value - txJson
//				//loop over all key valye pairs and collect borrowing txs
//				for _, txjson := range keyValuePairs {
//					tx := JsonToTransaction([]byte(txjson))
//					if tx.To ==
//
//				}
//			}
//		}
//	}
//
//	return btx
//
//}
