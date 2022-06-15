package verifier

import (
	"fmt"
	"github.com/hkalina/fantom-rpc-tester/client"
	"github.com/hkalina/fantom-rpc-tester/rpctypes"
	"log"
	"math/big"
)

type Verifier struct {
	cache  *BalanceCache
	prefix string
	debug  bool
}

func NewVerifier(id int, debug bool) *Verifier {
	v := Verifier{}
	v.cache = NewBalanceCache()
	v.prefix = fmt.Sprintf("THREAD-%d> ", id)
	v.debug = debug
	return &v
}

func (v *Verifier) VerifyRange(startBlock int64, endBlock int64, ftm *client.FtmBridge) {
	v.printf("Started on range %d-%d\n", startBlock, endBlock)
	for i := startBlock; i < endBlock; i++ {
		err := v.VerifyBlock(i, ftm)
		if err != nil {
			log.Fatalf(v.prefix+ "VerifyBlock %d failed: %s", i, err)
		}
	}
	v.printf("Finished successfully")
}

func (v *Verifier) VerifyBlock(blockNum int64, ftm *client.FtmBridge) error {
	v.printf("Verifying block %d...\n", blockNum)
	currBlockNum := big.NewInt(blockNum)
	prevBlockNum := big.NewInt(blockNum - 1)

	txs, err := ftm.GetBlockTxs(currBlockNum)
	if err != nil {
		return err
	}

	if v.debug {
		v.printTxs(txs)
	}
	movements := TxsIntoBlockMovements(txs)
	if v.debug {
		v.printMovements(movements)
	}
	for _, address := range movements.Addresses {
		amount := movements.Map[address]
		oldBalance, err := v.cache.GetBalanceFromCacheOrLoad(address, prevBlockNum, ftm.GetBalance)
		if err != nil {
			return fmt.Errorf("unable to get balance for block %s: %s", prevBlockNum.String(), err)
		}
		newBalance, err := ftm.GetBalance(address, currBlockNum)
		if err != nil {
			return fmt.Errorf("unable to get balance for block %s: %s", prevBlockNum.String(), err)
		}
		v.cache.StoreBalanceIntoCache(address, newBalance)
		computedBalance := new(big.Int).Add(oldBalance, amount)
		if computedBalance.Cmp(newBalance) != 0 {
			computedDiff := new(big.Int).Sub(computedBalance, oldBalance)
			realDiff := new(big.Int).Sub(newBalance, oldBalance)
			return fmt.Errorf(
				"unexpected balance for %s at block %d:\n" +
					" previous: %s\n computed: %s  (diff: %s)\n real:     %s  (diff: %s)",
				address, blockNum, oldBalance.String(), computedBalance.String(), computedDiff.String(), newBalance.String(), realDiff.String())
		}
	}

	return nil
}

func (v *Verifier) printf(format string, args ...interface{}) {
	log.Printf(v.prefix+ format, args...)
}

func (v *Verifier) printTxs(etxs []rpctypes.ExternalTx) {
	for _, etx := range etxs {
		v.printf("ExtTx %s: %s", etx.Hash, etx.ErrorMessage)
		for _, itx := range etx.InternalTxs {
			v.printf("    - InternalTx: %s -> %s: %s\n", itx.From, itx.To, itx.Value)
		}
	}
}

func (v *Verifier) printMovements(movements *BlockMovements) {
	v.printf("Block movements:")
	for _, address := range movements.Addresses {
		amount := movements.Map[address]
		v.printf("    - %s: %s\n", address, amount.String())
	}
}
