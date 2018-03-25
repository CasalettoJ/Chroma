package blockchain

import (
	"flag"
	"fmt"
	"os"
)

// CLI holds a blockchain and operates with it for given flags
type CLI struct{}

// PrintHelp prints CLI usage
func PrintHelp() {
	fmt.Println("Usage: ")
	fmt.Println("  getbalance -address {ADDRESS} - Get balance of ADDRESS")
	fmt.Println("  createblockchain -address {ADDRESS} - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from {FROM} -to {TO} -amount {AMOUNT} - Send AMOUNT of coins from FROM address to TO")
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

// CreateBlockchain creates a blockchain db
func (cli *CLI) CreateBlockchain(address string) {
	bc := CreateBlockchain(address)
	bc.DB.Close()
	fmt.Println("CHROMA chain created")
}

// PrintChain iterates through the chain and prints the data of each
func (cli *CLI) PrintChain() {
	bc := OpenBlockchain()
	defer bc.DB.Close()
	bci := bc.Iterator()

	for {
		block := bci.Next()
		if block.PrevHash != nil {
			fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		}
		fmt.Printf("Tx Hash: %x\n", block.HashTransactions())
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println("Transactions:")
		for i, tx := range block.Transactions {
			fmt.Printf("TX %d: %d inputs %d outputs\n", i, len(tx.Vin), len(tx.Vout))
			for j, input := range tx.Vin {
				fmt.Printf("Input %d:\n%+v\n", j, input)
			}
			for j, output := range tx.Vout {
				fmt.Printf("Output %d:\n%+v\n", j, output)
			}
		}
		fmt.Println()
		if bci.IsGenesisBlock() {
			break
		}
	}
}

// GetBalance prints the balance of a given address to the console
func (cli *CLI) GetBalance(address string) {
	total := 0
	bc := OpenBlockchain()
	defer bc.DB.Close()
	UTXOs := bc.GetUTXOs(address)

	for _, UTXO := range UTXOs {
		total += UTXO.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, total)
}

// Run runs cli flags
func (cli *CLI) Run() {
	ValidateArgs()
	createBlockchainCommand := flag.NewFlagSet(CLIcreateblockchain, flag.PanicOnError)
	createAddress := createBlockchainCommand.String(CLIaddress, "", "Reward Address")

	getBalanceCommand := flag.NewFlagSet(CLIgetbalance, flag.PanicOnError)
	balanceAddress := getBalanceCommand.String(CLIaddress, "", "Balance Address")

	printChainCommand := flag.NewFlagSet(CLIprintchain, flag.PanicOnError)

	switch os.Args[1] {
	case CLIcreateblockchain:
		CheckAnxiety(createBlockchainCommand.Parse(os.Args[2:]))
	case CLIprintchain:
		CheckAnxiety(printChainCommand.Parse(os.Args[2:]))
	case CLIgetbalance:
		CheckAnxiety(getBalanceCommand.Parse(os.Args[2:]))
	default:
		CLIFailure()
	}

	if printChainCommand.Parsed() {
		cli.PrintChain()
	}

	if createBlockchainCommand.Parsed() {
		if *createAddress == "" {
			CLIFailure()
		}
		cli.CreateBlockchain(*createAddress)
	}

	if getBalanceCommand.Parsed() {
		if *balanceAddress == "" {
			CLIFailure()
		}
		cli.GetBalance(*balanceAddress)
	}

}
