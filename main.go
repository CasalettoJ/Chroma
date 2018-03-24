package main

import (
	"github.com/casalettoj/chroma/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer bc.DB.Close()

	cli := &blockchain.CLI{BC: bc}
	cli.Run()
}
