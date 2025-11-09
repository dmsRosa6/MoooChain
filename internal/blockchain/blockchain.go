package blockchain

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"log"

	"github.com/dmsRosa6/MoooChain/internal/options"
	"github.com/dmsRosa6/MoooChain/internal/redisutils"
	"github.com/dmsRosa6/MoooChain/internal/transaction"
	"github.com/redis/go-redis/v9"
)

var (
	ErrBlockchainNotFound = errors.New("blockchain does not exist")
	GenesisData = "Genesis"
)

type Blockchain struct {
	LastHash []byte
	Database *redis.Client
	log      *log.Logger
	options  *options.Options
}

func InitBlockchain(r *redis.Client, log *log.Logger, options *options.Options, addr string) (*Blockchain, error) {
	ctx := context.Background()

	val, err := getBytes(r, ctx, redisutils.LastHashKeyKeyword)
	if err != nil {
		return nil, err
	}

	if options.DebugChain {
		_ = r.Eval(ctx, redisutils.InitDebugChainRedisFunction, []string{}, []string{})
	}

	bc := Blockchain{Database: r, log: log, options: options}

	if val == nil {
		log.Println("no blockchain found. creating new one...")
		var encoded bytes.Buffer
		encoder := gob.NewEncoder(&encoded)
		
		tx, err := transaction.CreateMintTx(addr, GenesisData)
		
		if err != nil {
			return nil, err
		}
		
		b := GenesisBlock(tx)
		jsonEncoded, err := json.Marshal(b)

		if err != nil {
			return nil, err
		}

		err = encoder.Encode(b)
	
		if err != nil {
			return nil, err
		}

		if err := setBytes(r, ctx, redisutils.BuildBlockKey(b.Hash), encoded.Bytes()); err != nil {
			return nil, err
		}

		if err := setString(r, ctx, redisutils.LastHashKeyKeyword, string(b.Hash)); err != nil {
			return nil, err
		}

		if err := setString(r, ctx, redisutils.BlockChainNameKeyword, redisutils.BlockChainName); err != nil {
			return nil, err
		}

		if options.DebugChain {
			if _, err := r.RPush(ctx, redisutils.BlockChainName, jsonEncoded).Result(); err != nil {
				return nil, err
			}
		}
	} else {
		log.Println("blockchain found.")
		bc.LastHash = val
	}

	return &bc, nil
}

func (bc *Blockchain) AddBlock(transactions []*transaction.Transaction) error {
	ctx := context.Background()

	lh, err := getBytes(bc.Database, ctx, redisutils.LastHashKeyKeyword)
	
	if err != nil {
		return err
	}
	
	if lh == nil {
		return ErrBlockchainNotFound
	}

	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	
	newBlock := CreateBlock(transactions, lh)
	data, err := json.Marshal(newBlock)
	
	if err != nil {
		return err
	}

	err = encoder.Encode(newBlock)

	if err != nil {
		return err
	}

	if err := setBytes(bc.Database, ctx, redisutils.BuildBlockKey(newBlock.Hash), data); err != nil {
		return err
	}

	if err := setBytes(bc.Database, ctx, redisutils.BuildPrevBlockKey(newBlock.Hash), newBlock.PrevHash); err != nil {
		return err
	}

	if err := setBytes(bc.Database, ctx, redisutils.LastHashKeyKeyword, newBlock.Hash); err != nil {
		return err
	}

	if bc.options.DebugChain {
		if _, err := bc.Database.RPush(ctx, redisutils.BlockChainName, string(data)).Result(); err != nil {
			return err
		}
	}

	bc.LastHash = newBlock.Hash

	return nil
}

func (bc *Blockchain) IterateBlockChain() (*BlockIterator, error) {
	ite, err := NewBlockIterator(bc.Database)

	if err != nil {
		return nil, err
	}

	return ite, nil
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
