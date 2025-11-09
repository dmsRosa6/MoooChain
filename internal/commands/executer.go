package commands

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/dmsRosa6/MoooChain/internal/blockchain"
	"github.com/dmsRosa6/MoooChain/internal/options"
	"github.com/redis/go-redis/v9"
)

type Executer struct {
	log *log.Logger
	blockchain *blockchain.Blockchain
	redis *redis.Client
	options *options.Options
}

func NewExecuter(l *log.Logger, options *options.Options,  redis *redis.Client) *Executer {
	log.Println("Redis initialized")
    return &Executer{log:l, options: options, redis: redis}
}

func (e *Executer) Execute(command Command, args []string) error {
    switch command {
    case CreateBlockChain:
        e.log.Println("Executing:", CommandLongName[command])
		addr := args[0]
		bc, err := blockchain.InitBlockchain(e.redis,e.log, e.options, addr)
		
		if err != nil {
			return err
		}

		e.blockchain = bc
		
		return nil

    case AddBlock:
       	//implement
        return nil 
	
	case Send:
       	
		
		return nil 

    case IterateBlockChain:
        e.log.Println("Executing:", CommandLongName[command])

		ite, _ := e.blockchain.IterateBlockChain();

		for ite.HasNext() {
			fmt.Println(ite.Next())
		}

        return nil

    case GetBlock:
        e.log.Println("Executing:", CommandLongName[command], "with args:", args)
		
		if e.blockchain == nil {
			return fmt.Errorf("blockchain not initialized")
		}
        
		return fmt.Errorf("not implemented")

	case DestroyBlockChain:
		ctx := context.Background()
		e.redis.FlushAll(ctx);
		e.log.Println("Deleting...")
		return nil
    case Exit:
		if e.blockchain == nil{
			return nil
		}

		err := e.blockchain.Database.Close()
		if err != nil{
			return err
		}
        
		log.Println("Exiting...")
        return nil

    default:
        return fmt.Errorf("unknown command: %v", command)
    }
}

func (e *Executer) CleanupChain() error{
	if e.blockchain == nil {
		return errors.New("blockchain null on chain cleanup operation")
	}

	if e.blockchain.Database == nil {
		return errors.New("db client null on chain cleanup operation")
	}
	ctx := context.Background()

	e.blockchain.Database.FlushAll(ctx)
	
	return nil 
}
