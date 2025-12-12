package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/dmsRosa6/MoooChain/crypto"
	"github.com/dmsRosa6/MoooChain/types"
)

type Header struct {
	Version   uint32     `json:"version"`
	PrevBlock types.Hash `json:"prev_block"`
	Nonce     uint64     `json:"nonce"`
	Timestamp int64      `json:"ts"`
	Height    uint64     `json:"height"`
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
	Header `json:"header"`
	Data   []Transaction `json:"data"`
	Signature *crypto.Signature
	Validator crypto.PubKey
	hash types.Hash
}

func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}

	return b.hash
}

func (b *Block) EncodeBlock(w io.Writer, encoder Encoder[*Block]) error {
	return encoder.Encode(w)
}

func DecodeBlock(r io.Reader, decoder Decoder[*Block])  error {
	return decoder.Decode(r)
}

func (b * Block) HeaderData() []byte{

	buf := bytes.Buffer{}

	encoder := gob.NewEncoder(&buf)

	encoder.Encode(b.Header)

	return buf.Bytes()
}

func (b * Block) Sign(privKey crypto.PrivKey) error{

	sign, err := privKey.Sign(b.HeaderData())

	if err != nil {
		return err
	}
	
	b.Validator = privKey.PubKey()
	b.Signature = sign

	return nil
}


func (b * Block) Verify() error{
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.HeaderData(), b.Validator) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func CreateBlock(txs []Transaction, prevBlock []byte) *Block {

	header := Header{
		PrevBlock: types.HashFromBytes(prevBlock),
		Nonce:     0,
	}

	block := Block{
		Header: header,
		Data:   txs,
	}

	//TODO finish this
	//proof := NewProof(&block)
	//nonce, hash := proof.Run()
	//block.Hash = hash
	//block.Nonce = nonce
	return &block
}

func GenesisBlock(mintTx Transaction) *Block {
	b := CreateBlock([]Transaction{mintTx}, []byte{})
	return b
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Data {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
