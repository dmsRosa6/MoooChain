package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)


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

func (b *Block) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        Data     string `json:"data"`
        Hash     string `json:"hash"`
        PrevHash string `json:"prev_hash"`
        Nonce    int    `json:"nonce"`
    }{
        Data:     string(b.Data),
        Hash:     hex.EncodeToString(b.Hash),
        PrevHash: hex.EncodeToString(b.PrevHash),
        Nonce:    b.Nonce,
    })
}

func (b *Block) UnmarshalJSON(data []byte) error {
	aux := struct {
		Data     string `json:"data"`
		Hash     string `json:"hash"`
		PrevHash string `json:"prev_hash"`
		Nonce    int    `json:"nonce"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var (
		hashBytes     []byte
		prevHashBytes []byte
		err           error
	)

	if aux.Hash != "" {
		hashBytes, err = hex.DecodeString(aux.Hash)
		if err != nil {
			return fmt.Errorf("invalid hex in hash: %w", err)
		}
	}

	if aux.PrevHash != "" {
		prevHashBytes, err = hex.DecodeString(aux.PrevHash)
		if err != nil {
			return fmt.Errorf("invalid hex in prev_hash: %w", err)
		}
	}

	b.Data = []byte(aux.Data)
	b.Hash = hashBytes
	b.PrevHash = prevHashBytes
	b.Nonce = aux.Nonce
	return nil
}