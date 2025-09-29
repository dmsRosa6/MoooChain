package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/dmsRosa6/MoooChain/internal/blockchain"
	"github.com/redis/go-redis/v9"
)

type Executer struct {
	log *log.Logger
	blockchain *blockchain.Blockchain
}

func NewExecuter(l *log.Logger) *Executer {
    return &Executer{log:l}
}

func (e *Executer) Execute(command Command, args []string) error {
    switch command {
    case CreateBlockChain:
        log.Println("Executing:", CommandLongName[command])
		redis := initRedis()
		log.Println("Redis initialized")
		bc, err := blockchain.InitBlockchain(redis,e.log)
		
		if err != nil {
			return err
		}

		e.blockchain = bc
		
		return nil

    case AddBlock:
        log.Println("Executing:", CommandLongName[command], "with args:", args)

		if e.blockchain == nil {
			return fmt.Errorf("blockchain not initialized")
		}

		err := e.blockchain.AddBlock(args[0])

		if err != nil {
			return err
		}

        return fmt.Errorf("not implemented")

    case IterateBlockChain:
        log.Println("Executing:", CommandLongName[command])

		if e.blockchain == nil {
			return fmt.Errorf("blockchain not initialized")
		}


        return fmt.Errorf("not implemented")

    case GetBlock:
        log.Println("Executing:", CommandLongName[command], "with args:", args)
		
		if e.blockchain == nil {
			return fmt.Errorf("blockchain not initialized")
		}
        
		return fmt.Errorf("not implemented")

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

func initRedis() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr: buildAddr(),
	})
	return client
}

func buildAddr() string {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	return host + ":" + port
}
