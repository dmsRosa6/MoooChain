package core

import (
	"testing"

	"github.com/dmsRosa6/MoooChain/log"
	"github.com/dmsRosa6/MoooChain/options"
	"github.com/dmsRosa6/MoooChain/store"
	"github.com/stretchr/testify/assert"
)

func NewTestBlockchain() (*Blockchain, *Block, error) {
	genesis := RandomBlock(0,0)
	memStore := store.NewMemStore()
	logger := log.NewNopLogger()
	opts := options.InitOptions(logger)

	bc, err := NewBlockchain(genesis, memStore, logger, opts)
	if err != nil {
		return nil, nil, err
	}

	return bc, genesis, nil
}

func TestCreateBlockchain(t *testing.T) {
	bc, genesis, err := NewTestBlockchain()
	assert.NoError(t, err)
	assert.NotNil(t, bc)
	assert.NotNil(t, genesis)
	assert.Equal(t, bc.Height(), uint32(0))
	assert.False(t, bc.HasBlock(0))
}



func TestAddBlock(t * testing.T){
    bc, _, err := NewTestBlockchain()

	var numBlocks uint64 = 10

    for i := 1; i < int(numBlocks) + 1; i++ {
		b := RandomBlockWithSig(0,uint64(i))

		err = bc.AddBlock(b)

		assert.Nil(t, err)
	}

	assert.Equal(t, numBlocks, bc.Height())

	assert.Equal(t, int(numBlocks + 1), len(bc.headers))
	
	b := RandomBlockWithSig(0,2)

	err = bc.AddBlock(b)

	assert.NotNil(t, err)
}
