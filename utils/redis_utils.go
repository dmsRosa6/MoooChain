package utils

import (
	"encoding/hex"
	"fmt"
	"os"

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
	BlockKeyword                = "Block:"
	PrevBlockKeyword            = "Block:prev:"
)

func BuildBlockKey(hash []byte) string {
	return BlockKeyword + hex.EncodeToString(hash)
}

func BuildPrevBlockKey(hash []byte) string {
	return PrevBlockKeyword + hex.EncodeToString(hash)
}

func InitRedis() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr: buildAddr(),
	})
	return client
}

func buildAddr() string {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	if host == "" {
		host = "localhost"
		if port == "" {
		}
		port = "6379"
	}

	return fmt.Sprintf("%s:%s", host, port)
}
