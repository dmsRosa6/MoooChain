package blockchain

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const (
	defaultCapacity = 50
)

type BlockIterator struct {
	blocks      []Block
	nextHash    []byte
	index       int
	capacity    int
	HasNextPage bool
	redis       *redis.Client
}

func NewBlockIteratorWithCapacity(c *redis.Client, capacity int) (*BlockIterator, error) {
	blockIterator := &BlockIterator{redis: c, nextHash: []byte(""), capacity: capacity}
	err := blockIterator.populate()

	if err != nil {
		return nil, err
	}

	return blockIterator, nil
}

func NewBlockIterator(c *redis.Client) (*BlockIterator, error) {
	blockIterator := &BlockIterator{redis: c, nextHash: []byte(""), capacity: defaultCapacity}
	err := blockIterator.populate()

	if err != nil {
		return nil, err
	}

	return blockIterator, nil
}

func (c *BlockIterator) verifyIfPopulateNeeded() error {
	if c.index < len(c.blocks) && c.HasNextPage {
		return c.populate()
	}

	return nil
}

func (c *BlockIterator) HasNext() bool {
	return c.index < len(c.blocks) || c.HasNextPage
}

func (c *BlockIterator) Next() Block {
	c.verifyIfPopulateNeeded()

	if c.index < len(c.blocks) {
		return Block{}
	}

	block := c.blocks[c.index]
	c.index++
	return block
}

func (c *BlockIterator) NextRange(num int) []Block {
	var blocks []Block
	if c.index+num > len(c.blocks) {
		blocks = c.blocks[c.index:]
		remain := num - len(blocks)
		c.index = len(c.blocks)
		c.verifyIfPopulateNeeded()
		if c.index < len(c.blocks) {
			return blocks
		}

		blocks = append(blocks, c.blocks[c.index:c.index+remain]...)
	} else {
		blocks = c.blocks[c.index : c.index+num]
		c.index += num
	}

	return blocks
}

func (c *BlockIterator) populate() error {
	ctx := context.Background()

	res, err := c.redis.FCall(ctx, "iterate_chain", []string{}, c.nextHash, strconv.Itoa(c.capacity)).Result()
	if err != nil {
		return err
	}

	nextHash, items, hasNext, err := parseIterateChainReply(res)
	if err != nil {
		return err
	}

	c.nextHash = nextHash
	c.blocks = items
	c.HasNextPage = hasNext

	return nil
}

func parseIterateChainReply(reply any) (nextHash []byte, items []Block, more bool, err error) {
	arr, ok := reply.([]any)
	if !ok {
		err = errors.New("unexpected reply type, expected array")
		return
	}
	if len(arr) < 3 {
		err = fmt.Errorf("unexpected array length: %d", len(arr))
		return
	}

	switch v := arr[0].(type) {
	case string:
		if v != "" {
			nextHash, err = hex.DecodeString(v)
			if err != nil {
				err = fmt.Errorf("invalid hex in nextHash: %w", err)
				return
			}
		}
	case []byte:
		if len(v) > 0 {
			nextHash = append([]byte(nil), v...)
		}
	case nil:
	default:
		s := fmt.Sprintf("%v", v)
		if s != "" {
			nextHash, err = hex.DecodeString(s)
			if err != nil {
				err = fmt.Errorf("invalid hex in nextHash (fallback): %w", err)
				return
			}
		}
	}

	resultArr, ok := arr[1].([]interface{})
	if !ok {
		if arr[1] == nil {
			items = []Block{}
		} else {
			err = errors.New("unexpected type for results array")
			return
		}
	} else {
		items = make([]Block, 0, len(resultArr))
		for i := range resultArr {
			var raw []byte

			switch it := resultArr[i].(type) {
			case string:
				raw = []byte(it)
			case []byte:
				raw = it
			case nil:
				continue
			default:
				raw = []byte(fmt.Sprintf("%v", it))
			}

			if len(raw) == 0 {
				continue
			}

			var blk Block
			if e := json.Unmarshal(raw, &blk); e != nil {
				err = fmt.Errorf("failed to unmarshal block[%d]: %w", i, e)
				return
			}

			items = append(items, blk)
		}
	}

	switch v := arr[2].(type) {
	case int64:
		more = v != 0
	case int:
		more = v != 0
	case string:
		more = v != "0" && v != ""
	case []byte:
		more = string(v) != "0" && len(v) > 0
	default:
		more = false
	}

	return
}
