package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

// TxOutput represents a transaction output
type TxOutput struct {
	Value        int
	ScriptPubKey string
}

// TxInput represents a transaction input
type TxInput struct {
	TxID      []byte
	Vout      int
	ScriptSig string
}

// Transaction is a collection of inputs and outputs with its hashed data as an ID
type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

// SetID hashes the data in the tx and sets it to the ID
func (tx *Transaction) SetID() {
	var buffer bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&buffer)
	CheckAnxiety(encoder.Encode(tx))
	hash = sha256.Sum256(buffer.Bytes())
	tx.ID = hash[:]
}

// NewCoinbaseTx returns a special TX to be awarded for mining a block.
func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coinbase award for %s", to)
	}
	txin := TxInput{TxID: []byte{}, Vout: -1, ScriptSig: data}
	txout := TxOutput{Value: TXcoinbaseaward, ScriptPubKey: to}
	tx := Transaction{ID: nil, Vin: []TxInput{txin}, Vout: []TxOutput{txout}}
	tx.SetID()
	return &tx
}
