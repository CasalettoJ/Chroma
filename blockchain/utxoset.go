package blockchain

import (
	"encoding/hex"

	conf "github.com/casalettoj/chroma/constants"
	util "github.com/casalettoj/chroma/utils"
	bolt "github.com/coreos/bbolt"
)

// FindUTXOsForPayment searches through the UTXOSet for unlockable UTXOs until the amount is reached
// returns the amount of all retrieved UTXOs and a map of TxIDs and UTXO indices
func FindUTXOsForPayment(bc *Blockchain, pubKeyHash []byte, amount int) (int, map[string][]int) {
	accumulated := 0
	UTXOIndices := make(map[string][]int)
	db := bc.DB

	util.CheckAnxiety(db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DButxobucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			txID := hex.EncodeToString(k)
			UTXOs := DeserializeTxOutputs(v)
			for UTXOIndex, UTXO := range UTXOs.Outputs {
				if accumulated > amount {
					break
				}
				if UTXO.Unlockable(pubKeyHash) {
					accumulated += UTXO.Value
					UTXOIndices[txID] = append(UTXOIndices[txID], UTXOIndex)
				}
			}
		}
		return nil
	}))
	return accumulated, UTXOIndices
}

// GetUTXOsForAddress returns all unspent tx outputs for a given address
func GetUTXOsForAddress(bc *Blockchain, pubKeyHash []byte) []TxOutput {
	db := bc.DB
	var UTXOs []TxOutput
	util.CheckAnxiety(db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DButxobucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			utxoutputs := DeserializeTxOutputs(v)
			for _, utxo := range utxoutputs.Outputs {
				if utxo.Unlockable(pubKeyHash) {
					UTXOs = append(UTXOs, utxo)
				}
			}
		}
		return nil
	}))
	return UTXOs
}

// ReindexUTXOs deletes the current UTXO set from db and creates a new set
func ReindexUTXOs(bc *Blockchain) {
	db := bc.DB
	util.CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(conf.DButxobucket))
		if err != bolt.ErrBucketNotFound {
			util.CheckAnxiety(err)
		}
		_, err = tx.CreateBucket([]byte(conf.DButxobucket))
		util.CheckAnxiety(err)
		return nil
	}))
	UTXOsByTxID := bc.GetUTXOs()
	util.CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DButxobucket))
		for txID, utxos := range UTXOsByTxID {
			key, err := hex.DecodeString(txID)
			util.CheckAnxiety(err)
			util.CheckAnxiety(bucket.Put(key, utxos.Serialize()))
		}
		return nil
	}))
}

// UpdateUTXOs takes the newest block, removes all outputs that were used as inputs in its transactions
// and adds the outputs of each Tx as new UTXOs in the set.
func UpdateUTXOs(bc *Blockchain, b *Block) {
	db := bc.DB
	util.CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		utxoBucket := tx.Bucket([]byte(conf.DButxobucket))
		for _, transaction := range b.Transactions {
			// If the tx is a coinbase tx, ignore the inputs entirely
			if !transaction.IsCoinbaseTx() {
				for _, input := range transaction.Vin {
					// Check the last TX's output's index and if it wasn't the index used in the input (Vout)
					// then it is still unspent and should be in the new UTXOs of the last TX.
					updatedUTXOs := TxOutputs{}
					prevTxUTXOsBytes := utxoBucket.Get(input.TxID)
					prevTxUTXOs := DeserializeTxOutputs(prevTxUTXOsBytes)
					for prevUTXOIndex, prevUTXO := range prevTxUTXOs.Outputs {
						if input.Vout != prevUTXOIndex {
							updatedUTXOs.Outputs = append(updatedUTXOs.Outputs, prevUTXO)
						}
					}
					// Then if the TX has no more UTXOs remove it from the bucket
					// Otherwise, update the TXID-indexed TxOutputs with the updated structure
					if len(updatedUTXOs.Outputs) == 0 {
						util.CheckAnxiety(utxoBucket.Delete(input.TxID))
					} else {
						util.CheckAnxiety(utxoBucket.Put(input.TxID, updatedUTXOs.Serialize()))
					}
				}
			}

			// Next, place all of the new TxOutputs from the new block into the UTXOset
			newUTXOs := TxOutputs{}
			for _, output := range transaction.Vout {
				newUTXOs.Outputs = append(newUTXOs.Outputs, output)
			}
			util.CheckAnxiety(utxoBucket.Put(transaction.ID, newUTXOs.Serialize()))
		}
		return nil
	}))
}
