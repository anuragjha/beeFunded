package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"log"
	//s "../p5security"
	s "../identity"
)

type ClientId struct {
	PrivateKey *rsa.PrivateKey `json:"privateKey"`
	PublicKey  *rsa.PublicKey  `json:"publicKey"`
	//HashForKey 	hash.Hash
	Label string `json:"label"`
}

func (id *ClientId) ClientIdToJsonByteArray() []byte {
	jsonByteArray, err := json.Marshal(id)
	if err != nil {
		log.Println("Cannot convert Clientidentity to json bytes, err - ", err)
	}
	return jsonByteArray
}

func JsonToClientId(jsonStr string) (ClientId, error) {

	cid := ClientId{}
	err := json.Unmarshal([]byte(jsonStr), &cid)
	if err != nil {
		log.Println("Cannot convert json to ClientIdentity, err - ", err)
		return ClientId{}, err
	}
	return cid, nil
}

func NewClientId(label string) ClientId {
	id := ClientId{}
	id.PrivateKey, id.PublicKey = s.GeneratePubPrivKeyPair()
	//id.HashForKey = GenerateHashForKey(label)
	id.Label = label

	return id
}

func ExistingClientId(priv *rsa.PrivateKey, pub *rsa.PublicKey, label string) ClientId {
	id := ClientId{}
	id.PrivateKey = priv
	id.PublicKey = pub
	//id.HashForKey = GenerateHashForKey(label)
	id.Label = label

	return id
}

func (id *ClientId) GenSignature(message []byte) []byte {
	//digest := GenDigest(id.HashForKey, message)//id.HashForKey.Sum(message)
	//signature, err := rsa.SignPSS(rand.Reader, id.privateKey, crypto.SHA256, digest, nil)
	hashMsg := sha3.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, id.PrivateKey, crypto.SHA256, hashMsg[:])
	if err != nil {
		log.Fatal("Err in generating signature : err : ", err)
	}
	return signature
}

func (id *ClientId) GetMyPublicIdentity() s.PublicIdentity {
	pid := s.PublicIdentity{}
	pid.PublicKey = id.PublicKey
	//pid.HashForKey = id.HashForKey
	pid.Label = id.Label

	return pid
}
