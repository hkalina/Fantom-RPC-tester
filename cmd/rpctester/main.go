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

	itxs, err := ftm.GetBlockInternalTxs(&startBlock)
	if err != nil {
		log.Fatal("GetBlockInternalTxs failed: ", err)
	}

	bytes, err := json.Marshal(itxs)
	if err != nil {
		log.Fatal("JSON marshal failed: ", err)
	}
	fmt.Printf("%s\n", string(bytes))
}
