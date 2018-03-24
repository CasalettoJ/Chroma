package blockchain

// TxInput represents a transaction input
type TxInput struct {
	TxID      []byte
	Vout      int
	ScriptSig string
}

// InitiatedBy returns whether the input was initiated by a given address
func (txi *TxInput) InitiatedBy(address string) bool {
	return txi.ScriptSig == address
}
