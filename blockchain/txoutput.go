package blockchain

import (
	"bytes"
	"encoding/gob"
)

// TxOutput represents a transaction output
type TxOutput struct {
	Value        int
	ScriptPubKey string
}

// Unlockable returns whether the output can be unlocked by a given address
func (txo *TxOutput) Unlockable(address string) bool {
	return txo.ScriptPubKey == address
}

// TxOutputs is a structure containing an array of TXOs
type TxOutputs struct {
	Outputs []TxOutput
}

// Serialize returns a byte array serialization
func (txos *TxOutputs) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	CheckAnxiety(encoder.Encode(txos))
	return result.Bytes()
}

// DeserializeTxOutputs deserializes a byte array into a Block struct
func DeserializeTxOutputs(bbytes []byte) *TxOutputs {
	var txOutputs TxOutputs
	decoder := gob.NewDecoder(bytes.NewReader(bbytes))
	CheckAnxiety(decoder.Decode(&txOutputs))
	return &txOutputs
}
