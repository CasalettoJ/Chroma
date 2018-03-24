package blockchain

import (
	bolt "github.com/coreos/bbolt"
)

// Blockchain is the entire structure of sequential blocks
type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

// AddBlock adds a new block with given data to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	var lasthash []byte
	CheckAnxiety(bc.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DBblocksbucket))
		lasthash = bucket.Get([]byte(DBlasthash))
		return nil
	}))
	newBlock := NewBlock(data, lasthash)
	CheckAnxiety(bc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DBblocksbucket))
		CheckAnxiety(bucket.Put(newBlock.Hash, newBlock.Serialize()))
		CheckAnxiety(bucket.Put([]byte(DBlasthash), newBlock.Hash))
		bc.Tip = newBlock.Hash
		return nil
	}))
}

// Iterator give iterator
func (bc *Blockchain) Iterator() *Iterator {
	iterator := &Iterator{bc.Tip, bc.DB}
	return iterator
}

// NewBlockchain establishes a blockchain with a genesis block
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(DBdbfile, 0600, nil)
	CheckAnxiety(err)

	CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DBblocksbucket))
		if bucket != nil {
			tip = bucket.Get([]byte(DBlasthash))
		} else {
			genesisBlock := GenerateGenesisBlock()
			bucket, err := tx.CreateBucket([]byte(DBblocksbucket))
			CheckAnxiety(err)
			CheckAnxiety(bucket.Put(genesisBlock.Hash, genesisBlock.Serialize()))
			CheckAnxiety(bucket.Put([]byte(DBlasthash), genesisBlock.Hash))
			tip = genesisBlock.Hash
		}

		return nil
	}))
	return &Blockchain{tip, db}
}
