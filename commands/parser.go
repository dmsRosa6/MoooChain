package commands

import (
	"fmt"
	"strings"
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