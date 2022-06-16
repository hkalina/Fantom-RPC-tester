package verifier

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/hkalina/fantom-rpc-tester/rpctypes"
	"math/big"
)

// BalancesChanges stores changes of balances on individual account during a block
type BalancesChanges struct {
	Map map[common.Address]*big.Int
	Addresses []common.Address // for deterministic Map iterating
}

func (bm *BalancesChanges) Add(address common.Address, amount *big.Int) {
	if address == (common.Address{}) || address == common.HexToAddress("0xFC00FACE00000000000000000000000000000000") {
		return // skip validation for zero/SFC address
	}
	item, exists := bm.Map[address]
	if exists {
		bm.Map[address] = item.Add(item, amount)
	} else {
		bm.Map[address] = amount
		bm.Addresses = append(bm.Addresses, address)
	}
}

// AggregateTxsIntoBalancesChanges aggregates internal txs to one balance change per account
func AggregateTxsIntoBalancesChanges(txs []rpctypes.ExternalTx) *BalancesChanges {
	movements := new(BalancesChanges)
	movements.Map = make(map[common.Address]*big.Int)
	for _, tx := range txs {
		for _, itx := range tx.InternalTxs {
			movements.Add(itx.To, itx.Value)
			movements.Add(itx.From, new(big.Int).Neg(itx.Value))
		}
	}
	return movements
}
