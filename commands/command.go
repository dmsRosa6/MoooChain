package commands

import "strings"

type Command int

const (
	InvalidCommand   Command = -1
	CreateBlockChain Command = iota
	DestroyBlockChain
	AddBlock
	IterateBlockChain
	GetBlock
	Send
	Exit
)

type CommandInfo struct {
	Names    []string
	NumArgs  int
	LongName string
}

var (
	StringToCommand   = make(map[string]Command)
	CommandNumberArgs = make(map[Command]int)
	CommandLongName   = make(map[Command]string)
)

func registerCommand(cmd Command, info CommandInfo) {
	for _, name := range info.Names {
		StringToCommand[strings.ToUpper(name)] = cmd
	}
	CommandNumberArgs[cmd] = info.NumArgs
	CommandLongName[cmd] = info.LongName
}

func init() {
	registerCommand(CreateBlockChain, CommandInfo{
		Names:    []string{"CREATEBLOCKCHAIN", "INIT"},
		NumArgs:  1,
		LongName: "Create BlockChain",
	})
	registerCommand(AddBlock, CommandInfo{
		Names:    []string{"ADDBLOCK", "ADD"},
		NumArgs:  1,
		LongName: "Add Block",
	})
	registerCommand(IterateBlockChain, CommandInfo{
		Names:    []string{"ITERATEBLOCKCHAIN", "ITERATE"},
		NumArgs:  0,
		LongName: "Iterate BlockChain",
	})
	registerCommand(GetBlock, CommandInfo{
		Names:    []string{"GETBLOCK", "GET"},
		NumArgs:  1,
		LongName: "Get Block",
	})
	registerCommand(DestroyBlockChain, CommandInfo{
		Names:    []string{"DESTROYBLOCKCHAIN", "DESTROY"},
		NumArgs:  0,
		LongName: "Destroy BlockChain",
	})
	registerCommand(Send, CommandInfo{
		Names:    []string{"SEND"},
		NumArgs:  3,
		LongName: "Send",
	})
	registerCommand(Exit, CommandInfo{
		Names:    []string{"EXIT"},
		NumArgs:  0,
		LongName: "Exit",
	})
}
