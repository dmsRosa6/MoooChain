package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
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
	Validator *crypto.PubKey
	hash types.Hash
}

func (b *Block) Sign(key crypto.PrivKey) error{
    
	h := b.Hash(BlockHasher{})
    sign, err := key.Sign(h[:])

	if err != nil {
		return err
	}

	b.Signature = sign

	return nil
}

func (b *Block) Verify() error{

    if b.Signature == nil{
		return errors.New("transaction has no signature")
	}

	if isVerified := b.Signature.Verify(b.Data, b.PubKey); !isVerified{
        return errors.New("invalid transation signature")
	}

	return nil
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
