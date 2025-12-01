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