package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

const difficulty = 12


type ProofOfWork struct{
	Block *Block
	Target *big.Int
}

func NewProof(b * Block) *ProofOfWork{
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	pow := &ProofOfWork{b,target}

	return pow
}


func (pow *ProofOfWork) InitData(nonce int) []byte{
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			toBytes(int64(nonce)),
			toBytes(difficulty),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte){
	var intHash big.Int
	var hash [32]byte

	nonce := 0 

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		intHash.SetBytes(hash[:]) 
	
		if intHash.Cmp(pow.Target) == -1 {
			break
		}else{
			nonce++
		}
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	var hash [32]byte

	data := pow.InitData(pow.Block.Nonce)
	hash = sha256.Sum256(data)

	intHash.SetBytes(hash[:]) 

	return intHash.Cmp(pow.Target) == -1
}

func toBytes(num int64) []byte {
    buf := make([]byte, 8)
    binary.BigEndian.PutUint64(buf, uint64(num))
    return buf
}

func toHexString(num int64) string {
    return fmt.Sprintf("%x", num)
}

