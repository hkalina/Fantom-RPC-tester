package verifier

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hkalina/fantom-rpc-tester/client"
	"github.com/hkalina/fantom-rpc-tester/rpctypes"
	"log"
	"math/big"
)

func VerifyBlock(blockNum *big.Int, ftm *client.FtmBridge) error {
	log.Printf("Verifying block %s...\n", blockNum.String())
	prevBlockNum := new(big.Int).Sub(blockNum, big.NewInt(1))

	txs, err := ftm.GetBlockTxs(blockNum)
	if err != nil {
		return err
	}

	movements := MergeTxsIntoAggregatedBlockMovements(txs)
	for _, address := range movements.Addresses {
		amount := movements.Map[address]
		fmt.Printf("Move: %s: %s\n", address, amount.String())
		oldBalance, err := ftm.GetBalance(address, prevBlockNum)
		if err != nil {
			return fmt.Errorf("unable to get balance for block %s: %s", prevBlockNum.String(), err)
		}
		newBalance, err := ftm.GetBalance(address, blockNum)
		if err != nil {
			return fmt.Errorf("unable to get balance for block %s: %s", prevBlockNum.String(), err)
		}
		computedBalance := new(big.Int).Add(oldBalance, amount)
		if computedBalance.Cmp(newBalance) != 0 {
			computedDiff := new(big.Int).Sub(computedBalance, oldBalance)
			realDiff := new(big.Int).Sub(newBalance, oldBalance)
			return fmt.Errorf(
				"unexpected balance for %s at block %s:\n" +
					" previous: %s\n computed: %s  (diff: %s)\n real:     %s  (diff: %s)",
				address, blockNum.String(), oldBalance.String(), computedBalance.String(), computedDiff.String(), newBalance.String(), realDiff.String())
		} else {
			fmt.Printf("New balance OK\n")
		}
	}

	return nil
}

type BlockMovements struct {
	Map map[common.Address]*big.Int
	Addresses []common.Address // for deterministic Map iterating
}

func (bm *BlockMovements) Add(address common.Address, amount *big.Int) {
	if address == (common.Address{}) {
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

func MergeTxsIntoAggregatedBlockMovements(txs []rpctypes.ExternalTx) *BlockMovements {
	movements := new(BlockMovements)
	movements.Map = make(map[common.Address]*big.Int)
	for _, tx := range txs {
		for _, itx := range tx.InternalTxs {
			fmt.Printf("InternalTx: %s -> %s: %s\n", itx.From, itx.To, itx.Value)
			movements.Add(itx.To, itx.Value)
			movements.Add(itx.From, new(big.Int).Neg(itx.Value))
		}
	}
	return movements
}
