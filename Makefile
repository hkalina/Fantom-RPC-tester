rpctester:
	go build -o build/rpctester ./cmd/rpctester

test: rpctester
	build/rpctester http://rpcapi-tracing.fantom.network/ 40428976 40428997 2
