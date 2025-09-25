package blockchain

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const (
	LastHashKeyKeyword          = "LastHash"
	BlockChainNameKeyword       = "BlockChainName"
	DebugChainKeyword           = "DebugChain"
	BlockChainName              = "Moochain"
	BlockKeyword                = "Block:"
	GenesisBlockKeyword         = "GenesisBlockHash"
	InitDebugChainRedisFunction = "init_debug_chain"
	IterateChainRedisFunction   = "iterate_chain"
	PrevBlockKeyword            = "Block:prev:"
)

var (
	ErrBlockchainNotFound = errors.New("blockchain does not exist")
	DebugChain            bool
)

func init() {
	DebugChain = false

	val := os.Getenv("DEBUG_CHAIN")

	if val != "" {
		convertedVal, err := strconv.ParseBool(val)

		if err != nil {
			log.Printf("Invalid DEBUG_CAHIN value %q, defaulting to FALSE", val)
		} else {
			DebugChain = convertedVal
		}
	}

}

type Blockchain struct {
	LastHash []byte
	Database *redis.Client
	log      *log.Logger
}

func InitBlockchain() (*Blockchain, error) {

	log := configLog()
	redis := initRedis()

	ctx := context.Background()

	val, err := redis.Get(ctx, LastHashKeyKeyword).Bytes()

	if err != nil {
		return nil, err
	}

	if DebugChain {
		redis.Eval(ctx, InitDebugChainRedisFunction, []string{}, []string{})
	}

	bc := Blockchain{Database: redis, log: log}

	if val == nil {
		log.Println("no blockchain found. creating new one...")

		b := GenesisBlock()
		data, err := json.Marshal(b)

		if err != nil {
			return nil, err
		}

		key := hex.EncodeToString(b.Hash)
		_, err = bc.Database.Set(ctx, buildBlockKey(key), data, 0).Result()

		if err != nil {
			return nil, err
		}

		_, err = bc.Database.Set(ctx, LastHashKeyKeyword, b.Hash, 0).Result()

		if err != nil {
			return nil, err
		}

		_, err = bc.Database.Set(ctx, BlockChainNameKeyword, BlockChainName, 0).Result()

		if err != nil {
			return nil, err
		}

		if DebugChain {

			_, err = bc.Database.RPush(ctx, string(b.PrevHash)+":"+string(b.Hash)+":"+string(b.Data), BlockChainName).Result()

			if err != nil {
				return nil, err
			}
		}

	} else {
		log.Println("blockchain found.")
		bc.LastHash = []byte(val)
	}

	return &bc, nil
}

func (bc *Blockchain) AddBlock(blockData string) error {
	ctx := context.Background()

	lh, err := bc.Database.Get(ctx, LastHashKeyKeyword).Bytes()

	if err != nil {
		return err
	}

	if lh == nil {
		return ErrBlockchainNotFound
	}

	newBlock := CreateBlock(blockData, lh)

	data, err := json.Marshal(newBlock)

	if err != nil {
		return err
	}

	key := hex.EncodeToString(newBlock.Hash)
	_, err = bc.Database.Set(ctx, buildBlockKey(key), data, 0).Result()

	if err != nil {
		return err
	}

	_, err = bc.Database.Set(ctx, buildPrevBlockKey(key), lh, 0).Result()

	if err != nil {
		return err
	}

	_, err = bc.Database.Set(ctx, LastHashKeyKeyword, newBlock.Hash, 0).Result()

	if err != nil {
		return err
	}

	_, err = bc.Database.Set(ctx, LastHashKeyKeyword, newBlock.Hash, 0).Result()

	if err != nil {
		return err
	}

	if DebugChain {
		_, err = bc.Database.RPush(ctx, string(newBlock.PrevHash)+":"+string(newBlock.Hash)+":"+string(newBlock.Data), BlockChainName).Result()

		if err != nil {
			return err
		}
	}

	return nil
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

func configLog() *log.Logger {
	return log.New(os.Stdout, "Moochain:", log.LstdFlags)

}

func buildBlockKey(hash string) string {
	return BlockKeyword + hash
}

func buildPrevBlockKey(prev string) string {
	return PrevBlockKeyword + prev
}
