package main

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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
