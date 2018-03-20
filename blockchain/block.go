package blockchain

import (
	"time"
)

// Block is the collection of data and headers for a single entry in the blockchain
type Block struct {
	Timestamp int64
	Data      []byte
	PrevHash  []byte
	Hash      []byte
	Nonce     int
}

// NewBlock creates a new block
func NewBlock(data string, prevHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

// GenerateGenesisBlock creates a new genesis block for a new blockchain with a special message
func GenerateGenesisBlock() *Block {
	return NewBlock("09 F9 11 02 9D 74 E3 5B D8 41 56 C5 63 56 88 C0", []byte{})
}
