package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const maxNonce = math.MaxInt64
const targetBits = 18

// ProofOfWork is a struture containing difficulty target and block being mined.
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork does the obvious
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target = target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}

// PrepareData returns a []byte of all headers, the target bits, and nonce.
func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevHash,
		pow.block.HashTransactions(),
		Int64ToByteArray(pow.block.Timestamp),
		Int64ToByteArray(int64(targetBits)),
		Int64ToByteArray(int64(nonce)),
	}, []byte{})
	return data
}

// Run runs the proof of work algorithm until mined
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < maxNonce {
		preparedData := pow.PrepareData(nonce)
		hash = sha256.Sum256(preparedData)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

// IsValid returns whether the PoW is valid.
func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int

	data := pow.PrepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
