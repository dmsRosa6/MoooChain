package blockchain

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

const(
	LAST_HASH = "lh"
)

type Blockchain struct{
	LastHash []byte
	Database *redis.Client
	log		 *log.Logger
}

func InitBlockchain() (*Blockchain, error) {
	
	log := configLog() 
	b := GenesisBlock()
	bc := Blockchain{ b.Hash , initRedis(), log}
	ctx := context.Background()
	val, err := bc.Database.Get(ctx, LAST_HASH).Bytes()

	if err != nil {
		return nil, err		
	}

	if val == nil {
		log.Println("no blockchain found. creating new one...")
		data, err := json.Marshal(b)
		
		if err != nil {
			return nil, err
		}
		
		key := hex.EncodeToString(b.Hash)
		_, err = bc.Database.Set(ctx, key, data, 0).Result()

		if err != nil {
			return nil, err
		}

		_, err = bc.Database.Set(ctx, LAST_HASH, b.Hash, 0).Result()

		if err != nil {
			return nil, err
		}
	}else{
		log.Println("blockchain found.")
		bc.LastHash = []byte(val)
	}
	
	return &bc, nil 
}

func (bc *Blockchain) AddBlock(data string) {
	newBlock := CreateBlock(data,bc.LastHash)

	bc.LastHash = newBlock.Hash
}

func initRedis() *redis.Client{
	
	client := redis.NewClient(&redis.Options{
		Addr: buildAddr(),
	})
	return client
}

func buildAddr() string{
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	return host + ":" + port 
}

func configLog() *log.Logger{
	return log.New(os.Stdout, "Moochain:", log.LstdFlags)
	
}