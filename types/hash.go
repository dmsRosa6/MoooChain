package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Hash [32]uint8

//TODO maybe i should no panic here
func HashFromBytes(arr []byte) Hash {

	if len(arr) != 32 {
		err := fmt.Sprintf("byte array length is %d, should be exactly 32", len(arr))
		panic(err)
	}

	var value [32]uint8

	copy(value[:], arr)

	return Hash(value)
}

func (h *Hash) IsZero() bool {

	for i := 0; i < len(h); i++ {
		if h[i] != 0 {
			return  false
		}
	}

	return true
}

func RandomBytes(size int) []byte {
	if size < 0 {
		msg := fmt.Sprintf("request size is %d, needed a value greater than 0", size)

		panic(msg)
	}

	token := make([]byte, size)
	if _, err := rand.Read(token); err != nil {
        panic(err)
    }


	return token
}

func RandomHash() Hash {
	return HashFromBytes(RandomBytes(32))
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}
