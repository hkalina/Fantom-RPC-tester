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

	movements := MergeBlockMovements(txs)
	for address, amount := range movements {
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
			return fmt.Errorf("unexpected balance for %s at block %s:\n previous: %s\n computed: %s\n real:     %s",
				address, blockNum.String(), oldBalance.String(), computedBalance.String(), newBalance.String())
		} else {
			fmt.Printf("New balance OK\n")
		}
	}

	return nil
}

type BlockMovements map[common.Address]*big.Int

func (bm BlockMovements) Add(address common.Address, amount *big.Int) {
	if address == (common.Address{}) {
		return // skip validation for zero address
	}
	item, exists := bm[address]
	if exists {
		bm[address] = item.Add(item, amount)
	} else {
		bm[address] = amount
	}
}

func MergeBlockMovements(txs []rpctypes.ExternalTx) BlockMovements {
	movements := make(BlockMovements)
	for _, tx := range txs {
		for _, itx := range tx.InternalTxs {
			movements.Add(itx.To, itx.Value)
		}
	}
	return movements
}
