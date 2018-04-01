package blockchain

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58"

	conf "github.com/casalettoj/chroma/constants"
	util "github.com/casalettoj/chroma/utils"
	wallet "github.com/casalettoj/chroma/wallet"
)

// Transaction is a collection of inputs and outputs with its hashed data as an ID
type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

// IsCoinbaseTx identifies if the tx given is a coinbase tx based on inputs.
func (tx *Transaction) IsCoinbaseTx() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxID) == 0 && tx.Vin[0].Vout == -1
}

// SetID hashes the data in the tx and sets it to the ID
func (tx *Transaction) SetID() {
	var buffer bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&buffer)
	util.CheckAnxiety(encoder.Encode(tx))
	hash = sha256.Sum256(buffer.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) String() (lines []string) {
	lines = append(lines, fmt.Sprintf("___TX %x: (Vin: %d, Vout: %d)___\n", tx.ID, len(tx.Vin), len(tx.Vout)))
	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("Input %d:\nTxID: %x\nVout: %d\nSig: %x\nPubKey: %x\n", i, input.TxID, input.Vout, input.Signature, input.PubKey))

	}
	lines = append(lines, fmt.Sprintln())
	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("Output %d:\nValue: %d\nPubKeyHash: %x\n", i, output.Value, output.PubKeyHash))
	}
	lines = append(lines, fmt.Sprintf("___\n\n"))
	return
}

// NewTransaction returns a new transaction
func NewTransaction(bc *Blockchain, wallets *wallet.Wallets, to, from string, amount int) *Transaction {
	var vin []TxInput
	vout := []TxOutput{*NewUTXO(amount, to)}

	fromWallet := wallets.GetWallet(from)
	fromAddress := fromWallet.GetChromaAddress()
	fromPubKeyHash := base58.Decode(string(fromAddress))
	fromPubKeyHash = fromPubKeyHash[1 : len(fromPubKeyHash)-4]
	totalIn, usedTxOutputs := FindUTXOsForPayment(bc, fromPubKeyHash, amount)

	if totalIn < amount {
		log.Panic(fmt.Printf("Insufficient Funds: Found %d and needed at least %d\nTx Found: %+v\n", totalIn, amount, usedTxOutputs))
	}

	for txID, outputs := range usedTxOutputs {
		txIDBytes, err := hex.DecodeString(txID)
		util.CheckAnxiety(err)
		for _, outputIndex := range outputs {
			input := TxInput{Vout: outputIndex, Signature: nil, PubKey: fromWallet.PublicKey, TxID: txIDBytes}
			vin = append(vin, input)
		}
	}

	if totalIn > amount {
		change := *NewUTXO(totalIn-amount, from)
		vout = append(vout, change)
	}

	newTx := Transaction{Vin: vin, Vout: vout}
	newTx.SetID()

	return &newTx
}

// NewCoinbaseTx returns a special TX to be awarded for mining a block.
func NewCoinbaseTx(to, data string) *Transaction {
	// Fill pubkey with random data
	if data == "" {
		randomData := make([]byte, 20)
		_, err := rand.Read(randomData)
		util.CheckAnxiety(err)
		data = string(randomData)
	}
	txin := TxInput{TxID: []byte{}, Vout: -1, Signature: nil, PubKey: []byte(data)}
	txout := *NewUTXO(conf.TXcoinbaseaward, to)
	tx := Transaction{ID: nil, Vin: []TxInput{txin}, Vout: []TxOutput{txout}}
	tx.SetID()
	return &tx
}
