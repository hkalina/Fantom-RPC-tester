package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"log"
	"math/big"
)

// TxTrace represents output of debug_trace for one transaction.
type TxTrace struct {
	Result Call   `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Call represents one call operation in a TxTrace.
type Call struct {
	Type         string         `json:"type"`
	From         common.Address `json:"from"`
	To           common.Address `json:"to"`
	Value        *hexutil.Big   `json:"value"`
	GasUsed      *hexutil.Big   `json:"gasUsed"`
	Revert       bool           `json:"revert,omitempty"`
	ErrorMessage string         `json:"error,omitempty"`
	Calls        []Call         `json:"calls,omitempty"`
}

// InternalTx represents one internal transaction derived from a Call.
type InternalTx struct {
	From    common.Address `json:"from"`
	To      common.Address `json:"to"`
	Value   *hexutil.Big   `json:"value"`
	GasUsed *hexutil.Big   `json:"gasUsed"`
}

func (data *Call) InternalTxs() (txs []InternalTx) {
	if data.Revert != false || data.ErrorMessage != "" {
		return
	}
	if data.Type == "CALL" && data.Value != nil && data.Value.ToInt().Sign() != 0 {
		txs = append(txs, InternalTx{
			From:    data.From,
			To:      data.To,
			Value:   data.Value,
			GasUsed: data.GasUsed,
		})
	}
	for _, child := range data.Calls {
		txs = append(txs, child.InternalTxs()...)
	}
	return
}

func (ftm *FtmBridge) TraceBlockByNumber(block *big.Int) ([]TxTrace, error) {
	var result []TxTrace
	tracer := "callTracer"
	options := tracers.TraceConfig{
		Tracer: &tracer,
	}

	if err := ftm.rpc.Call(&result, "debug_traceBlockByNumber", (*hexutil.Big)(block).String(), options); err != nil {
		log.Printf("failed debug_traceBlockByNumber: %s", err)
		return nil, err
	}
	return result, nil
}
