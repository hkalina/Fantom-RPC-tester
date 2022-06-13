package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
)

var ftm *FtmBridge
var startBlock big.Int
var endBlock big.Int

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s http://rpc1/ [blockNumberStart] [blockNumberEnd]\n", os.Args[0])
		os.Exit(1)
	}

	startBlock.SetString(os.Args[2], 10)
	endBlock.SetString(os.Args[3], 10)

	if startBlock.Cmp(&endBlock) != -1 {
		fmt.Printf("The test [blockNumberStart] has to be less than [blockNumberEnd]\n")
		os.Exit(1)
	}

	ftm = NewFtmBridge(os.Args[1])
	defer ftm.Close()

	/*
		bal, err := ftm.GetBalance(common.HexToAddress("0x83A6524Be9213B1Ce36bCc0DCEfb5eb51D87aD10"), &startBlock)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("balance: %s\n", bal)
	*/

	trace, err := ftm.TraceBlockByNumber(&startBlock)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range trace {
		bytes, err := json.Marshal(tx.Result.InternalTxs())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(bytes))
	}
}
