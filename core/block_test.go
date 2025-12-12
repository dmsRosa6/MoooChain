package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/dmsRosa6/MoooChain/crypto"
	"github.com/dmsRosa6/MoooChain/types"
	"github.com/stretchr/testify/assert"
)

func TestHeader_Encode_Then_Decode(t *testing.T) {
	header := Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Nonce:     1,
		Timestamp: time.Now().UnixNano(),
		Height:    0,
	}

	buf := bytes.Buffer{}

	err := header.EncodeHeader(&buf)
	assert.Nil(t, err)

	newHeader, err := DecodeHeader(&buf)
	assert.Nil(t, err)

	assert.Equal(t, header.Version, newHeader.Version)
	assert.Equal(t, header.PrevBlock, newHeader.PrevBlock)
	assert.Equal(t, header.Nonce, newHeader.Nonce)
	assert.Equal(t, header.Timestamp, newHeader.Timestamp)
	assert.Equal(t, header.Height, newHeader.Height)

	assert.Equal(t, header, newHeader)
}

func TestBlock_Encode_Then_Decode(t *testing.T) {
    header := Header{
        Version:   1,
        PrevBlock: types.RandomHash(),
        Nonce:     1,
        Timestamp: time.Now().UnixNano(),
        Height:    0,
    }

    block := Block{Header: header, Data: []Transaction{}}

    var buf bytes.Buffer

    err := block.EncodeBlock(&buf, BlockEncoder{})
    assert.Nil(t, err)

    dec := BlockDecoder{}
    err = DecodeBlock(&buf, &dec)
    assert.Nil(t, err)

    newBlock := dec.Decode(&buf)
    assert.NotNil(t, newBlock)

    assert.Equal(t, block, newBlock)
}


func TestBlock_Hash(t *testing.T) {
	header := Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Nonce:     1,
		Timestamp: time.Now().UnixNano(),
		Height:    0,
	}

	block := Block{Header: header, Data: []Transaction{}}

	assert.Nil(t, block.hash)

	h := block.Hash(BlockHasher{})
	assert.Equal(t, h.IsZero(), false)

	h1 := block.Hash(BlockHasher{})
	assert.Equal(t, h, h1)
	assert.Equal(t, h.IsZero(), false)
}

func TestSignVerifyBlockSuccess(t *testing.T) {
	header := Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Nonce:     1,
		Timestamp: time.Now().UnixNano(),
		Height:    0,
	}

	block := Block{Header: header, Data: []Transaction{}}

	privKey := crypto.NewPrivKey()

	err := block.Sign(privKey)

	assert.Nil(t,err)
	
	err = block.Verify()
	
	assert.Nil(t,err)

}

func TestSignVerifyBlockBadSignature(t *testing.T) {
	header := Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Nonce:     1,
		Timestamp: time.Now().UnixNano(),
		Height:    0,
	}

	block := Block{Header: header, Data: []Transaction{}}

	privKey := crypto.NewPrivKey()

	privKey2 := crypto.NewPrivKey()

	err := block.Sign(privKey)

	assert.Nil(t,err)

	sign2, err := privKey2.Sign(block.HeaderData())
	
	assert.Nil(t,err)

	block.Signature = sign2

	err = block.Verify()
	
	assert.NotNil(t,err)

}

func TestSignVerifyBlockBadHash(t *testing.T) {
	header := Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Nonce:     1,
		Timestamp: time.Now().UnixNano(),
		Height:    0,
	}

	block := Block{Header: header, Data: []Transaction{}}

	privKey := crypto.NewPrivKey()

	err := block.Sign(privKey)

	assert.Nil(t,err)

	sign2, err := privKey.Sign([]byte("bad"))
	
	assert.Nil(t,err)

	block.Signature = sign2

	err = block.Verify()
	
	assert.NotNil(t,err)

}