package blockchain

import (
	"github.com/go-redis/redis"
)

const(
	host = "localhost"
	port = "6379"
	password = ""
	db = 0
	protocol = 2
)

type Blockchain struct{
	LastHash []byte
	Database *redis.Client
}

func InitBlockchain() *Blockchain{
	b := GenesisBlock()
	bc := Blockchain{ b.Hash , initRedis()}
	return &bc
}

func (bc *Blockchain) AddBlock(data string) {
	newBlock := CreateBlock(data,bc.LastHash)

	bc.LastHash = newBlock.Hash
}

func initRedis() *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr: buildAddr(host, port),
		Password: password,
	})
	return client
}

func buildAddr(host string, port string) string{
	return host + ":" + port
}