package pow

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
	"math/rand"
	"strings"
	"time"

	//"../p1"
	"../data_structure/mpt"
	//"../p2"
	"../blockchain"
	//"../p2/block"
	"../block"
	//"../p3/data"
	"../sync_blockchain"
)

// func to InitializeNonce
func InitializeNonce(length int) string {
	bytes := make([]byte, length)
	source := rand.NewSource(time.Now().UnixNano())

	for i := range bytes {
		bytes[i] = byte(source.Int63()) //
	}

	nonce := hex.EncodeToString(bytes)
	return nonce
}

//check if POW is satisfied - return true or false
func POW(parentHash string, nonce string, mptRootHash string, difficulty int) bool {

	bytes := sha3.Sum256([]byte(parentHash + nonce + mptRootHash))
	work := hex.EncodeToString(bytes[:])
	proof := string(work[:difficulty])

	against := ""
	for i := 0; i < difficulty; i++ {
		against += "0"
	}
	//against := "0000000"

	//log.Println("in POW Comparing proof - against ", proof, " : ", against)
	if strings.Compare(proof, against) == 0 {
		//log.Println("WoW in POW Comparing proof - against ", proof, " : ", against)
		return true
	}
	return false
}

// FindNonce is used once - to insert 1st block
func FindNonce(parentHash string, mpt *mpt.MerklePatriciaTrie, difficulty int) string {

	// y = SHA3(parentHash + nonce + mptRootHash)

	for {
		nonce := InitializeNonce(8) //NextNonce(Nonce)
		if POW(parentHash, nonce, mpt.Root, difficulty) {
			return nonce
		}
	}
}

// GetCanonicalChains func returns slice of canonical blockchains
func GetCanonicalChains(SBC *sync_blockchain.SyncBlockChain) []blockchain.Blockchain {
	maxHeight := SBC.GetLength() // - 6 //
	blocksAtMaxHeight, _ := SBC.Get(maxHeight)

	canonicalChains := make([]blockchain.Blockchain, len(blocksAtMaxHeight))
	for i := range blocksAtMaxHeight {
		bc := blockchain.Blockchain{}
		bc.Initial()
		canonicalChains[i] = bc
	}

	for lastblocks := range blocksAtMaxHeight {
		canonicalChains[lastblocks].UnsafeInsert(blocksAtMaxHeight[lastblocks])

	}

	//if len(canonicalChains) >= 1 {
	for _, chain := range canonicalChains {
		for height := maxHeight - 1; height > 0; height-- {
			existingChildBlocks, _ := chain.Get(height + 1)
			potentialParentBlocks, _ := SBC.Get(height)
			for _, potentialParentBlock := range potentialParentBlocks {
				if block.Block(existingChildBlocks[0]).Header.ParentHash == block.Block(potentialParentBlock).Header.Hash {
					chain.UnsafeInsert(potentialParentBlock)
				}
			}

		}
	}

	return canonicalChains
}

///// pow ////
