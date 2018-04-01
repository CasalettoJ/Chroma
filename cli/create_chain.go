package cli

import (
	"fmt"

	chroma "github.com/casalettoj/chroma/blockchain"
)

// CreateBlockchain creates a blockchain db
func CreateBlockchain(address string) {
	bc := chroma.CreateBlockchain(address)
	defer bc.DB.Close()
	chroma.ReindexUTXOs(bc)
	fmt.Println("CHROMA chain created")
}
