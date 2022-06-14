package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"log"
	"math/big"
)

type BlockResult struct {
	Number hexutil.Big `json:"number"`
	Hash common.Hash `json:"hash"`
	Txs []BlockTransaction `json:"transactions"`
	BaseFeePerGas hexutil.Big `json:"baseFeePerGas"`
}

type BlockTransaction struct {
	Hash common.Hash `json:"hash"`
	From common.Address `json:"from"`
	To common.Address `json:"to"`
	Value hexutil.Big `json:"value"`
	GasLimit hexutil.Big `json:"gas"`
	GasPrice hexutil.Big `json:"gasPrice"`
}

func (ftm *FtmBridge) GetBlock(block *big.Int) (*BlockResult, error) {
	var result BlockResult
	if err := ftm.rpc.Call(&result, "eth_getBlockByNumber", (*hexutil.Big)(block).String(), true); err != nil {
		log.Printf("failed eth_getBlockByNumber: %s", err)
		return nil, err
	}
	return &result, nil
}
