package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/dmsRosa6/MoooChain/types"
	"github.com/stretchr/testify/assert"
)

func TestHeader_Encode_Then_Decode(t *testing.T){
	header := Header{
		Version: 1,
		PrevBlock: types.RandomHash(),
		Nonce: 1,
		Timestamp: time.Now().UnixNano(),
		Height: 0,
	}

	buf := bytes.Buffer{}

	err := header.EncodeHeader(&buf)
	assert.Nil(t,err)

	newHeader, err := DecodeHeader(&buf)
	assert.Nil(t,err)


	assert.Equal(t,header.Version ,newHeader.Version)
	assert.Equal(t,header.PrevBlock ,newHeader.PrevBlock)
	assert.Equal(t,header.Nonce ,newHeader.Nonce)
	assert.Equal(t,header.Timestamp ,newHeader.Timestamp)
	assert.Equal(t,header.Height ,newHeader.Height)

	assert.Equal(t, header, newHeader)
}


func TestBlock_Encode_Then_Decode(t *testing.T){
	header := Header{
		Version: 1,
		PrevBlock: types.RandomHash(),
		Nonce: 1,
		Timestamp: time.Now().UnixNano(),
		Height: 0,
	}

	block := Block{Header: header, Data: []Transaction{}}

	buf := bytes.Buffer{}

	err := block.EncodeBlock(&buf)
	assert.Nil(t,err)

	newBlock, err := DecodeBlock(&buf)
	assert.Nil(t,err)

	newHeader := newBlock.Header

	assert.Equal(t,header.Version ,newHeader.Version)
	assert.Equal(t,header.PrevBlock ,newHeader.PrevBlock)
	assert.Equal(t,header.Nonce ,newHeader.Nonce)
	assert.Equal(t,header.Timestamp ,newHeader.Timestamp)
	assert.Equal(t,header.Height ,newHeader.Height)

	assert.Equal(t, header, newHeader)
	assert.Equal(t, block.Data, newBlock.Data)

	assert.Equal(t, block, newBlock)

}


func TestBlock_Hash(t *testing.T){
	header := Header{
		Version: 1,
		PrevBlock: types.RandomHash(),
		Nonce: 1,
		Timestamp: time.Now().UnixNano(),
		Height: 0,
	}

	block := Block{Header: header, Data: []Transaction{}}

	assert.Nil(t, block.cachedHash)

	h := block.Hash()
	assert.Equal(t, h.IsZero(), false)

	h1 := block.Hash()
	assert.Equal(t, h, h1)
	assert.Equal(t, h.IsZero(), false)
}