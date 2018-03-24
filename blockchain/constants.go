package blockchain

const (
	// DBdbfile is the database filename
	DBdbfile = "chroma_db"
	// DBblocksbucket is the name of the bolt bucket the blocks are stored in
	DBblocksbucket = "blocks"
	// DBlasthash is the key the hash of the tip of the chain is stored in
	DBlasthash = "l"

	// CLIaddblock is the argument for adding a block to the chain
	CLIaddblock = "addblock"
	// CLIprintchain is the argument for printing the chain to the console
	CLIprintchain = "printchain"
	// CLIdata is the option for the addblock flag
	CLIdata = "data"
)
