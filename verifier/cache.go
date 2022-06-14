package verifier

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/tinylru"
	"math/big"
)

var cache tinylru.LRU

func InitCache() {
	cache.Resize(1_000_000)
}

func GetBalanceFromCacheOrLoad(address common.Address, block *big.Int, loader func(common.Address, *big.Int) (*big.Int, error)) (*big.Int, error) {
	val, ok := cache.Get(address)
	if ok {
		return val.(*big.Int), nil
	}
	balance, err := loader(address, block)
	cache.Set(address, balance)
	return balance, err
}

func StoreBalanceIntoCache(address common.Address, balance *big.Int) {
	cache.Set(address, balance)
}
