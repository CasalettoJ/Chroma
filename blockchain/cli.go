package blockchain

import (
	"flag"
	"fmt"
	"os"
)

// CLI holds a blockchain and operates with it for given flags
type CLI struct {
	BC *Blockchain
}

// PrintHelp prints CLI usage
func PrintHelp() {
	fmt.Println("Usage: ")
	fmt.Println("  addblock -data {BLOCK_DATA}  -- Adds a block with given data to the chain.")
	fmt.Println("  printchain -- prints all of the blocks in the chain")
}

// CLIFailure prints CLI usage and exits with an error
func CLIFailure() {
	PrintHelp()
	os.Exit(1)
}

// ValidateArgs ensures flag validity
func ValidateArgs() {
	if len(os.Args) < 2 {
		CLIFailure()
	}
}

// AddBlock creates a new block on the chain
func (cli *CLI) AddBlock(data string) {
	cli.BC.AddBlock(data)
	fmt.Println("Block added.")
}

// PrintChain iterates through the chain and prints the data of each
func (cli *CLI) PrintChain() {
	bci := cli.BC.Iterator()

	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// Run runs cli flags
func (cli *CLI) Run() {
	addBlockCommand := flag.NewFlagSet(CLIaddblock, flag.PanicOnError)
	printChainCommand := flag.NewFlagSet(CLIprintchain, flag.PanicOnError)
	addblockData := addBlockCommand.String(CLIdata, "", "Data to be added")

	switch os.Args[1] {
	case CLIaddblock:
		CheckAnxiety(addBlockCommand.Parse(os.Args[2:]))
	case CLIprintchain:
		CheckAnxiety(printChainCommand.Parse(os.Args[2:]))
	default:
		CLIFailure()
	}

	if addBlockCommand.Parsed() {
		if *addblockData == "" {
			CLIFailure()
		}
		cli.AddBlock(*addblockData)
	}

	if printChainCommand.Parsed() {
		cli.PrintChain()
	}

}
