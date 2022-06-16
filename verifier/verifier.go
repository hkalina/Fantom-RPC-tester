package verifier

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hkalina/fantom-rpc-tester/client"
	"github.com/hkalina/fantom-rpc-tester/rpctypes"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"time"
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
	for i := startBlock; i <= endBlock; i++ {
		err := v.VerifyBlockRepeatedly(i, ftm)
		if err != nil {
			log.Fatalf(v.prefix+"VerifyBlock %d failed: %s", i, err)
		}
	}
	v.printf("Finished successfully")
}

func (v *Verifier) VerifyBlockRepeatedly(blockNum int64, ftm *client.FtmBridge) (err error) {
	for attempt := 1; attempt <= 5; attempt++ {
		err = v.VerifyBlock(blockNum, ftm)
		if err == nil || strings.HasPrefix(err.Error(), "unexpected balance") {
			return nil
		}
		v.printf("VerifyBlock(%d) failed (attempt %d): %s\n", blockNum, attempt, err)
		time.Sleep((time.Duration(rand.Intn(1000)) + 100) * time.Millisecond)
	}
	return err
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
	balancesChanges := AggregateTxsIntoBalancesChanges(txs)
	if v.debug {
		v.printBalancesChanges(balancesChanges)
	}
	newBalances := make(map[common.Address]*big.Int, 50)
	for _, address := range balancesChanges.Addresses {
		balanceChange := balancesChanges.Map[address]

		// get the account balance at the beginning of the block
		oldBalance, err := v.cache.GetBalanceFromCacheOrLoad(address, prevBlockNum, ftm.GetBalance)
		if err != nil {
			return fmt.Errorf("unable to get balance for block %s: %s", prevBlockNum.String(), err)
		}

		// get the account balance at the end of the block
		newBalance, err := ftm.GetBalance(address, currBlockNum)
		if err != nil {
			return fmt.Errorf("unable to get balance for block %s: %s", prevBlockNum.String(), err)
		}
		newBalances[address] = newBalance

		// compare balance computed from the block internal txs and the old balance with the new balance
		computedBalance := new(big.Int).Add(oldBalance, balanceChange)
		if computedBalance.Cmp(newBalance) != 0 {
			computedDiff := new(big.Int).Sub(computedBalance, oldBalance)
			realDiff := new(big.Int).Sub(newBalance, oldBalance)
			return fmt.Errorf(
				"unexpected balance for %s at block %d:\n"+
					" previous: %s\n computed: %s  (diff: %s)\n real:     %s  (diff: %s)",
				address, blockNum, oldBalance.String(), computedBalance.String(), computedDiff.String(), newBalance.String(), realDiff.String())
		}
	}

	// override cached balance only when the block succeed (allow retry)
	for address, newBalance := range newBalances {
		v.cache.StoreBalanceIntoCache(address, newBalance)
	}

	return nil
}

func (v *Verifier) printf(format string, args ...interface{}) {
	log.Printf(v.prefix+format, args...)
}

func (v *Verifier) printTxs(etxs []rpctypes.ExternalTx) {
	for _, etx := range etxs {
		v.printf("ExtTx %s: %s", etx.Hash, etx.ErrorMessage)
		for _, itx := range etx.InternalTxs {
			v.printf("    - InternalTx: %s -> %s: %s\n", itx.From, itx.To, itx.Value)
		}
	}
}

func (v *Verifier) printBalancesChanges(changes *BalancesChanges) {
	v.printf("Block balances changes:")
	for _, address := range changes.Addresses {
		amount := changes.Map[address]
		v.printf("    - %s: %s\n", address, amount.String())
	}
}
