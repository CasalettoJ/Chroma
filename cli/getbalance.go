package cli

import (
	"fmt"

	"github.com/casalettoj/chroma/blockchain"
)

// getBalance prints the balance of a given address to the console
func getBalance(address string) {
	total := 0
	bc := blockchain.OpenBlockchain()
	defer bc.DB.Close()
	bc.GetBalance(address)

	fmt.Printf("Balance of '%s': %d\n", address, total)
}
