package blockchain

import (
	"context"
	"errors"
	"log"

	"github.com/dmsRosa6/MoooChain/internal/codec"
	"github.com/dmsRosa6/MoooChain/internal/options"
	"github.com/redis/go-redis/v9"
)

const (
	LastHashKeyKeyword          = "LastHash"
	BlockChainNameKeyword       = "BlockChainName"
	DebugChainKeyword           = "DebugChain"
	BlockChainName              = "Moochain"
	GenesisBlockKeyword         = "GenesisBlockHash"
	InitDebugChainRedisFunction = "init_debug_chain"
	IterateChainRedisFunction   = "iterate_chain"
	GetBlockByHeightRedisFunction = ""
)

var (
	ErrBlockchainNotFound = errors.New("blockchain does not exist")
)

type Blockchain struct {
	LastHash []byte
	Database *redis.Client
	log      *log.Logger
	options  *options.Options
}

func InitBlockchain(r *redis.Client, log *log.Logger, options *options.Options) (*Blockchain, error) {
	ctx := context.Background()

	val, err := getBytes(r, ctx, LastHashKeyKeyword)
	if err != nil {
		return nil, err
	}

	if options.DebugChain {
		_ = r.Eval(ctx, InitDebugChainRedisFunction, []string{}, []string{})
	}

	bc := Blockchain{Database: r, log: log, options: options}

	if val == nil {
		log.Println("no blockchain found. creating new one...")

		b := GenesisBlock()
		data, err := b.MarshalJSON()
		if err != nil {
			return nil, err
		}

		if err := setBytes(r, ctx, codec.BuildBlockKey(b.Hash), data); err != nil {
			return nil, err
		}

		if err := setString(r, ctx, LastHashKeyKeyword, string(b.Hash)); err != nil {
			return nil, err
		}

		if err := setString(r, ctx, BlockChainNameKeyword, BlockChainName); err != nil {
			return nil, err
		}

		if options.DebugChain {
			if _, err := r.RPush(ctx, BlockChainName, string(data)).Result(); err != nil {
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
	data, err := newBlock.MarshalJSON()
	if err != nil {
		return err
	}

	if err := setBytes(bc.Database, ctx, codec.BuildBlockKey(newBlock.Hash), data); err != nil {
		return err
	}

	if err := setBytes(bc.Database, ctx, codec.BuildPrevBlockKey(newBlock.Hash), newBlock.PrevHash); err != nil {
		return err
	}

	if err := setBytes(bc.Database, ctx, LastHashKeyKeyword, newBlock.Hash); err != nil {
		return err
	}

	if bc.options.DebugChain {
		if _, err := bc.Database.RPush(ctx, BlockChainName, string(data)).Result(); err != nil {
			return err
		}
	}

	bc.LastHash = newBlock.Hash

	return nil
}

func getBytes(r *redis.Client, ctx context.Context, key string) ([]byte, error) {
	val, err := r.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func setBytes(r *redis.Client, ctx context.Context, key string, value []byte) error {
	return r.Set(ctx, key, value, 0).Err()
}

func setString(r *redis.Client, ctx context.Context, key string, value string) error {
	return r.Set(ctx, key, value, 0).Err()
}
