package blockchain

import (
	"context"
	"encoding/hex"
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

func InitBlockchain(redis *redis.Client, log *log.Logger) (*Blockchain, error) {

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
		data, err := b.ToJSON()

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

			_, err = bc.Database.RPush(ctx, BlockChainName, data).Result()

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

	data, err := newBlock.ToJSON()

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


	if DebugChain {
		_, err = bc.Database.RPush(ctx, BlockChainName, data).Result()

		if err != nil {
			return err
		}
	}

	return nil
}

func buildBlockKey(hash string) string {
	return BlockKeyword + hash
}

func buildPrevBlockKey(prev string) string {
	return PrevBlockKeyword + prev
}
