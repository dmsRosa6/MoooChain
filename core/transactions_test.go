package core

import (
	"testing"

	"github.com/dmsRosa6/MoooChain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSignTransactionSuccess(t *testing.T) {
	privKey := crypto.NewPrivKey()

	tx, err := CreateMintTx("Foo", "Bar")

	assert.Nil(t, err)
	assert.NotNil(t, tx)

	err = tx.Sign(privKey)

	assert.Nil(t, err)

	assert.NotNil(t, tx.Signature)
	assert.NotNil(t, tx.PubKey)
}

func TestSignAndVerifyTransactionSuccess(t *testing.T) {
	privKey := crypto.NewPrivKey()

	tx, err := CreateMintTx("Foo", "Bar")

	assert.Nil(t, err)
	assert.NotNil(t, tx)

	err = tx.Sign(privKey)

	assert.Nil(t, err)

	assert.NotNil(t, tx.Signature)
	assert.NotNil(t, tx.PubKey)

	assert.Nil(t, tx.Verify())
}

func TestSignAndVerifyTransactionError(t *testing.T) {
	privKey := crypto.NewPrivKey()

	tx, err := CreateMintTx("Foo", "Bar")

	assert.Nil(t, err)
	assert.NotNil(t, tx)

	err = tx.Sign(privKey)

	assert.Nil(t, err)

	assert.NotNil(t, tx.Signature)
	assert.NotNil(t, tx.PubKey)

	tx.PubKey = crypto.NewPrivKey().PubKey()

	assert.NotNil(t, tx.Verify())
}


