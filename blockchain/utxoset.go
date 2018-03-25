package blockchain

import (
	"encoding/hex"

	bolt "github.com/coreos/bbolt"
)

// GetUTXOsForAddress returns all unspent tx outputs for a given address
func GetUTXOsForAddress(bc *Blockchain, address string) []TxOutput {
	db := bc.DB
	var UTXOs []TxOutput
	CheckAnxiety(db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DButxobucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			utxoutputs := DeserializeTxOutputs(v)
			for _, utxo := range utxoutputs.Outputs {
				if utxo.Unlockable(address) {
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
	CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(DButxobucket))
		if err != bolt.ErrBucketNotFound {
			CheckAnxiety(err)
		}
		_, err = tx.CreateBucket([]byte(DButxobucket))
		CheckAnxiety(err)
		return nil
	}))
	UTXOs := bc.GetUTXOs()
	CheckAnxiety(db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DButxobucket))
		for txID, utxo := range UTXOs {
			key, err := hex.DecodeString(txID)
			CheckAnxiety(err)
			CheckAnxiety(bucket.Put(key, utxo.Serialize()))
		}
		return nil
	}))
}
