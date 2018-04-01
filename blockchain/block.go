package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"

	util "github.com/casalettoj/chroma/utils"
)

// Block is the collection of data and headers for a single entry in the blockchain
type Block struct {
	Timestamp    int64
	Transactions []*Transaction
	PrevHash     []byte
	Hash         []byte
	Nonce        int
}

// Serialize returns a byte array serialization
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	util.CheckAnxiety(encoder.Encode(b))
	return result.Bytes()
}

// HashTransactions returns a []byte hash representation of all txIDs in a Tx
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

func (b *Block) String() (lines []string) {
	lines = append(lines, fmt.Sprintf("====BLOCK %x====\n", b.Hash))
	if b.PrevHash != nil {
		lines = append(lines, fmt.Sprintf("Prev. hash: %x\n", b.PrevHash))
	}
	lines = append(lines, fmt.Sprintf("Tx Hash: %x\n", b.HashTransactions()))
	lines = append(lines, fmt.Sprintln("Transactions:"))
	for _, tx := range b.Transactions {
		lines = append(lines, tx.String()...)
	}
	return
}

// DeserializeBlock deserializes a byte array into a Block struct
func DeserializeBlock(bbytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(bbytes))
	util.CheckAnxiety(decoder.Decode(&block))
	return &block
}

// NewBlock creates a new block
func NewBlock(transactions []*Transaction, prevHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

// GenerateGenesisBlock creates a new genesis block for a new blockchain with a special message
func GenerateGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
