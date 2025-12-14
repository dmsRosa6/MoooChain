package commands

import (
	"fmt"
	"log"

	"github.com/dmsRosa6/MoooChain/core"
	"github.com/dmsRosa6/MoooChain/options"
	"github.com/redis/go-redis/v9"
)

type Executer struct {
	log        *log.Logger
	blockchain *core.Blockchain
	redis      *redis.Client
	options    *options.Options
}

func NewExecuter(l *log.Logger, options *options.Options, redis *redis.Client) *Executer {
	log.Println("Redis initialized")
	return &Executer{log: l, options: options, redis: redis}
}

func (e *Executer) Execute(command Command, args []string) error {
	switch command {
	case CreateBlockChain:
		e.log.Println("Executing:", CommandLongName[command])
	
		return nil

	case AddBlock:
		//implement
		return nil

	case Send:

		return nil

	case GetBlock:
		e.log.Println("Executing:", CommandLongName[command], "with args:", args)

		if e.blockchain == nil {
			return fmt.Errorf("blockchain not initialized")
		}

		return fmt.Errorf("not implemented")

	case DestroyBlockChain:
		
	case Exit:

	default:
		return fmt.Errorf("unknown command: %v", command)
	}

	return nil
}