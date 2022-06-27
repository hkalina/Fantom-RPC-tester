Go-opera Integration Tests
==========================

The consistency testing tool, which compares reported internal transactions with reported balances at given blocks.
The internal transactions are obtained by replaying block transactions using `debug_traceBlockByNumber`
RPC tracing interface. The balances are obtained using `eth_getBalance` RPC call.

The tool to be called as:
```
rpctester http://rpcapi-tracing.fantom.network/ {blockFromNumber} {blockToNumber} {parallelThreads} [debug]
```

The tool does not use any persistent storage, but accounts balances are cached in memory in thread-local LRU cache.
