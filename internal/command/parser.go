package commands

import (
	"fmt"
	"strings"
)

type Command int

const(
	InvalidCommand Command = -1
	CreateBlockChain Command = iota
	AddBlock
	IterateBlockChain
	GetBlock
	Exit
)

var(
	StringToCommand = map[string]Command{
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

	CommandNumberArgs = map[Command]int{
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
	s = strings.TrimSpace(s)
	
	if s == "" {
		return InvalidCommand, nil, fmt.Errorf("empty input")
	}

	tokens := strings.Fields(s)

	commandAsString := strings.ToUpper(strings.TrimSpace(tokens[0]))
	args := tokens[1:]


	command, ok := StringToCommand[commandAsString]

	if !ok {
		return InvalidCommand, nil, fmt.Errorf("command is not valid: %s", commandAsString)
	}

	numberOfArgs := CommandNumberArgs[command]

	if len(args) != numberOfArgs {
		return InvalidCommand, nil, 
		fmt.Errorf("invalid number of arguments for %s operation. got: %d, expected: %d", 
			CommandLongName[command], 
			len(args),
			numberOfArgs,
		)
	}


	return command, args, nil
}