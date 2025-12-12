package core

import "errors"

type Validator interface {
    Validate(any) error
}

type BlockValidator struct {
    bc Blockchain
}

func NewBlockValidator(bc Blockchain) *BlockValidator {
    return &BlockValidator{bc}
}

//TODO Incomplete
func (bv *BlockValidator) Validate(v any) error {
    _, ok := v.(Block)
    if !ok {
        return errors.New("expected Block")
    }
    return nil
}
