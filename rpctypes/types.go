package rpctypes

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// InternalTx represents one internal transaction.
type InternalTx struct {
	From    common.Address `json:"from"`
	To      common.Address `json:"to"`
	Value   *big.Int       `json:"value"`
	GasUsed *big.Int       `json:"gasUsed"`
}

type ExternalTx struct {
	Hash         common.Hash
	From         common.Address
	To           *common.Address
	InternalTxs  []InternalTx
	GasUsed      big.Int
	GasPrice     big.Int
	Revert       bool
	ErrorMessage string
}
