package tokens

import (
	"encoding/json"
	"fmt"
	"log"
	//s "../p5security"
	"../client"
	s "../identity"
)

type TransactionBeat struct {
	Tx      Transaction
	FromPid s.PublicIdentity
	TxSig   []byte
	Hops    int
}

func NewTransactionBeat(tx Transaction, fromPid s.PublicIdentity, fromSig []byte) TransactionBeat {
	return TransactionBeat{
		Tx:      tx,
		TxSig:   fromSig,
		FromPid: fromPid,
		Hops:    1,
	}
}

func PrepareTransactionBeat(tx Transaction, cid client.ClientId) TransactionBeat {

	pid := cid.GetMyPublicIdentity()

	return TransactionBeat{
		Tx:      tx,
		TxSig:   tx.CreateTxSig(cid),
		FromPid: pid,
		Hops:    1,
	}
}

//EncodeToJson func encodes HeartBeatData to json byte array
func (data *TransactionBeat) EncodeToJsonByteArray() []byte {

	dataEncodedBytearray, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Here in err condition of EncodeToJsonByteArray of transactionBeat.go")
		return []byte("{}")
	}
	return dataEncodedBytearray
}

//EncodeToJson func encodes HeartBeatData to json string
func (data *TransactionBeat) EncodeToJson() string {

	dataEncodedBytearray, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Here in err condition of EncodeToJson of transactionBeat.go")
		return "{}" //empty heartbeat
	}
	return string(dataEncodedBytearray)
}

//DecodeToHeartBeatData func decodes json string to HeartBeatData
func DecodeToTransactionBeat(transactionBeatJson []byte) TransactionBeat {
	tb := TransactionBeat{}
	err := json.Unmarshal(transactionBeatJson, &tb)
	if err != nil {
		log.Println("Err in DecodeToTransactionBeat in transactionBeat.go - err : ", err)
		log.Println("Error transactionBeatJson : ", transactionBeatJson)
	} else {

	}
	return tb
}

func (data *TransactionBeat) VerifyTxSigInTxBeat() bool {

	return VerifyTxSig(data.Tx.From, data.Tx, data.TxSig)

}
