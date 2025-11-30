package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/dmsRosa6/MoooChain/transaction"
	"github.com/dmsRosa6/MoooChain/types"
)

type Header struct {
    Version uint32 `json:"version"`
    PrevBlock types.Hash `json:"prev_block"`
    Nonce    uint64    `json:"nonce"`
	Timestamp int64  `json:"ts"`
	Height uint64 `json:"height"`
}

type Block struct {
	Header  `json:"header"`
    Data     []*transaction.Transaction `json:"data"`
}

func (h * Header) EncodeHeader(w io.Writer) error{
	if err := binary.Write(w, binary.LittleEndian, &h.Version); err != nil{
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.PrevBlock); err != nil{
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.Nonce); err != nil{
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.Timestamp); err != nil{
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.Height); err != nil{
		return err
	}

	return nil
}

func DecodeHeader(r io.Reader) (*Header, error){
	var version uint32
	var prevBlock types.Hash
	var nonce uint64
	var timestamp int64
	var height uint64
	
	if err := binary.Read(r, binary.LittleEndian, version); err != nil{
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, prevBlock); err != nil{
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, nonce); err != nil{
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, timestamp); err != nil{
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, height); err != nil{
		return nil, err
	}

	return &Header{	
					Version: version,
		 			PrevBlock: prevBlock, 
					Nonce: nonce,
					Timestamp: timestamp,
					Height: height,
				  },
				  nil
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
