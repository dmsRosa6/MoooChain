package types

import "encoding/hex"

// We used 20 as size to be more compact alike to eth
// we could use whatever
type Address [20]uint8

func (a Address) ToString() string{
	return hex.EncodeToString(a[:])
}