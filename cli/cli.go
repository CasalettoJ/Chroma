package cli

import (
	"flag"
	"fmt"
	"os"

	conf "github.com/casalettoj/chroma/constants"
	util "github.com/casalettoj/chroma/utils"
)

// Run runs cli flags
func Run() {
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
		PrintChain()
	}

	if createBlockchainCommand.Parsed() {
		ValidateRequiredOption(*createAddress)
		CreateBlockchain(*createAddress)
	}

	if getBalanceCommand.Parsed() {
		ValidateRequiredOption(*balanceAddress)
		GetBalance(*balanceAddress)
	}

	if sendCommand.Parsed() {
		ValidateRequiredOption(*sendTo)
		ValidateRequiredOption(*sendFrom)
		Send(*sendFrom, *sendTo, *sendAmount)
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
