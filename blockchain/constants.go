package blockchain

const (
	// Message memes something
	Message = "09 F9 11 02 9D 74 E3 5B D8 41 56 C5 63 56 88 C0"
	// DBdbfile is the database filename
	DBdbfile = "chroma_db"
	// DBblocksbucket is the name of the bolt bucket the blocks are stored in
	DBblocksbucket = "blocks"
	// DBlasthash is the key the hash of the tip of the chain is stored in
	DBlasthash = "l"

	// CLIcreateblockchain is the command to create a new DB
	CLIcreateblockchain = "createblockchain"
	// CLIprintchain is the argument for printing the chain to the console
	CLIprintchain = "printchain"
	// CLIaddress is an option flag for an address
	CLIaddress = "address"

	// TXcoinbaseaward is the amount of coins awarded for mining a block
	TXcoinbaseaward = 10
)
