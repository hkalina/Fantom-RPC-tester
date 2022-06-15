package client

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/hkalina/fantom-rpc-tester/rpctypes"
	"log"
	"math/big"
)

type FtmBridge struct {
	rpc *rpc.Client
	eth *ethclient.Client
}

func NewFtmBridge(rpcUrl string) *FtmBridge {
	rpcClient, err := rpc.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}

	ethClient, err := ethclient.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}

	return &FtmBridge{
		rpc: rpcClient,
		eth: ethClient,
	}
}

func (ftm *FtmBridge) Close() {
	if ftm.rpc != nil {
		ftm.rpc.Close()
		ftm.eth.Close()
	}
}

func (ftm *FtmBridge) GetBalance(address common.Address, block *big.Int) (*big.Int, error) {
	return ftm.eth.BalanceAt(context.Background(), address, block)
}

func allTracesSucceed(trace []TxTrace) bool {
	for i, trc := range trace {
		if trc.Error != "" {
			log.Printf("trace %d error: %s", i, trc.Error)
			return false
		}
	}
	return true
}

func (ftm *FtmBridge) GetBlockTxs(blockNum *big.Int, debug bool) (etxs []rpctypes.ExternalTx, err error) {
	block, err := ftm.GetBlock(blockNum)
	if err != nil {
		return nil, fmt.Errorf("GetBlock failed: %s", err)
	}

	var trace []TxTrace
	var ok bool
	for i := 0; i < 3; i++ {
		trace, err = ftm.TraceBlockByNumber(blockNum)
		if err != nil {
			return nil, fmt.Errorf("TraceBlockByNumber failed: %s", err)
		}
		if allTracesSucceed(trace) {
			ok = true
			break
		}
	}
	if !ok {
		return nil, fmt.Errorf("all TraceBlockByNumber attempts failed")
	}

	txsHashes := make([]common.Hash, 0)
	for i, tx := range block.Txs {
		etx := rpctypes.ExternalTx{
			Hash:     tx.Hash,
			GasPrice: big.Int(tx.GasPrice),
			From:     tx.From,
			To:       tx.To,
		}
		txsHashes = append(txsHashes, etx.Hash)
		etx.InternalTxs = trace[i].Result.InternalTxs() // extract internal txs from trace
		//etx.GasUsed = (*big.Int)(trace[i].Result.GasUsed)
		etx.Revert = trace[i].Result.Revert
		etx.ErrorMessage = trace[i].Result.ErrorMessage
		if debug && etx.ErrorMessage != "" {
			log.Printf("trace of tx %s revert: %s (revert: %t)", tx.Hash, etx.ErrorMessage, etx.Revert)
		}
		etxs = append(etxs, etx)
	}

	receipts, err := ftm.GetReceipts(txsHashes)
	for i, receipt := range receipts {
		if receipt.TxHash != etxs[i].Hash {
			return nil, fmt.Errorf("receipt for %s returned when %s requested", receipt.TxHash, etxs[i].Hash)
		}
		if receipt.BlockNumber.Cmp(blockNum) != 0 {
			return nil, fmt.Errorf("block number differes for %s - expected %s, got %s", receipt.TxHash, blockNum.String(), receipt.BlockNumber.String())
		}
		etxs[i].GasUsed.SetUint64(receipt.GasUsed)

		// add fee internal tx
		feeAmount := new(big.Int).Mul(&etxs[i].GasPrice, &etxs[i].GasUsed)
		etxs[i].InternalTxs = append(etxs[i].InternalTxs, rpctypes.InternalTx{
			From:    etxs[i].From,
			To:      common.Address{}, // use zero-address as destination for fees (for now)
			Value:   feeAmount,
			GasUsed: nil,
		})
	}

	return etxs, nil
}

func (ftm *FtmBridge) GetReceipts(
	txs []common.Hash,
) ([]*types.Receipt, error) {
	receipts := make([]*types.Receipt, len(txs))
	if len(txs) == 0 {
		return receipts, nil
	}

	reqs := make([]rpc.BatchElem, len(txs))
	for i := range reqs {
		reqs[i] = rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{txs[i].Hex()},
			Result: &receipts[i],
		}
	}
	if err := ftm.rpc.BatchCallContext(context.Background(), reqs); err != nil {
		return nil, err
	}
	for i := range reqs {
		if reqs[i].Error != nil {
			return nil, reqs[i].Error
		}
		if receipts[i] == nil {
			return nil, fmt.Errorf("got empty receipt for %x", txs[i].Hex())
		}
	}
	return receipts, nil
}
