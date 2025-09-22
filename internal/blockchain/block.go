package blockchain


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
