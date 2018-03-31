package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
	conf "github.com/casalettoj/chroma/constants"
	util "github.com/casalettoj/chroma/utils"
	"golang.org/x/crypto/ripemd160"
)

// Wallet holds a private key
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// GetChromaAddress returns a public CHROMA address for a wallet
func (wa *Wallet) GetChromaAddress() []byte {
	hashedKey := HashPublicKey(wa.PublicKey)
	payload := append([]byte{conf.Version}, hashedKey...)
	checksum := Checksum(payload)
	address := base58.Encode(append(payload, checksum...))
	return []byte(address)
}

// NewWallet creates a new wallet
func NewWallet() *Wallet {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	util.CheckAnxiety(err)
	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	public = append([]byte{conf.UncompressedPubKeyPrefix}, public...) // Append prefix for an uncompressed public key -- may do key compression later
	return &Wallet{PrivateKey: *private, PublicKey: public}
}

// HashPublicKey takes a public key and hashes its SHA256 hash w/ RIPEMD160 to return a public key hash
func HashPublicKey(pk []byte) []byte {
	hashedKey := sha256.Sum256(pk)
	RIPEMD160hasher := ripemd160.New()
	_, err := RIPEMD160hasher.Write(hashedKey[:])
	util.CheckAnxiety(err)
	return RIPEMD160hasher.Sum(nil)
}

// Checksum hashes a byte array twice with sha256 and returns a bytearray of AddressChecksumLen length
func Checksum(pl []byte) []byte {
	hashedKey := sha256.Sum256(pl)
	hashedKey = sha256.Sum256(hashedKey[:])
	return hashedKey[:conf.AddressChecksumLen]
}
