package main

import (
	"fmt"
	"strconv"

	"github.com/chromanetwork/chroma/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()

	bc.AddBlock("Block One")
	bc.AddBlock("Block Two")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.IsValid()))
		fmt.Println()
	}
}
