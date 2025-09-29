package cursor

import (
	"github.com/dmsRosa6/MoooChain/internal/blockchain"
	"github.com/redis/go-redis/v9"
)

const(
	PageSize = 500
)

type Cursor struct{
	blocks []blockchain.Block
	lashHash []byte
	height int
	hasNext bool
	page int
	client redis.Client
}

func NewCursor(c redis.Client) *Cursor{
	return &Cursor{
		client: c,
		hasNext: true,
	}
}

func (c *Cursor) Next()