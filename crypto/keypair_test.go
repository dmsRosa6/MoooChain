package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestGeneratePrivateKey(t *testing.T) {

	var priv PrivKey

	assert.NotPanics(t, func() {priv = NewPrivKey()})

	assert.NotNil(t, priv)
}

func TestSignVerifySuccess(t *testing.T) {

	var priv PrivKey

	assert.NotPanics(t, func() {priv = NewPrivKey()})

	assert.NotNil(t, priv)

	msg := []byte("Hello")

	sig, err := priv.Sign(msg)

	assert.Nil(t, err)


	pub := priv.PubKey()

	assert.NotNil(t, pub)

	isVerified := sig.Verify(msg, pub)

	assert.True(t, isVerified)
}


func TestSignVerifyFail(t *testing.T) {

	var priv1, priv2 PrivKey
	msg := []byte("Hello")

	assert.NotPanics(t, func() {priv1 = NewPrivKey()})
	assert.NotNil(t, priv1)


	sig, err := priv1.Sign(msg)
	assert.Nil(t, err)

	pub := priv1.PubKey()
	assert.NotNil(t, pub)

	assert.NotPanics(t, func() {priv2 = NewPrivKey()})
	assert.NotNil(t, priv2)
	assert.Nil(t, err)

	pub2 := priv2.PubKey()
	assert.NotNil(t, pub2)

	//create 2 keys sign with one and verify with another
	isVerified := sig.Verify(msg, 	pub2)

	other := sig.Verify([]byte("Some other message"), pub2)

	assert.False(t, isVerified)
	assert.False(t, other)
}

func TestGeneratePubKey(t *testing.T) {

	var priv PrivKey

	assert.NotPanics(t, func() {priv = NewPrivKey()})

	pub := priv.PubKey()

	assert.NotNil(t, priv)
	assert.NotNil(t, pub)
}

func TestEncryptDecrypt(t *testing.T) {

	var priv PrivKey

	assert.NotPanics(t, func() {priv = NewPrivKey()})

	pub := priv.PubKey()

	assert.NotNil(t, priv)
	assert.NotNil(t, pub)


}