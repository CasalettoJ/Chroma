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
* 	Implement newtransaction
* 	Implement mineblock
*	Implement printpendingtransactions
 */
