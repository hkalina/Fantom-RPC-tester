package main

import (
	"fmt"
	"github.com/hkalina/fantom-rpc-tester/client"
	"github.com/hkalina/fantom-rpc-tester/verifier"
	"log"
	"math/big"
	"os"
	"strconv"
)

var ftm *client.FtmBridge

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s http://rpc1/ [blockNumberStart] [blockNumberEnd]\n", os.Args[0])
		os.Exit(1)
	}

	startBlock, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Fatal("Invalid start block argument")
	}
	endBlock, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		log.Fatal("Invalid end block argument")
	}
	if ! (startBlock < endBlock) {
		fmt.Printf("The start block argument has to be less than end block argument\n")
		os.Exit(1)
	}

	ftm = client.NewFtmBridge(os.Args[1])
	defer ftm.Close()

	for i := startBlock; i < endBlock; i++ {
		err := verifier.VerifyBlock(big.NewInt(i), ftm)
		if err != nil {
			log.Fatalf("VerifyBlock %d failed: %s", i, err)
		}
	}
}
