package cli

import (
	"fmt"
	"os"

	"github.com/casalettoj/chroma/wallet"

	chroma "github.com/casalettoj/chroma/blockchain"
)

// send creates a TX and CoinbaseTX and mines a new transaction
func send(from, to string, amount int) {
	if amount <= 0 {
		fmt.Println("Invalid amount.")
		os.Exit(1)
	}

	bc := chroma.OpenBlockchain()
	defer bc.DB.Close()

	wallets := wallet.OpenWallets()

	newTx := chroma.NewTransaction(bc, wallets, to, from, amount)
	coinbaseTx := chroma.NewCoinbaseTx(from, "")
	Txs := []*chroma.Transaction{coinbaseTx, newTx}
	newBlock := bc.MineBlock(Txs)
	chroma.UpdateUTXOs(bc, newBlock)
	fmt.Printf("Sent %d to %s.\n", amount, to)
}
