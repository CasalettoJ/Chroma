package main

import (
	"github.com/casalettoj/chroma/blockchain"
)

func main() {
	cli := &blockchain.CLI{}
	cli.Run()
}

/*
* TODO
* 	Save buckets for UTXO and TXs to DB
*	Get all UTXOs from blockchain
*	create UTXOSet structure and functions
*	Serialize TX and TxOutput
*
 */
