package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"github.com/dmsRosa6/MoooChain/types"
)

type PrivKey struct{
	key *ecdsa.PrivateKey
}

func (k PrivKey)Sign(data []byte) (*Signature, error){

	//TODO I should probably use SignASN1 to simplify
	r,s, err := ecdsa.Sign(rand.Reader, k.key, data)

	if err != nil {
		return nil, errors.New("error while signing")
	}

	return &Signature{r: r, s: s}, nil
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
	r *big.Int
	s *big.Int
}

//TODO change this to a error instead of a bool maybe
func (sig Signature) Verify(signedData []byte, pubKey PubKey) bool{

	return ecdsa.Verify(pubKey.key, signedData, sig.r, sig.s)
}