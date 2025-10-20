package transaction

import "fmt"

type Transaction struct {
	ID []byte
	Inputs []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value string
	PubKey string
}

type TxInput struct{
	ID []byte
	Out int
	Sig string
}


func BaseTx(to, data string) (*Transaction, error){
	if data == "" {
		data = fmt.Sprintf("to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{"moo", to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	
	id, err := generate256HashWithTx(tx)

	if err != nil {
		return nil, err
	}

	tx.ID = id

	return &tx, nil
}

func (tx *Transaction) IsBaseTx() bool{
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (in *TxInput) CanUnlock(data string) bool{
	return in.Sig == data
}

func (out *TxOutput) CanBeLocked(data string) bool{
	return out.PubKey == data
}
