package redisutils

import "encoding/hex"

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
