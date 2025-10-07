package blockchain

import (
	"context"
	"encoding/hex"
	"errors"
	"log"

	"github.com/dmsRosa6/MoooChain/internal/options"
	"github.com/redis/go-redis/v9"
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
	GetBlockByHeightRedisFunction = ""
)

var (
	ErrBlockchainNotFound = errors.New("blockchain does not exist")
)


type Blockchain struct {
	LastHash []byte
	Database *redis.Client
	log      *log.Logger
	options *options.Options
}

func InitBlockchain(r *redis.Client, log *log.Logger, options *options.Options) (*Blockchain, error) {
	ctx := context.Background()

	val, err := getBytes(r, ctx, LastHashKeyKeyword)
	if err != nil {
		return nil, err
	}

	if options.DebugChain {
		r.Eval(ctx, InitDebugChainRedisFunction, []string{}, []string{})
	}

	bc := Blockchain{Database: r, log: log, options: options}

	if val == nil {
		log.Println("no blockchain found. creating new one...")

		b := GenesisBlock()
		data, err := b.ToJSON()
		if err != nil {
			return nil, err
		}

		key := hex.EncodeToString(b.Hash)

		if err := setKey(r, ctx, buildBlockKey(key), data); err != nil {
			return nil, err
		}

		if err := setKey(r, ctx, LastHashKeyKeyword, b.Hash); err != nil {
			return nil, err
		}

		if err := setKey(r, ctx, BlockChainNameKeyword, BlockChainName); err != nil {
			return nil, err
		}

		if options.DebugChain {
			if _, err := r.RPush(ctx, BlockChainName, data).Result(); err != nil {
				return nil, err
			}
		}
	} else {
		log.Println("blockchain found.")
		bc.LastHash = val
	}

	return &bc, nil
}

func (bc *Blockchain) AddBlock(blockData string) error {
	ctx := context.Background()

	lh, err := getBytes(bc.Database, ctx, LastHashKeyKeyword)
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

	if err := setKey(bc.Database, ctx, buildBlockKey(key), data); err != nil {
		return err
	}
	if err := setKey(bc.Database, ctx, buildPrevBlockKey(key), lh); err != nil {
		return err
	}
	if err := setKey(bc.Database, ctx, LastHashKeyKeyword, newBlock.Hash); err != nil {
		return err
	}

	if bc.options.DebugChain {
		if _, err := bc.Database.RPush(ctx, BlockChainName, data).Result(); err != nil {
			return err
		}
	}

	return nil
}

func (bc *Blockchain) GetBlockByHeight(height int){
	context := context.Background()

	bc.Database.Eval(context,,BlockChainName,nil)	
}

func buildBlockKey(hash string) string {
	return BlockKeyword + hash
}

func buildPrevBlockKey(prev string) string {
	return PrevBlockKeyword + prev
}

func getBytes(r *redis.Client, ctx context.Context, key string) ([]byte, error) {
	val, err := r.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func setKey(r *redis.Client, ctx context.Context, key string, value any) error {
	_, err := r.Set(ctx, key, value, 0).Result()
	return err
}
