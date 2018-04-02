package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/btcsuite/btcutil/base58"

	conf "github.com/casalettoj/chroma/constants"
	util "github.com/casalettoj/chroma/utils"
	bolt "github.com/coreos/bbolt"
)

// Blockchain is the entire structure of sequential blocks
type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

// GetBalance returns the balance of the address given for the current bc
func (bc *Blockchain) GetBalance(address string) int {
	total := 0
	pubKeyHash := base58.Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	UTXOs := GetUTXOsForAddress(bc, pubKeyHash)

	for _, UTXO := range UTXOs {
		total += UTXO.Value
	}
	return total
}

// FindTransaction traverses the blockchain looking for the TX matching ID
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Println(len(block.Transactions))
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		if bci.IsGenesisBlock() {
			break
		}
	}
	return Transaction{}, errors.New("tx Not found in chain")
}

// SignTransaction signs a transaction with a private key
func (bc *Blockchain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey) {
	prevTxs := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		prevTx, err := bc.FindTransaction(vin.TxID)
		util.CheckAnxiety(err)
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	tx.Sign(privateKey, prevTxs)
}

// VerifyTransaction verifies the signatures of a transaction's inputs
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbaseTx() {
		return true
	}
	prevTxs := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		prevTx, err := bc.FindTransaction(vin.TxID)
		if err != nil {
			return false
		}
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	verified := tx.Verify(prevTxs)
	return verified
}

// Iterator give iterator
func (bc *Blockchain) Iterator() *Iterator {
	iterator := &Iterator{bc.Tip, bc.DB}
	return iterator
}

// MineBlock mines a block with the given transactions
func (bc *Blockchain) MineBlock(Txs []*Transaction) *Block {
	var lastHash []byte

	for _, tx := range Txs {
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid Tx in Block")
		}
	}

	util.CheckAnxiety(bc.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DBblocksbucket))
		lastHash = bucket.Get([]byte(conf.DBlasthash))
		return nil
	}))

	newBlock := NewBlock(Txs, lastHash)

	util.CheckAnxiety(bc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DBblocksbucket))
		util.CheckAnxiety(bucket.Put(newBlock.Hash, newBlock.Serialize()))
		util.CheckAnxiety(bucket.Put([]byte(conf.DBlasthash), newBlock.Hash))
		bc.Tip = newBlock.Hash
		return nil
	}))

	return newBlock
}

// GetUTXOs gets all UTXOs in the blockchain
func (bc *Blockchain) GetUTXOs() map[string]TxOutputs {
	UTXOs := make(map[string]TxOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()
		// For every transaction in the block...
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
			// For every output in the transaction...
			for outIndex, out := range tx.Vout {
				spent := false
				// If this transaction has spent outputs..
				if spentTXOs[txID] != nil {
					// For every output index spent...
					for _, spentOut := range spentTXOs[txID] {
						// If this output index is the same as the index of the output being checked then it has been spent
						if outIndex == spentOut {
							spent = true
						}
					}
				}
				// If the output hasn't been spent and it can be unlocked by the address, add tx to the UTXset.
				if !spent {
					outputs := UTXOs[txID]
					outputs.Outputs = append(outputs.Outputs, out)
					UTXOs[txID] = outputs
				}
			}
			// If the transaction isn't a coinbase TX, then...
			if !tx.IsCoinbaseTx() {
				// For every input in the transaction...
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.TxID)
					// add the input's output index to the list of spent TXOs for the tx its referenced UTXO was created in.
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}
		if bci.IsGenesisBlock() {
			break
		}
	}

	return UTXOs
}

// OpenBlockchain opens a preexisting blockchain and returns Tip and DB
func OpenBlockchain() *Blockchain {
	if !util.DoesDBExist() {
		fmt.Println("No existing Chroma chain.  Create DB first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(conf.DBdbfile, 0600, nil)
	util.CheckAnxiety(err)

	util.CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DBblocksbucket))
		tip = bucket.Get([]byte(conf.DBlasthash))
		return nil
	}))
	bc := &Blockchain{DB: db, Tip: tip}
	return bc
}

// CreateBlockchain establishes a blockchain with a genesis block
func CreateBlockchain(address string) *Blockchain {
	if util.DoesDBExist() {
		fmt.Println("Chroma chain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(conf.DBdbfile, 0600, nil)
	util.CheckAnxiety(err)

	util.CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		genesisBlock := GenerateGenesisBlock(NewCoinbaseTx(address, conf.Message))
		bucket, err := tx.CreateBucket([]byte(conf.DBblocksbucket))
		util.CheckAnxiety(err)
		util.CheckAnxiety(bucket.Put(genesisBlock.Hash, genesisBlock.Serialize()))
		util.CheckAnxiety(bucket.Put([]byte(conf.DBlasthash), genesisBlock.Hash))
		tip = genesisBlock.Hash
		return nil
	}))
	bc := &Blockchain{tip, db}
	return bc
}
