package main

import (
	"fmt"
	"github.com/hkalina/fantom-rpc-tester/client"
	"github.com/hkalina/fantom-rpc-tester/verifier"
	"log"
	"os"
	"strconv"
	"sync"
)

var ftm *client.FtmBridge
var debug bool
var wg sync.WaitGroup

func main() {
	if len(os.Args) != 5 && len(os.Args) != 6 {
		fmt.Printf("Usage: %s http://rpc1/ [blockNumberStart] [blockNumberEnd] [threads] [debug]\n", os.Args[0])
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

	threads, err := strconv.ParseInt(os.Args[4], 10, 32)
	if err != nil {
		log.Fatal("Invalid the threads count")
	}

	if len(os.Args) == 6 && os.Args[5] == "debug" {
		debug = true
	}

	blocksPerThread := (endBlock - startBlock) / threads + 1

	start := startBlock
	wg.Add(int(threads))
	for thread := 0; thread < int(threads); thread++ {
		verif := verifier.NewVerifier(thread, debug)
		end := start + blocksPerThread
		if end > endBlock {
			end = endBlock
		}

		go func(startBlock, endBlock int64) {
			defer wg.Done()

			ftm = client.NewFtmBridge(os.Args[1])
			defer ftm.Close()

			verif.VerifyRange(startBlock, endBlock, ftm)

		}(start, end)
		start += blocksPerThread + 1
	}
	wg.Wait()
}
