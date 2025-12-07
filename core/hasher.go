package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/dmsRosa6/MoooChain/types"
)

type Hasher[T any] interface{
	Hash(T) types.Hash
}

type BlockHasher struct {

}

func (bh BlockHasher) Hash(block *Block) types.Hash{
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(block); err != nil{
		panic(err)
	}

	h := sha256.Sum256(buf.Bytes())

	return types.Hash(h)
}