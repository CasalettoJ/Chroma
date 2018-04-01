package cli

import (
	"fmt"

	"github.com/casalettoj/chroma/blockchain"
	"github.com/casalettoj/chroma/wallet"
)

// PrintWallets prints the address of every wallet in the wallet file.
func PrintWallets() {
	bc := blockchain.OpenBlockchain()
	wallets := wallet.OpenWallets()
	fmt.Println("Wallet Addresses:")
	for address := range wallets.Wallets {
		balance := bc.GetBalance(address)
		fmt.Printf("%s %d\n", address, balance)
	}
}
