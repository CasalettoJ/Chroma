package cli

import (
	"fmt"

	"github.com/casalettoj/chroma/wallet"
)

// createNewWallet creates a new private/public key pair and adds it to the wallets file.
func createNewWallet() {
	wallets := wallet.OpenWallets()
	address := wallets.AddNewWallet()
	wallets.SaveWallets()
	fmt.Printf("New wallet created. Address: %s\n", address)
}
