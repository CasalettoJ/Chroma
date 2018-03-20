package blockchain

// Blockchain is the entire structure of sequential blocks
type Blockchain struct {
	Blocks []*Block
}

// AddBlock adds a new block with given data to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevHash := bc.Blocks[len(bc.Blocks)-1].Hash
	newBlock := NewBlock(data, prevHash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// NewBlockchain establishes a blockchain with a genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenerateGenesisBlock()}}
}
