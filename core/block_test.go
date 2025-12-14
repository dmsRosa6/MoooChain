package core

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/dmsRosa6/MoooChain/crypto"
	"github.com/dmsRosa6/MoooChain/types"
	"github.com/stretchr/testify/assert"
)

func RandomInt64() uint64 {
    return rand.Uint64()
}

func RandomHeader(height uint64) *Header {
	return &Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Nonce:     RandomInt64(),
		Timestamp: time.Now().UnixNano(),
		Height:    height,
	}
}

// TODO Add tx's in the future
func RandomBlock(dataSize int, height uint64) *Block {
	header := RandomHeader(height)
	data := make([]Transaction, dataSize)
	for i := 0; i < dataSize; i++ {
		
	}
	return &Block{
		Header: header,
		Data:   data,
	}
}

func RandomBlockWithSig(dataSize int, height uint64) *Block {
	privkey := crypto.NewPrivKey()
	b := RandomBlock(dataSize, height)
	b.Sign(privkey)

	return b
}

// --- Tests ---

func TestHeader_EncodeDecode(t *testing.T) {
	header := RandomHeader(0)
	buf := bytes.Buffer{}

	assert.NoError(t, header.EncodeHeader(&buf))

	newHeader, err := DecodeHeader(&buf)
	assert.NoError(t, err)
	assert.Equal(t, header, newHeader)
}

func TestBlock_EncodeDecode(t *testing.T) {
	block := RandomBlock(3,0)
	buf := bytes.Buffer{}

	assert.NoError(t, block.EncodeBlock(&buf, BlockEncoder{}))

	dec := BlockDecoder{}
	assert.NoError(t, DecodeBlock(&buf, &dec))

	newBlock := dec.Decode(&buf)
	assert.NotNil(t, newBlock)
}

func TestBlock_Hash(t *testing.T) {
	block := RandomBlock(0,0)
	assert.NotNil(t, block.hash)

	h := block.Hash(BlockHasher{})
	assert.False(t, h.IsZero())

	h2 := block.Hash(BlockHasher{})
	assert.Equal(t, h, h2)
}

func TestSignVerifyBlock_Success(t *testing.T) {
	block := RandomBlock(2,1)
	priv := crypto.NewPrivKey()

	assert.NoError(t, block.Sign(priv))
	assert.NoError(t, block.Verify())
}

func TestSignVerifyBlock_BadSignature(t *testing.T) {
	block := RandomBlock(2,2)
	priv1 := crypto.NewPrivKey()
	priv2 := crypto.NewPrivKey()

	assert.NoError(t, block.Sign(priv1))

	// forge with different key
	sign2, err := priv2.Sign(block.HeaderData())
	assert.NoError(t, err)
	block.Signature = sign2

	assert.Error(t, block.Verify())
}

func TestSignVerifyBlock_BadHash(t *testing.T) {
	block := RandomBlock(2,2)
	priv := crypto.NewPrivKey()

	assert.NoError(t, block.Sign(priv))

	sign2, err := priv.Sign([]byte("bad"))
	assert.NoError(t, err)
	block.Signature = sign2

	assert.Error(t, block.Verify())
}
