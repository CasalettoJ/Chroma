package cli

import (
	"flag"
	"fmt"
	"os"

	chroma "github.com/casalettoj/chroma/blockchain"
	conf "github.com/casalettoj/chroma/constants"
	util "github.com/casalettoj/chroma/utils"
)

// CLI holds a blockchain and operates with it for given flags
type CLI struct{}

// CreateBlockchain creates a blockchain db
func (cli *CLI) CreateBlockchain(address string) {
	bc := chroma.CreateBlockchain(address)
	defer bc.DB.Close()
	chroma.ReindexUTXOs(bc)
	fmt.Println("CHROMA chain created")
}

// PrintChain iterates through the chain and prints the data of each
func (cli *CLI) PrintChain() {
	bc := chroma.OpenBlockchain()
	defer bc.DB.Close()
	bci := bc.Iterator()

	for {
		fmt.Println()
		fmt.Println("===============================================")
		block := bci.Next()
		if block.PrevHash != nil {
			fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		}
		fmt.Printf("Tx Hash: %x\n", block.HashTransactions())
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println("Transactions:")
		for i, tx := range block.Transactions {
			fmt.Println()
			if i != 0 {
				fmt.Println("____________")
			}
			fmt.Printf("TX %d: %d inputs %d outputs\n", i, len(tx.Vin), len(tx.Vout))
			fmt.Println()
			for j, input := range tx.Vin {
				fmt.Printf("Input %d:\n%+v\n", j, input)
			}
			fmt.Println()
			for j, output := range tx.Vout {
				fmt.Printf("Output %d:\n%+v\n", j, output)
			}
		}
		fmt.Println("===============================================")
		fmt.Println()
		if bci.IsGenesisBlock() {
			break
		}
	}
}

// GetBalance prints the balance of a given address to the console
func (cli *CLI) GetBalance(address string) {
	total := 0
	bc := chroma.OpenBlockchain()
	defer bc.DB.Close()
	UTXOs := chroma.GetUTXOsForAddress(bc, address)

	for _, UTXO := range UTXOs {
		total += UTXO.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, total)
}

// Send creates a TX and CoinbaseTX and mines a new transaction
func (cli *CLI) Send(from, to string, amount int) {
	if amount <= 0 {
		fmt.Println("Invalid amount.")
		os.Exit(1)
	}

	bc := chroma.OpenBlockchain()
	defer bc.DB.Close()

	newTx := chroma.NewTransaction(bc, to, from, amount)
	coinbaseTx := chroma.NewCoinbaseTx(from, "")
	Txs := []*chroma.Transaction{coinbaseTx, newTx}
	newBlock := bc.MineBlock(Txs)
	chroma.UpdateUTXOs(bc, newBlock)
	fmt.Printf("Sent %d to %s.", amount, to)
}

// Run runs cli flags
func (cli *CLI) Run() {
	ValidateArgs()
	createBlockchainCommand := flag.NewFlagSet(conf.CLIcreateblockchain, flag.PanicOnError)
	createAddress := createBlockchainCommand.String(conf.CLIaddress, "", "Reward Address")

	getBalanceCommand := flag.NewFlagSet(conf.CLIgetbalance, flag.PanicOnError)
	balanceAddress := getBalanceCommand.String(conf.CLIaddress, "", "Balance Address")

	printChainCommand := flag.NewFlagSet(conf.CLIprintchain, flag.PanicOnError)

	sendCommand := flag.NewFlagSet(conf.CLIsend, flag.PanicOnError)
	sendTo := sendCommand.String(conf.CLIto, "", "To Address")
	sendFrom := sendCommand.String(conf.CLIfrom, "", "From Address")
	sendAmount := sendCommand.Int(conf.CLIamount, 0, "Amout to send")

	switch os.Args[1] {
	case conf.CLIcreateblockchain:
		util.CheckAnxiety(createBlockchainCommand.Parse(os.Args[2:]))
	case conf.CLIprintchain:
		util.CheckAnxiety(printChainCommand.Parse(os.Args[2:]))
	case conf.CLIgetbalance:
		util.CheckAnxiety(getBalanceCommand.Parse(os.Args[2:]))
	case conf.CLIsend:
		util.CheckAnxiety(sendCommand.Parse(os.Args[2:]))
	default:
		Failure()
	}

	if printChainCommand.Parsed() {
		cli.PrintChain()
	}

	if createBlockchainCommand.Parsed() {
		ValidateRequiredOption(*createAddress)
		cli.CreateBlockchain(*createAddress)
	}

	if getBalanceCommand.Parsed() {
		ValidateRequiredOption(*balanceAddress)
		cli.GetBalance(*balanceAddress)
	}

	if sendCommand.Parsed() {
		ValidateRequiredOption(*sendTo)
		ValidateRequiredOption(*sendFrom)
		cli.Send(*sendFrom, *sendTo, *sendAmount)
	}
}

// ValidateRequiredOption quits if an option is not supplied
func ValidateRequiredOption(option string) {
	if option == "" {
		Failure()
	}
}

// PrintHelp prints CLI usage
func PrintHelp() {
	fmt.Println("Usage: ")
	fmt.Println("  getbalance -address {ADDRESS} - Get balance of ADDRESS")
	fmt.Println("  createblockchain -address {ADDRESS} - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from {FROM} -to {TO} -amount {AMOUNT} - Send AMOUNT of coins from FROM address to TO")
}

// Failure prints CLI usage and exits with an error
func Failure() {
	PrintHelp()
	os.Exit(1)
}

// ValidateArgs ensures flag validity
func ValidateArgs() {
	if len(os.Args) < 2 {
		Failure()
	}
}
