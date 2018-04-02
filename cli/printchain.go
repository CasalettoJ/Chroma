package cli

import (
	"fmt"

	"github.com/casalettoj/chroma/blockchain"
)

// printChain iterates through the chain and prints the data of each
func printChain() {
	bc := blockchain.OpenBlockchain()
	defer bc.DB.Close()
	bci := bc.Iterator()

	for {
		fmt.Println()
		block := bci.Next()
		fmt.Println(block)
		fmt.Println()
		if bci.IsGenesisBlock() {
			break
		}
	}
}
