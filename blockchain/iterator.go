package blockchain

import (
	conf "github.com/casalettoj/chroma/constants"
	util "github.com/casalettoj/chroma/utils"
	bolt "github.com/coreos/bbolt"
)

// Iterator is the same as blockchain but it's got a different name.
type Iterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

// Next returns the current hash and decrements the current block hash to its previous
func (i *Iterator) Next() *Block {
	block := i.Peek()
	i.CurrentHash = block.PrevHash
	return block
}

// IsGenesisBlock returns whether the current hash points to the genesis block
func (i *Iterator) IsGenesisBlock() bool {
	return i.CurrentHash == nil
}

// Peek returns the block at the current hash of the iterator
func (i *Iterator) Peek() *Block {
	var block *Block
	util.CheckAnxiety(i.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DBblocksbucket))
		encodedBlock := bucket.Get(i.CurrentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	}))
	return block
}
