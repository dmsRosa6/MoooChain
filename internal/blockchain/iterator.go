package blockchain

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const(
	defaultCapacity = 50
)

type BlockIterator struct{
	blocks []Block
	nextHash []byte
	index int
	client redis.Client
	capacity int
	HasNextPage bool
	redis *redis.Client
}

func NewBlockIteratorWithCapacity(c redis.Client, capacity int) *BlockIterator{
	blockIterator := &BlockIterator{client: c, nextHash: []byte(""), capacity: capacity}
	blockIterator.populate()

	return blockIterator
}


func NewBlockIterator(c redis.Client) *BlockIterator{
	blockIterator := &BlockIterator{client: c, nextHash: []byte(""),capacity: defaultCapacity}
	blockIterator.populate()

	return blockIterator
}

func (c *BlockIterator) verifyIfPopulateNeeded(){
	if c.index < len(c.blocks) && c.HasNextPage{
		c.populate()
	}
}

func (c *BlockIterator) Next() Block{
	c.verifyIfPopulateNeeded()

	if c.index < len(c.blocks) {
		return Block{}
	}

	block := c.blocks[c.index]
	c.index++
	return block
}

func (c *BlockIterator) NextRange(num int) []Block{
	var blocks []Block
	if c.index + num > len(c.blocks) {
		blocks = c.blocks[c.index:]
		remain := num - len(blocks)
		c.index = len(c.blocks) 
		c.verifyIfPopulateNeeded()
		if c.index < len(c.blocks) {
			return blocks
		}
		
		blocks = append(blocks, c.blocks[c.index:c.index+remain]...)
	}else{
		blocks = c.blocks[c.index:c.index+num]
		c.index += num
	}

	return blocks
}

func (c *BlockIterator) populate() error {
	ctx := context.Background()

	res, err := c.redis.Do(ctx, "FCALL", "iterate_chain", 0, c.nextHash, strconv.Itoa(c.capacity)).Result()
	if err != nil {
		return err
	}

	nextHash, items, hasNext, err := ParseIterateChainReply(res)
	if err != nil {
		return err
	}

	c.nextHash = nextHash
	c.blocks = items
	c.HasNextPage = hasNext

	return nil
}

