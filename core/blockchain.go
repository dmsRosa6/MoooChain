package core

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"log"

	"github.com/dmsRosa6/MoooChain/options"
	"github.com/dmsRosa6/MoooChain/utils"
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

	val, err := getBytes(r, ctx, utils.LastHashKeyKeyword)
	if err != nil {
		return nil, err
	}

	if options.DebugChain {
		_ = r.Eval(ctx, utils.InitDebugChainRedisFunction, []string{}, []string{})
	}

	bc := Blockchain{Database: r, log: log, options: options}

	if val == nil {
		log.Println("no blockchain found. creating new one...")
		var encoded bytes.Buffer
		encoder := gob.NewEncoder(&encoded)
		
		tx, err := CreateMintTx(addr, GenesisData)
		
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

		if err := setBytes(r, ctx, utils.BuildBlockKey(b.Hash), encoded.Bytes()); err != nil {
			return nil, err
		}

		if err := setString(r, ctx, utils.LastHashKeyKeyword, string(b.Hash)); err != nil {
			return nil, err
		}

		if err := setString(r, ctx, utils.BlockChainNameKeyword, utils.BlockChainName); err != nil {
			return nil, err
		}

		if options.DebugChain {
			if _, err := r.RPush(ctx, utils.BlockChainName, jsonEncoded).Result(); err != nil {
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

	lh, err := getBytes(bc.Database, ctx, utils.LastHashKeyKeyword)
	
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

	if err := setBytes(bc.Database, ctx, utils.BuildBlockKey(newBlock.Hash), data); err != nil {
		return err
	}

	if err := setBytes(bc.Database, ctx, utils.BuildPrevBlockKey(newBlock.Hash), newBlock.PrevHash); err != nil {
		return err
	}

	if err := setBytes(bc.Database, ctx, utils.LastHashKeyKeyword, newBlock.Hash); err != nil {
		return err
	}

	if bc.options.DebugChain {
		if _, err := bc.Database.RPush(ctx, utils.BlockChainName, string(data)).Result(); err != nil {
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
