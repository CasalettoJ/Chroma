package cli

import (
	"fmt"

	"github.com/casalettoj/chroma/blockchain"
)

// createBlockchain creates a blockchain db
func createBlockchain(address string) {
	bc := blockchain.CreateBlockchain(address)
	defer bc.DB.Close()
	blockchain.ReindexUTXOs(bc)
	fmt.Println("CHROMA chain created")
}
