rpctester:
	go build -o build/rpctester ./cmd/rpctester

test:
	build/rpctester https://rpcapi-tracing.fantom.network/ 40428976 40428977
