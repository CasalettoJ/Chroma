package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"

	bolt "github.com/coreos/bbolt"
)

// Blockchain is the entire structure of sequential blocks
type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

// Iterator give iterator
func (bc *Blockchain) Iterator() *Iterator {
	iterator := &Iterator{bc.Tip, bc.DB}
	return iterator
}

// GetUTXOs gets all UTXOs unlockable by a given address
func (bc *Blockchain) GetUTXOs(address string) []TxOutput {
	var UTXOs []TxOutput
	UTXs := bc.GetUTXs(address)

	for _, tx := range UTXs {
		for _, output := range tx.Vout {
			if output.Unlockable(address) {
				UTXOs = append(UTXOs, output)
			}
		}
	}

	return UTXOs
}

// GetUTXs gets all UTXs unlockable by given address
func (bc *Blockchain) GetUTXs(address string) []Transaction {
	var unspentTxs []Transaction
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
				// If the output hasn't been spent and it can be unlocked by the address, add it to the UTXOset.
				if !spent && out.Unlockable(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			// If the transaction isn't a coinbase TX, then...
			if !tx.IsCoinbaseTx() {
				// For every input in the transaction...
				for _, in := range tx.Vin {
					// If this address was the one to unlock the output...
					if in.InitiatedBy(address) {
						inTxID := hex.EncodeToString(in.TxID)
						// then add the input's output index to the list of spent TXOs for the tx its referenced UTXO was created in.
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		if bci.IsGenesisBlock() {
			break
		}
	}

	return unspentTxs
}

// OpenBlockchain opens a preexisting blockchain and returns Tip and DB
func OpenBlockchain() *Blockchain {
	if !DoesDBExist() {
		fmt.Println("No existing Chroma chain.  Create DB first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(DBdbfile, 0600, nil)
	CheckAnxiety(err)

	CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DBblocksbucket))
		tip = bucket.Get([]byte(DBlasthash))
		return nil
	}))
	bc := &Blockchain{DB: db, Tip: tip}
	return bc
}

// CreateBlockchain establishes a blockchain with a genesis block
func CreateBlockchain(address string) *Blockchain {
	if DoesDBExist() {
		fmt.Println("Chroma chain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(DBdbfile, 0600, nil)
	CheckAnxiety(err)

	CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		genesisBlock := GenerateGenesisBlock(NewCoinbaseTx(address, Message))
		bucket, err := tx.CreateBucket([]byte(DBblocksbucket))
		CheckAnxiety(err)
		CheckAnxiety(bucket.Put(genesisBlock.Hash, genesisBlock.Serialize()))
		CheckAnxiety(bucket.Put([]byte(DBlasthash), genesisBlock.Hash))
		tip = genesisBlock.Hash
		return nil
	}))
	return &Blockchain{tip, db}
}
