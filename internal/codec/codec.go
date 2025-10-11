package codec

import "encoding/hex"

const (
	BlockKeyword     = "Block:"
	PrevBlockKeyword = "Block:prev:"
)

func HexEncode(b []byte) string {
	return hex.EncodeToString(b)
}

func HexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func BuildBlockKey(hash []byte) string {
	return BlockKeyword + HexEncode(hash)
}

func BuildPrevBlockKey(hash []byte) string {
	return PrevBlockKeyword + HexEncode(hash)
}
