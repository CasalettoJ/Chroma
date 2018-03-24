package blockchain

import (
	bolt "github.com/coreos/bbolt"
)

// Iterator is the same as blockchain but it's got a different name.
type Iterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

// Next returns the next block decrementing from the Tip
func (i *Iterator) Next() *Block {
	var block *Block
	CheckAnxiety(i.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DBblocksbucket))
		encodedBlock := bucket.Get(i.CurrentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	}))
	i.CurrentHash = block.PrevHash
	return block
}
