package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
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

func (ftm *FtmBridge) GetCode(contract common.Address, block hexutil.Big) string {
	var code string
	if err := ftm.rpc.Call(&code, "eth_getCode", contract, block.String()); err != nil {
		log.Printf("failed eth_getCode: %s", err)
		return ""
	}
	return code
}

func (ftm *FtmBridge) GetBalance(address common.Address, block hexutil.Big) string {
	var code string
	if err := ftm.rpc.Call(&code, "eth_getBalance", address, block.String()); err != nil {
		log.Printf("failed eth_getBalance: %s", err)
		return ""
	}
	return code
}
