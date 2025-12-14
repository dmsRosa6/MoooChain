package core

import (
	"context"
	"errors"

	"github.com/dmsRosa6/MoooChain/log"
	"github.com/dmsRosa6/MoooChain/options"
	"github.com/dmsRosa6/MoooChain/store"
)

var (
	ErrBlockchainNotFound = errors.New("blockchain does not exist")
	GenesisData           = "Genesis"
)

type Blockchain struct {
	headers []*Header
	store store.Store
	validator Validator
	log      log.Logger
	options  *options.Options
}

func (bc *Blockchain) Height() uint64{
	return uint64(len(bc.headers) - 1)
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) HasBlock(height uint64) bool{
	return height < bc.Height()
}

func (bc *Blockchain) addBlockNoValidation(b * Block) error{
	bc.headers = append(bc.headers, b.Header)
    ctx := context.Background()

	err := bc.store.Set(ctx, "Genesis",b)

	return  err
}

func NewBlockchain(genesis *Block, r store.Store, log log.Logger, options *options.Options) (*Blockchain, error) {
	bc := &Blockchain{store: r, log: log, options: options}

	blockValidator := NewBlockValidator(bc) 
	bc.validator = blockValidator

	err := bc.addBlockNoValidation(genesis)


	return bc, err
}

func (bc *Blockchain) AddBlock(block *Block) error {

	if err := bc.validator.Validate(block); err != nil{
		return err
	}

	return bc.addBlockNoValidation(block)
}