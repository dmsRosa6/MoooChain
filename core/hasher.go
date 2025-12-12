package core

import (
	"crypto/sha256"

	"github.com/dmsRosa6/MoooChain/types"
)

type Hasher[T any] interface{
	Hash(T) types.Hash
}

type BlockHasher struct {

}

func (bh BlockHasher) Hash(block *Block) types.Hash{
	h := sha256.Sum256(block.HeaderData())

	return types.Hash(h)
}