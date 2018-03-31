package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"os"

	conf "github.com/casalettoj/chroma/constants"
	util "github.com/casalettoj/chroma/utils"
)

// Wallets holds private keys mapped by
type Wallets struct {
	Wallets map[string]*Wallet
}

// AddNewWallet creates a new private/public key pair and adds it to the wallet.
func (ws *Wallets) AddNewWallet() string {
	wallet := NewWallet()
	address := string(wallet.GetChromaAddress())
	ws.Wallets[address] = wallet
	return address
}

// GetWallet returns the wallet stored at the address specified
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// GetAddresses returns a string array of address keys in the Wallets collection
func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for k := range ws.Wallets {
		addresses = append(addresses, k)
	}
	return addresses
}

// SaveWallets saves the wallets data to a file
func (ws Wallets) SaveWallets() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	util.CheckAnxiety(encoder.Encode(ws))
	util.CheckAnxiety(ioutil.WriteFile(conf.WalletFile, content.Bytes(), 0644))
}

// OpenWallets creates a new wallets file if none exists otherwise loads from file and returns result
func OpenWallets() *Wallets {
	_, err := os.Stat(conf.WalletFile)
	if os.IsNotExist(err) {
		wallets := Wallets{}
		wallets.Wallets = make(map[string]*Wallet)
		wallets.SaveWallets()
		return &wallets
	}
	content, err := ioutil.ReadFile(conf.WalletFile)
	util.CheckAnxiety(err)
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	util.CheckAnxiety(decoder.Decode(&wallets))
	return &wallets
}
