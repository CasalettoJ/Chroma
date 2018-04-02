package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

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

func (tx *Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("___TX %x: (Vin: %d, Vout: %d)___\n", tx.ID, len(tx.Vin), len(tx.Vout)))
	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("Input %d:\nTxID: %x\nVout: %d\nSig: %x\nPubKey: %x\n", i, input.TxID, input.Vout, input.Signature, input.PubKey))

	}
	lines = append(lines, fmt.Sprintln())
	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("Output %d:\nValue: %d\nPubKeyHash: %x\n", i, output.Value, output.PubKeyHash))
	}
	lines = append(lines, fmt.Sprintf("___\n\n"))
	return strings.Join(lines, "")
}

// IsCoinbaseTx identifies if the tx given is a coinbase tx based on inputs.
func (tx *Transaction) IsCoinbaseTx() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxID) == 0 && tx.Vin[0].Vout == -1
}

// Hash returns a sha256 hash of serialized Tx data
func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

// Serialize returns a byte slice representation of the tx
func (tx *Transaction) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	util.CheckAnxiety(encoder.Encode(tx))
	return buffer.Bytes()
}

// TrimmedCopy returns a copy of the transaction with inputs stripped of their PubKey and Signature fields.
func (tx *Transaction) TrimmedCopy() Transaction {
	var vin []TxInput
	var vout []TxOutput

	for _, in := range tx.Vin {
		vin = append(vin, TxInput{TxID: in.TxID, Vout: in.Vout, Signature: nil, PubKey: nil})
	}
	for _, out := range tx.Vout {
		vout = append(vout, TxOutput{Value: out.Value, PubKeyHash: out.PubKeyHash})
	}

	return Transaction{Vin: vin, Vout: vout, ID: tx.ID}
}

// Sign signs a transaction with the private key of a wallet and a hash of the transaction
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	if tx.IsCoinbaseTx() {
		return
	}
	txCopy := tx.TrimmedCopy()
	for _, vin := range tx.Vin {
		if prevTxs[hex.EncodeToString(vin.TxID)].ID == nil {
			log.Panic("ERROR: Invalid Previous Tx List")
		}
	}
	for i, in := range txCopy.Vin {
		prevTx := prevTxs[hex.EncodeToString(in.TxID)]
		txCopy.Vin[i].Signature = nil
		txCopy.Vin[i].PubKey = prevTx.Vout[in.Vout].PubKeyHash // Set the pubkey in order to hash accurately w/ prev output
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[i].PubKey = nil // Don't need that anymore

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.ID)
		util.CheckAnxiety(err)

		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[i].Signature = signature
	}
}

// Verify checks the given transaction and verifies every input's signature
func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for i, in := range tx.Vin {
		// Generate a hash for the txCopy using the referenced output's public key hash as the pubkey.
		prevTx := prevTxs[hex.EncodeToString(in.TxID)]
		txCopy.Vin[i].Signature = nil
		txCopy.Vin[i].PubKey = prevTx.Vout[in.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[i].PubKey = nil

		// Take the input's signature and split it into it's pair
		r := big.Int{}
		s := big.Int{}
		sigLength := len(in.Signature)
		r.SetBytes(in.Signature[:sigLength/2])
		s.SetBytes(in.Signature[sigLength/2:])

		// Take the public key of the input and split it, then create a new raw public key from it
		x := big.Int{}
		y := big.Int{}
		keyLength := len(in.PubKey)
		x.SetBytes(in.PubKey[:keyLength/2])
		y.SetBytes(in.PubKey[keyLength/2:])
		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

		// Check if the raw public key
		if !ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) {
			return false
		}
	}

	return true
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
	newTx.ID = newTx.Hash()
	bc.SignTransaction(&newTx, fromWallet.PrivateKey)

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
	tx.ID = tx.Hash()
	return &tx
}
