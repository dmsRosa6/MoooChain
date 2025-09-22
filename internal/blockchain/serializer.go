package blockchain

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/vmihailenco/msgpack"
)

var UseJSON bool


type Serializer interface {
    Marshal(v interface{}) ([]byte, error)
    Unmarshal(data []byte, v interface{}) error
}

type JSONSerializer struct{}
func (JSONSerializer) Marshal(v interface{}) ([]byte, error) { return json.Marshal(v) }
func (JSONSerializer) Unmarshal(data []byte, v interface{}) error { return json.Unmarshal(data, v) }

type MsgPackSerializer struct{}
func (MsgPackSerializer) Marshal(v interface{}) ([]byte, error) { return msgpack.Marshal(v) }
func (MsgPackSerializer) Unmarshal(data []byte, v interface{}) error { return msgpack.Unmarshal(data, v) }

var ser Serializer

func init() {
    val := os.Getenv("DEBUG_JSON")
    UseJSON, _ = strconv.ParseBool(val)
	
	if UseJSON {
        ser = JSONSerializer{}
    } else {
        ser = MsgPackSerializer{}
    }
}