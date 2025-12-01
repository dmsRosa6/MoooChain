package utils

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

type Iterator struct {
	object      []any
	nextKey     []byte
	index       int
	capacity    int
	HasNextPage bool
	redis       *redis.Client
}

func NewIteratorWithCapacity(c *redis.Client, capacity int) (*Iterator, error) {
	iterator := &Iterator{redis: c, nextKey: []byte(""), capacity: capacity}
	err := iterator.populate()

	if err != nil {
		return nil, err
	}

	return iterator, nil
}

func NewIterator(c *redis.Client) (*Iterator, error) {
	iterator := &Iterator{redis: c, nextKey: []byte(""), capacity: defaultCapacity}
	err := iterator.populate()

	if err != nil {
		return nil, err
	}

	return iterator, nil
}

func (c *Iterator) verifyIfPopulateNeeded() error {
	if c.index < len(c.object) && c.HasNextPage {
		return c.populate()
	}

	return nil
}

func (c *Iterator) HasNext() bool {
	return c.index < len(c.object) || c.HasNextPage
}

func (c *Iterator) Next() any {
	c.verifyIfPopulateNeeded()

	if c.index < len(c.object) {
		return nil
	}

	object := c.object[c.index]
	c.index++
	return object
}

func (c *Iterator) NextRange(num int) []any {
	var object []any
	if c.index+num > len(c.object) {
		object = c.object[c.index:]
		remain := num - len(object)
		c.index = len(c.object)
		c.verifyIfPopulateNeeded()
		if c.index < len(c.object) {
			return object
		}

		object = append(object, c.object[c.index:c.index+remain]...)
	} else {
		object = c.object[c.index : c.index+num]
		c.index += num
	}

	return object
}

func (c *Iterator) populate() error {
	ctx := context.Background()

	res, err := c.redis.FCall(ctx, "iterate_chain", []string{}, c.nextKey, strconv.Itoa(c.capacity)).Result()
	if err != nil {
		return err
	}

	nextKey, items, hasNext, err := parseIterateChainReply(res)
	if err != nil {
		return err
	}

	c.nextKey = nextKey
	c.object = items
	c.HasNextPage = hasNext

	return nil
}

func parseIterateChainReply(reply any) (nextKey []byte, items []any, more bool, err error) {
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
			nextKey, err = hex.DecodeString(v)
			if err != nil {
				err = fmt.Errorf("invalid hex in nextKey: %w", err)
				return
			}
		}
	case []byte:
		if len(v) > 0 {
			nextKey = append([]byte(nil), v...)
		}
	case nil:
	default:
		s := fmt.Sprintf("%v", v)
		if s != "" {
			nextKey, err = hex.DecodeString(s)
			if err != nil {
				err = fmt.Errorf("invalid hex in nextKey (fallback): %w", err)
				return
			}
		}
	}

	resultArr, ok := arr[1].([]interface{})
	if !ok {
		if arr[1] == nil {
			items = []any{}
		} else {
			err = errors.New("unexpected type for results array")
			return
		}
	} else {
		items = make([]any, 0, len(resultArr))
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

			var blk any
			if e := json.Unmarshal(raw, &blk); e != nil {
				err = fmt.Errorf("failed to unmarshal [%d]: %w", i, e)
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
