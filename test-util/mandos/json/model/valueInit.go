package mandosjsonmodel

import (
	"math/big"
)

// JSONBytes stores the parsed byte slice value but also the original parsed string
type JSONBytes struct {
	Value    []byte
	Original string
}

// JSONBigInt stores the parsed big int value but also the original parsed string
type JSONBigInt struct {
	Value    *big.Int
	Original string
}

// JSONUint64 stores the parsed uint64 value but also the original parsed string
type JSONUint64 struct {
	Value    uint64
	Original string
}
