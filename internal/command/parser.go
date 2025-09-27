package commands

import (
	"fmt"
	"strings"
)

type Command int

const(
	CreateBlockChain Command = iota
	AddBlock
	IterateBlockChain
	GetBlock
	Exit
)

var(
	Commands = map[string]Command{
		"CREATEBLOCKCHAIN" : CreateBlockChain,
		"ADDBLOCK" : AddBlock,
		"ITERATEBLOCKCHAIN" : IterateBlockChain,
		"GETBLOCK" : GetBlock,
		"INIT" : CreateBlockChain,
		"ADD" : AddBlock,
		"ITERATE" : IterateBlockChain,
		"GET" : GetBlock,
		"EXIT" : Exit,
	}

	CommandArgs = map[Command]int{
		CreateBlockChain : 0,
		AddBlock : 1,
		IterateBlockChain : 0,
		GetBlock : 1,
		Exit : 0,
	}

	CommandLongName = map[Command]string{
		CreateBlockChain : "Create BlockChain",
		AddBlock : "Add Block",
		IterateBlockChain : "Iterate BlockChain",
		GetBlock : "Get Block",
		Exit : "Exit",
	}
)


type Parser struct{

}

func NewParser() *Parser{
	return &Parser{}
}

func (p *Parser) Parse(s string) (Command, []string, error){
	tokens := strings.Split(s, " ")

	if len(tokens) == 0 {
		return -1, nil, fmt.Errorf("empty input")
	}

	commandAsString := strings.ToUpper(tokens[0])
	args := tokens[1:]


	command, ok := Commands[commandAsString]

	if !ok {
		return -1, nil, fmt.Errorf("command %s is not a valid command", commandAsString)
	}

	numberOfArgs := CommandArgs[command]

	if len(args) != numberOfArgs {
		return -1, nil, 
		fmt.Errorf("invalid number of arguments for %s operation. got: %d, expected: %d", 
			CommandLongName[command], 
			len(args),
			numberOfArgs,
		)
	}


	return command, args, nil
}