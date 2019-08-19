package uri_routing

import (
	//"../p1"
	mptp "../data_structure/mpt"
	//"../p2"
	"../blockchain"
	//b "../p2/block"
	b "../block"
	"../consensus/pow"
	//"../p4"
	//p5 "../tokens"
	balance "../balance_book"
	"log"
	"strconv"
)

var BalanceBook balance.BalanceBook

func InitBalanceBook() {
	//currency //p5
	BalanceBook = balance.NewBalanceBook()
}

//func DefaultTokens() {
//	pubKeyStr := ID.GetMyPublicIdentity().PublicKey.N.String()
//	value := Wallet.Balance[p5.TOKENUNIT]
//	BalanceBook.UpdateABalanceInBook(pubKeyStr, value) //todo p5 todo p5
//}

// func to generate transactionMPT
func GenerateTransactionsMPT() mptp.MerklePatriciaTrie {
	mpt := mptp.MerklePatriciaTrie{}
	mpt.Initial()

	random := 5 //int((time.Now().UnixNano() / 100000 % 5))

	txs := TxPool.ReadFromTransactionPool(random)
	for _, tx := range txs {

		//todo - check if transaction valid -- need 2 things -  balancebook to see if enough balance and
		//
		//			todo									 -  list of all tx id map[string- txid]bool-false

		chains := pow.GetCanonicalChains(&SBC)

		//check if tx id already present on chain
		txl := BuildTransactionsList(chains[0])
		if _, ok := txl[tx.Id]; !ok { // txansaction not on canonical chain

			log.Println("In GenerateTransactionsMPT - tx not in canonical chain - so moving on ... :-)")

			bb := balance.NewBalanceBook()
			bb.BuildBalanceBook(chains[0], 2)

			// check if available balance is enough //todo check
			senderKeyForBook := balance.GetKeyForBook(tx.From.PublicKey)
			available, err := bb.Book.Get(senderKeyForBook)
			if err == nil {
				availBal, _ := strconv.ParseFloat(available, 64)

				amountPromised := bb.CheckAmountPromisedByOne(tx.From)

				if availBal >= tx.Tokens+tx.Fees+amountPromised && tx.Tokens >= 0 && tx.Fees >= 0 /* + amount promised */ {
					log.Println("In GenerateTransactionsMPT - Enough Balance available - so moving on ... :-)")
					/// code goes here
					mpt.Insert(tx.Id, tx.TransactionToJson())

				} else if tx.TxType == "promise" && tx.Tokens < 0 &&
					availBal >= tx.Fees+amountPromised && (amountPromised >= -tx.Tokens) && tx.Fees >= 0 {
					//take back promise //check to see if the promise of sum = or > than the amount in take back promise
					log.Println("In GenerateTransactionsMPT - Enough Promised to take back some/all promise - so moving on ... :-} ")
					mpt.Insert(tx.Id, tx.TransactionToJson())

				}

			} else if tx.TxType == "start" {
				log.Println("In GenerateTransactionsMPT - Start Transaction for miner - so moving on ... :-)")
				/// code goes here
				mpt.Insert(tx.Id, tx.TransactionToJson())
			}
		}

		TxPool.DeleteFromTransactionPool(tx.Id) //delete from TransactionPool
	}

	return mpt
}

//func MarkTxInTxPoolAsUsed(mpt p1.MerklePatriciaTrie) {
//	usedTxPool := mpt.GetAllKeyValuePairs()
//	for _, txJson := range usedTxPool {
//		tx := p5.DecodeToTransaction([]byte(txJson))
//		TxPool.Pool[tx] = true
//	}
//}

func BuildTransactionsList(chain blockchain.Blockchain) map[string]bool {

	txl := make(map[string]bool)
	//loop over the blockchain[0] of chains
	var i int32

	for i = 2; i <= chain.Length; i++ {
		blks, found := chain.Get(i)
		if found && len(blks) > 0 {
			blk := b.Block(blks[0])
			mpt := mptp.MerklePatriciaTrie(blk.Value)
			keyValuePairs := mpt.GetAllKeyValuePairs() //key - txid value - txJson
			//loop over all key valye pairs and collect borrowing txs
			for txId, _ := range keyValuePairs {
				txl[txId] = false
			}
		}
	}

	return txl

}
