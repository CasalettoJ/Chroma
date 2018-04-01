package cli

import (
	"fmt"

	chroma "github.com/casalettoj/chroma/blockchain"
)

// printChain iterates through the chain and prints the data of each
func printChain() {
	bc := chroma.OpenBlockchain()
	defer bc.DB.Close()
	bci := bc.Iterator()

	for {
		fmt.Println()
		block := bci.Next()
		blockStrings := block.String()
		for _, s := range blockStrings {
			fmt.Printf("%s", s)
		}
		fmt.Println()
		if bci.IsGenesisBlock() {
			break
		}
	}
}
