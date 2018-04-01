package cli

import (
	"fmt"
	"os"

	chroma "github.com/casalettoj/chroma/blockchain"
)

// Send creates a TX and CoinbaseTX and mines a new transaction
func Send(from, to string, amount int) {
	if amount <= 0 {
		fmt.Println("Invalid amount.")
		os.Exit(1)
	}

	bc := chroma.OpenBlockchain()
	defer bc.DB.Close()

	newTx := chroma.NewTransaction(bc, to, from, amount)
	coinbaseTx := chroma.NewCoinbaseTx(from, "")
	Txs := []*chroma.Transaction{coinbaseTx, newTx}
	newBlock := bc.MineBlock(Txs)
	chroma.UpdateUTXOs(bc, newBlock)
	fmt.Printf("Sent %d to %s.", amount, to)
}
