package blockchain

import "encoding/json"


type Block struct {
    Data     []byte `json:"data"`
    Hash     []byte `json:"hash"`
    PrevHash []byte `json:"prev_hash"`
    Nonce    int    `json:"nonce"`
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

func (b *Block) ToJSON() ([]byte, error) {
    return json.Marshal(b)
}