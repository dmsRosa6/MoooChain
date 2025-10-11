package cursor

import (
	"github.com/dmsRosa6/MoooChain/internal/blockchain"
	"github.com/redis/go-redis/v9"
)

const(
	defaultCapacity = 50
)

type Cursor struct{
	blocks []blockchain.Block
	lastHash []byte
	index int
	client redis.Client
	capacity int
	HasNextPage bool
}

func NewCursorWithCapacity(c redis.Client, capacity int) *Cursor{
	cursor := &Cursor{client: c, lastHash: []byte(""), capacity: capacity}
	//c.populate()

	return cursor
}


func NewCursor(c redis.Client) *Cursor{
	cursor := &Cursor{client: c, lastHash: []byte(""),capacity: defaultCapacity}
	//c.populate()

	return cursor
}

func (c *Cursor) verifyIfPopulateNeeded(){
	if c.index < len(c.blocks) && c.HasNextPage{
	//	c.populate()
	}
}

func (c *Cursor) Next() blockchain.Block{
	c.verifyIfPopulateNeeded()

	if c.index < len(c.blocks) {
		return blockchain.Block{}
	}

	block := c.blocks[c.index]
	c.index++
	return block
}

func (c *Cursor) NextRange(num int) []blockchain.Block{
	var blocks []blockchain.Block
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


