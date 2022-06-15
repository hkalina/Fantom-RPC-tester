package verifier

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/tinylru"
	"math/big"
)

type BalanceCache struct {
	cache tinylru.LRU
}

func NewBalanceCache() *BalanceCache {
	c := BalanceCache{}
	c.cache.Resize(100_000_000)
	return &c
}

func (c *BalanceCache) GetBalanceFromCacheOrLoad(address common.Address, block *big.Int, loader func(common.Address, *big.Int) (*big.Int, error)) (*big.Int, error) {
	val, ok := c.cache.Get(address)
	if ok {
		return val.(*big.Int), nil
	}
	balance, err := loader(address, block)
	c.cache.Set(address, balance)
	return balance, err
}

func (c *BalanceCache) StoreBalanceIntoCache(address common.Address, balance *big.Int) {
	c.cache.Set(address, balance)
}
