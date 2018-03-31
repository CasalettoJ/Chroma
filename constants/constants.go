package constants

const (
	// Message memes something
	Message = "09 F9 11 02 9D 74 E3 5B D8 41 56 C5 63 56 88 C0"
	// DBdbfile is the database filename
	DBdbfile = "chroma_db"
	// DBblocksbucket is the name of the bolt bucket the blocks are stored in, keyed by hash.
	DBblocksbucket = "blocks"
	//DBtxbucket is the name of the bolt bucket transactions are stored in, keyed by ID.
	DBtxbucket = "transactions"
	//DButxobucket is the name of the bolt bucket UTXOs are stored in, keyed by TXID
	DButxobucket = "utxoset"
	// DBlasthash is the key the hash of the tip of the chain is stored in
	DBlasthash = "lasthash"

	// TXcoinbaseaward is the amount of coins awarded for mining a block
	TXcoinbaseaward = 1000

	// CLIcreateblockchain is the command to create a new DB
	CLIcreateblockchain = "createblockchain"
	// CLIprintchain is the command for printing the chain to the console
	CLIprintchain = "printchain"
	// CLIgetbalance is the command for retrieving the balance of an address
	CLIgetbalance = "getbalance"
	// CLIsend is the command for creating a new transaction
	CLIsend = "send"

	// CLIaddress is an option flag for an address
	CLIaddress = "address"
	// CLIfrom is the option flag for a sender address
	CLIfrom = "from"
	// CLIto is the option flag for a recipient address
	CLIto = "to"
	// CLIamount is the option flag for an amount of coins
	CLIamount = "amount"

	// Version is the 1-byte version of the wallet.
	Version = byte(0x00)
	// UncompressedPubKeyPrefix is the 1-byte prefix of an uncompressed public key. Like bitcoin!
	UncompressedPubKeyPrefix = byte(0x04)
	// WalletFile is the wallet filename
	WalletFile = "wallet.dat"
	// AddressChecksumLen is the number of bytes to take after hashing public key for checksum
	AddressChecksumLen = 4
)
