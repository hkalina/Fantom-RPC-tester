package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/hkalina/fantom-rpc-tester/rpctypes"
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
	ErrorMessage string `json:"error,omitempty"`
	Calls        []Call `json:"calls,omitempty"`
}

func (data *Call) InternalTxs() (txs []rpctypes.InternalTx) {
	if data.Revert != false || data.ErrorMessage != "" {
		return
	}
	// TODO: check data.Type?
	if data.Value != nil && data.Value.ToInt().Sign() != 0 {
		txs = append(txs, rpctypes.InternalTx{
			From:    data.From,
			To:      data.To,
			Value:   (*big.Int)(data.Value),
			GasUsed: (*big.Int)(data.GasUsed),
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
