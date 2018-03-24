package blockchain

import (
	"fmt"
	"os"

	bolt "github.com/coreos/bbolt"
)

// Blockchain is the entire structure of sequential blocks
type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

// Iterator give iterator
func (bc *Blockchain) Iterator() *Iterator {
	iterator := &Iterator{bc.Tip, bc.DB}
	return iterator
}

// OpenBlockchain opens a preexisting blockchain and returns Tip and DB
func OpenBlockchain() *Blockchain {
	if !DoesDBExist() {
		fmt.Println("No existing Chroma chain.  Create DB first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(DBdbfile, 0600, nil)
	CheckAnxiety(err)

	CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DBblocksbucket))
		tip = bucket.Get([]byte(DBlasthash))
		return nil
	}))
	bc := &Blockchain{DB: db, Tip: tip}
	return bc
}

// CreateBlockchain establishes a blockchain with a genesis block
func CreateBlockchain(address string) *Blockchain {
	if DoesDBExist() {
		fmt.Println("Chroma chain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(DBdbfile, 0600, nil)
	CheckAnxiety(err)

	CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		genesisBlock := GenerateGenesisBlock(NewCoinbaseTx(Message, address))
		bucket, err := tx.CreateBucket([]byte(DBblocksbucket))
		CheckAnxiety(err)
		CheckAnxiety(bucket.Put(genesisBlock.Hash, genesisBlock.Serialize()))
		CheckAnxiety(bucket.Put([]byte(DBlasthash), genesisBlock.Hash))
		tip = genesisBlock.Hash
		return nil
	}))
	return &Blockchain{tip, db}
}
