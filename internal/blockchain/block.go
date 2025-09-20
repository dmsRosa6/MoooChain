package blockchain

type Block struct{
	Data []byte
	Hash []byte
	PrevHash []byte
	Nonce int
}

type Blockchain struct{
	Blocks []*Block
}

func (b *Block) DeriveHash() []byte{

	return []byte{}
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

func InitBlockchain() *Blockchain{
	bc := Blockchain{Blocks: []*Block{GenesisBlock()}}
	return &bc
}

func (b *Blockchain) AddBlock(data string) {
	len := len(b.Blocks)
	prev := b.Blocks[len-1]

	newBlock := CreateBlock(data,prev.Hash)

	b.Blocks = append(b.Blocks, newBlock)
}