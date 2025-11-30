package types

import (
	"crypto/rand"
	"fmt"
)


type Hash [32]uint8


func HashFromBytes(arr []byte) Hash{

	if len(arr) != 32 {
		err := fmt.Sprintf("byte array length is %d, should be exactly 32", len(arr))
		panic(err)
	}


	var value [32]uint8

	copy(value[:], arr)
	
	return Hash(value)
}


func RandomBytes(size int) []byte{
	if size < 0 {
		msg := fmt.Sprintf("request size is %d, needed a value greater than 0", size)
		
		panic(msg)
	}

	token := make([]byte,size)
	rand.Read(token)

	return token
}

func RandomHash() Hash{
	return HashFromBytes(RandomBytes(32))
}


