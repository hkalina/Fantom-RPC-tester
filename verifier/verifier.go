package verifier

import (
	"fmt"
	"github.com/hkalina/fantom-rpc-tester/client"
	"github.com/hkalina/fantom-rpc-tester/rpctypes"
	"log"
	"math/big"
)

func VerifyBlock(blockNum *big.Int, ftm *client.FtmBridge, debug bool) error {
	log.Printf("Verifying block %s...\n", blockNum.String())
	prevBlockNum := new(big.Int).Sub(blockNum, big.NewInt(1))

	txs, err := ftm.GetBlockTxs(blockNum)
	if err != nil {
		return err
	}

	if debug {
		printTxs(txs)
	}
	movements := TxsIntoBlockMovements(txs)
	if debug {
		printMovements(movements)
	}
	for _, address := range movements.Addresses {
		amount := movements.Map[address]
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
		}
	}

	return nil
}

func printTxs(txs []rpctypes.ExternalTx) {
	for _, tx := range txs {
		log.Printf("ExtTx %s:", tx.Hash)
		for _, itx := range tx.InternalTxs {
			log.Printf("    - InternalTx: %s -> %s: %s\n", itx.From, itx.To, itx.Value)
		}
	}
}

func printMovements(movements *BlockMovements) {
	log.Printf("Block movements:")
	for _, address := range movements.Addresses {
		amount := movements.Map[address]
		log.Printf("    - %s: %s\n", address, amount.String())
	}
}
