package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
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

// Serialize returns a byte array serialization
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// DeserializeBlock deserializes a byte array into a Block struct
func DeserializeBlock(bbytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(bbytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
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
