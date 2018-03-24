package blockchain

// TxOutput represents a transaction output
type TxOutput struct {
	Value        int
	ScriptPubKey string
}

// Unlockable returns whether the output can be unlocked by a given address
func (txo *TxOutput) Unlockable(address string) bool {
	return txo.ScriptPubKey == address
}
