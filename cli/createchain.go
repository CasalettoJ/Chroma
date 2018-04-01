package cli

import (
	"fmt"

	chroma "github.com/casalettoj/chroma/blockchain"
)

// createBlockchain creates a blockchain db
func createBlockchain(address string) {
	bc := chroma.CreateBlockchain(address)
	defer bc.DB.Close()
	chroma.ReindexUTXOs(bc)
	fmt.Println("CHROMA chain created")
}
