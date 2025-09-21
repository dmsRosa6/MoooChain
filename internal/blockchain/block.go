package blockchain

import (
	"bytes"
	"encoding/gob"
)

type Block struct{
	Data []byte
	Hash []byte
	PrevHash []byte
	Nonce int
}

func CreateBlock(data string, prevHash []byte) *Block{
	block := Block{
		Data: []byte(data),
		PrevHash: prevHash,
		Nonce: 0,
	}

	proof := NewProof(&block)
	nonce, hash := proof.Run() 
	block.Hash = hash 
	block.Nonce = nonce
	return &block
}

func GenesisBlock() *Block {
	b := CreateBlock("Genesis",[]byte{})
	return b
}

func (b *Block) Serizalize() ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)	

	err := encoder.Encode(b)

	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func Deserizalize(buf []byte) (*Block, error) {
	reader := bytes.NewReader(buf)
	encoder := gob.NewDecoder(reader)

	var b Block
	err := encoder.Decode(&b)

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func Handle(e error){
}