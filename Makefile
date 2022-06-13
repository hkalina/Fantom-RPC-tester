rpctester:
	go build -o build/rpctester ./cmd/rpctester

test:
	build/rpctester http://rpcapi.fantom.network/ 30000000 30000100
