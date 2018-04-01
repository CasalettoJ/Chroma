package cli

import (
	"fmt"

	chroma "github.com/casalettoj/chroma/blockchain"
)

// PrintChain iterates through the chain and prints the data of each
func PrintChain() {
	bc := chroma.OpenBlockchain()
	defer bc.DB.Close()
	bci := bc.Iterator()

	for {
		fmt.Println()
		fmt.Println("===============================================")
		block := bci.Next()
		if block.PrevHash != nil {
			fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		}
		fmt.Printf("Tx Hash: %x\n", block.HashTransactions())
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println("Transactions:")
		for i, tx := range block.Transactions {
			fmt.Println()
			if i != 0 {
				fmt.Println("____________")
			}
			fmt.Printf("TX %d: %d inputs %d outputs\n", i, len(tx.Vin), len(tx.Vout))
			fmt.Println()
			for j, input := range tx.Vin {
				fmt.Printf("Input %d:\n%+v\n", j, input)
			}
			fmt.Println()
			for j, output := range tx.Vout {
				fmt.Printf("Output %d:\n%+v\n", j, output)
			}
		}
		fmt.Println("===============================================")
		fmt.Println()
		if bci.IsGenesisBlock() {
			break
		}
	}
}
