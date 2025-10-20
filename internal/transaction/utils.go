package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

func generate256HashWithTx(tx Transaction) ([]byte, error){
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	
	if err != nil {
		return nil, err
	}

	hash = sha256.Sum256((encoded.Bytes()))

	return hash[:], nil
}