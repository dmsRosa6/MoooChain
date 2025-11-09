package blockchain

import (
	"bytes"
	"crypto/sha256"

	"github.com/dmsRosa6/MoooChain/internal/transaction"
)


type Block struct {
    Data     []*transaction.Transaction `json:"data"`
    Hash     []byte `json:"hash"`
    PrevHash []byte `json:"prev_hash"`
    Nonce    int    `json:"nonce"`
}

func CreateBlock(txs []*transaction.Transaction, prevHash []byte) *Block{
	block := Block{
		Data: txs,
		PrevHash: prevHash,
		Nonce: 0,
	}

	proof := NewProof(&block)
	nonce, hash := proof.Run() 
	block.Hash = hash 
	block.Nonce = nonce
	return &block
}

func GenesisBlock(mintTx *transaction.Transaction) *Block {
	b := CreateBlock([]*transaction.Transaction{mintTx},[]byte{})
	return b
}

func (b *Block) HashTransactions() []byte{
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Data{
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
