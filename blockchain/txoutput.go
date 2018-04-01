package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/btcsuite/btcutil/base58"
	util "github.com/casalettoj/chroma/utils"
)

// TxOutput represents a transaction output
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// LockTxO locks a TxO to a specific public key hash (1:len-4 bytes)
func (txo *TxOutput) LockTxO(address []byte) {
	pubKeyHash := base58.Decode(string(address[:]))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4] //Start at 1 for the version byte, stop 4 early for the checksum bytes
	txo.PubKeyHash = pubKeyHash
}

// Unlockable returns whether the output can be unlocked by a given address
func (txo *TxOutput) Unlockable(pubKeyHash []byte) bool {
	return bytes.Compare(txo.PubKeyHash, pubKeyHash) == 0
}

// NewUTXO creates a new transaction for a value and pubkeyhash string
func NewUTXO(value int, address string) (utxo *TxOutput) {
	utxo = &TxOutput{Value: value, PubKeyHash: nil}
	utxo.LockTxO([]byte(address))
	return
}

// TxOutputs is a structure containing an array of TXOs
type TxOutputs struct {
	Outputs []TxOutput
}

// Serialize returns a byte array serialization
func (txos *TxOutputs) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	util.CheckAnxiety(encoder.Encode(txos))
	return result.Bytes()
}

// DeserializeTxOutputs deserializes a byte array into a TxOutputs struct
func DeserializeTxOutputs(bbytes []byte) *TxOutputs {
	var txOutputs TxOutputs
	decoder := gob.NewDecoder(bytes.NewReader(bbytes))
	util.CheckAnxiety(decoder.Decode(&txOutputs))
	return &txOutputs
}
