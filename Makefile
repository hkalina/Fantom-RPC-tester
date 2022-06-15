rpctester:
	go build -o build/rpctester ./cmd/rpctester

test:
	build/rpctester http://rpcapi-tracing.fantom.network/ 40428976 40428977 1 debug
