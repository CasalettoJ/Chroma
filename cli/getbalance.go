package cli

import (
	"fmt"

	chroma "github.com/casalettoj/chroma/blockchain"
)

// GetBalance prints the balance of a given address to the console
func GetBalance(address string) {
	total := 0
	bc := chroma.OpenBlockchain()
	defer bc.DB.Close()
	bc.GetBalance(address)

	fmt.Printf("Balance of '%s': %d\n", address, total)
}
