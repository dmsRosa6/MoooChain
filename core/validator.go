package core

import (
	"errors"
	"fmt"
)

type Validator interface {
    Validate(any) error
}

type BlockValidator struct {
    bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
    return &BlockValidator{bc}
}

//TODO Incomplete
func (bv *BlockValidator) Validate(v any) error {
    block, ok := v.(*Block)
    if !ok {
        return errors.New("expected pointer to Block")
    }

    if bv.bc.HasBlock(block.Height) {
        return fmt.Errorf("the block with height (%d) already exists", block.Height)    
    }

    if err := block.Verify(); err != nil {
        return err
    }

    return nil
}
