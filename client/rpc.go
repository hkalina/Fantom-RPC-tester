package client

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/hkalina/fantom-rpc-tester/rpctypes"
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

func (ftm *FtmBridge) GetBlock(block *big.Int) (*types.Block, error) {
	return ftm.eth.BlockByNumber(context.Background(), block)
}

func (ftm *FtmBridge) GetBlockTxs(blockNum *big.Int) (out []rpctypes.ExternalTx, err error) {
	block, err := ftm.GetBlock(blockNum)
	if err != nil {
		return nil, fmt.Errorf("GetBlock failed: %s", err)
	}
	trace, err := ftm.TraceBlockByNumber(blockNum)
	if err != nil {
		return nil, fmt.Errorf("TraceBlockByNumber failed: %s", err)
	}

	for i, tx := range block.Transactions() {
		etx := rpctypes.ExternalTx{
			Hash: tx.Hash(),
		}
		if trace[i].Error != "" {
			return nil, fmt.Errorf("trace of tx %s error: %s", tx.Hash(), trace[i].Error)
		}

		etx.InternalTxs = trace[i].Result.InternalTxs() // extract internal rpctypes from trace
		etx.GasUsed = (*big.Int)(trace[i].Result.GasUsed)
		etx.Revert = trace[i].Result.Revert
		etx.ErrorMessage = trace[i].Result.ErrorMessage

		feeAmount := new(big.Int).Mul(tx.GasPrice(), etx.GasUsed)

		// derive tx sender
		etx.From, err = types.NewLondonSigner(tx.ChainId()).Sender(tx)
		if err != nil {
			return nil, fmt.Errorf("NewLondonSigner.Sender failed: %s", err)
		}

		// add fee internal tx
		etx.InternalTxs = append(etx.InternalTxs, rpctypes.InternalTx{
			From:    etx.From,
			To:      common.Address{},
			Value:   feeAmount,
			GasUsed: nil,
		})
		out = append(out, etx)
	}

	return out, nil
}
