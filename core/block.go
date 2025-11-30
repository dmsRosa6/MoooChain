package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/dmsRosa6/MoooChain/types"
)

type Header struct {
    Version uint32 `json:"version"`
    PrevBlock types.Hash `json:"prev_block"`
    Nonce    uint64    `json:"nonce"`
	Timestamp int64  `json:"ts"`
	Height uint64 `json:"height"`
}

func (h *Header) EncodeHeader(w io.Writer) error {
    return binary.Write(w, binary.LittleEndian, h)
}

func DecodeHeader(r io.Reader) (*Header, error) {
    var h Header
    if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
        return nil, err
    }
    return &h, nil
}


type Block struct {
	Header  `json:"header"`
    Data     []Transaction `json:"data"`
    
	cachedHash types.Hash
}


func (b *Block) Hash() types.Hash{
	buf := bytes.Buffer{}

	b.Header.EncodeHeader(&buf)


	if(b.cachedHash.IsZero()){
		b.cachedHash = types.Hash(sha256.Sum256(buf.Bytes()))
	}

	return b.cachedHash
}

func (b *Block) EncodeBlock(w io.Writer) error {
    return binary.Write(w, binary.LittleEndian, b)
}

func DecodeBlock(r io.Reader) (*Block, error) {
    var b Block
    if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
        return nil, err
    }
    return &b, nil
}

func CreateBlock(txs []Transaction, prevBlock []byte) *Block{
	
	header := Header{
		PrevBlock: types.HashFromBytes(prevBlock),
		Nonce: 0,
	}

	block := Block{
		Header: header,
		Data: txs,
	}

	//TODO finish this
	//proof := NewProof(&block)
	//nonce, hash := proof.Run() 
	//block.Hash = hash 
	//block.Nonce = nonce
	return &block
}

func GenesisBlock(mintTx Transaction) *Block {
	b := CreateBlock([]Transaction{mintTx},[]byte{})
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
