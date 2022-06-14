package verifier

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/hkalina/fantom-rpc-tester/rpctypes"
	"math/big"
)

type BlockMovements struct {
	Map map[common.Address]*big.Int
	Addresses []common.Address // for deterministic Map iterating
}

func (bm *BlockMovements) Add(address common.Address, amount *big.Int) {
	if address == (common.Address{}) || address == common.HexToAddress("0xFC00FACE00000000000000000000000000000000") {
		return // skip validation for zero address
	}
	item, exists := bm.Map[address]
	if exists {
		bm.Map[address] = item.Add(item, amount)
	} else {
		bm.Map[address] = amount
		bm.Addresses = append(bm.Addresses, address)
	}
}

func TxsIntoBlockMovements(txs []rpctypes.ExternalTx) *BlockMovements {
	movements := new(BlockMovements)
	movements.Map = make(map[common.Address]*big.Int)
	for _, tx := range txs {
		for _, itx := range tx.InternalTxs {
			movements.Add(itx.To, itx.Value)
			movements.Add(itx.From, new(big.Int).Neg(itx.Value))
		}
	}
	return movements
}
