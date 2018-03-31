package blockchain

import (
	"bytes"

	wallet "github.com/casalettoj/chroma/wallet"
)

// TxInput represents a transaction input
type TxInput struct {
	TxID      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

// ScriptSigCheck returns whether the input was initiated by a given address
func (txi *TxInput) ScriptSigCheck(pubKeyHash []byte) bool {
	hashLock := wallet.HashPublicKey(txi.PubKey)
	return bytes.Compare(hashLock, pubKeyHash) == 0
}
