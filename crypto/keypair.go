package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/dmsRosa6/MoooChain/types"
)

type PrivKey struct{
	key *ecdsa.PrivateKey
}


func NewPrivKey() PrivKey{
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	
	if err != nil {
		msg := fmt.Sprintf("error generating private key: %s", err.Error())
		panic(msg)
	}

	return PrivKey{key: key}
}

func (k PrivKey)PubKey() PubKey{

	return PubKey{key: &k.key.PublicKey}
}

type PubKey struct{
	key *ecdsa.PublicKey
}

func (k PubKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.key, k.key.X, k.key.Y)	
}

func (k PubKey) Address() types.Address {
	hash := sha256.Sum256(k.ToSlice())
	var addr types.Address
	copy(addr[:], hash[len(hash)-20:])
	return addr
}


type Signature struct{
	
}