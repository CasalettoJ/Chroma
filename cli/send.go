package cli

import (
	"fmt"
	"os"

	"github.com/casalettoj/chroma/blockchain"
	"github.com/casalettoj/chroma/wallet"
)

// send creates a TX and CoinbaseTX and mines a new transaction
func send(from, to string, amount int) {
	if amount <= 0 {
		fmt.Println("Invalid amount.")
		os.Exit(1)
	}

	bc := blockchain.OpenBlockchain()
	defer bc.DB.Close()

	wallets := wallet.OpenWallets()

	newTx := blockchain.NewTransaction(bc, wallets, to, from, amount)
	coinbaseTx := blockchain.NewCoinbaseTx(from, "")
	Txs := []*blockchain.Transaction{coinbaseTx, newTx}
	newBlock := bc.MineBlock(Txs)
	blockchain.UpdateUTXOs(bc, newBlock)
	fmt.Printf("Sent %d to %s.\n", amount, to)
}
