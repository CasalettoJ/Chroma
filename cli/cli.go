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
	validateArgs()
	createBlockchainCommand := flag.NewFlagSet(conf.CLIcreateblockchain, flag.PanicOnError)
	createAddress := createBlockchainCommand.String(conf.CLIaddress, "", "Reward Address")

	getBalanceCommand := flag.NewFlagSet(conf.CLIgetbalance, flag.PanicOnError)
	balanceAddress := getBalanceCommand.String(conf.CLIaddress, "", "Balance Address")

	printChainCommand := flag.NewFlagSet(conf.CLIprintchain, flag.PanicOnError)

	sendCommand := flag.NewFlagSet(conf.CLIsend, flag.PanicOnError)
	sendTo := sendCommand.String(conf.CLIto, "", "To Address")
	sendFrom := sendCommand.String(conf.CLIfrom, "", "From Address")
	sendAmount := sendCommand.Int(conf.CLIamount, 0, "Amout to send")

	newWalletCommand := flag.NewFlagSet(conf.CLInewwallet, flag.PanicOnError)

	printWalletsCommand := flag.NewFlagSet(conf.CLIprintwallets, flag.PanicOnError)

	switch os.Args[1] {
	case conf.CLIcreateblockchain:
		util.CheckAnxiety(createBlockchainCommand.Parse(os.Args[2:]))
	case conf.CLIprintchain:
		util.CheckAnxiety(printChainCommand.Parse(os.Args[2:]))
	case conf.CLIgetbalance:
		util.CheckAnxiety(getBalanceCommand.Parse(os.Args[2:]))
	case conf.CLIsend:
		util.CheckAnxiety(sendCommand.Parse(os.Args[2:]))
	case conf.CLInewwallet:
		util.CheckAnxiety(newWalletCommand.Parse(os.Args[2:]))
	case conf.CLIprintwallets:
		util.CheckAnxiety(printWalletsCommand.Parse(os.Args[2:]))
	default:
		failure()
	}

	if printChainCommand.Parsed() {
		printChain()
	}

	if createBlockchainCommand.Parsed() {
		validateRequiredOption(*createAddress)
		createBlockchain(*createAddress)
	}

	if getBalanceCommand.Parsed() {
		validateRequiredOption(*balanceAddress)
		getBalance(*balanceAddress)
	}

	if sendCommand.Parsed() {
		validateRequiredOption(*sendTo)
		validateRequiredOption(*sendFrom)
		send(*sendFrom, *sendTo, *sendAmount)
	}

	if newWalletCommand.Parsed() {
		createNewWallet()
	}

	if printWalletsCommand.Parsed() {
		printWallets()
	}
}

// validateRequiredOption quits if an option is not supplied
func validateRequiredOption(option string) {
	if option == "" {
		failure()
	}
}

// printHelp prints CLI usage
func printHelp() {
	fmt.Println("Usage: ")
	fmt.Println("  getbalance -address {ADDRESS} - Get balance of ADDRESS")
	fmt.Println("  newwallet - Create a new CHROMA address")
	fmt.Println("  printwallets - print all CHROMA addresses in the wallet")
	fmt.Println("  createblockchain -address {ADDRESS} - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from {FROM} -to {TO} -amount {AMOUNT} - Send AMOUNT of coins from FROM address to TO")
}

// failure prints CLI usage and exits with an error
func failure() {
	printHelp()
	os.Exit(1)
}

// validateArgs ensures flag validity
func validateArgs() {
	if len(os.Args) < 2 {
		failure()
	}
}
