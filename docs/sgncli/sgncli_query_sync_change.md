## sgncli query sync change

Query details of a single change

### Synopsis

Query details for a change. You can find the
change-id by running "<appcli> query sync changes".

Example:
$ <appcli> query sync change 1

```
sgncli query sync change [change-id] [flags]
```

### Options

```
      --height int    Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help          help for change
      --indent        Add indent to JSON response
      --ledger        Use a connected Ledger device
      --node string   <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --trust-node    Trust connected full node (don't verify proofs for responses)
```

### Options inherited from parent commands

```
      --chain-id string   Chain ID of tendermint node
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/Users/Frank/.sgncli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```

### SEE ALSO

* [sgncli query sync](sgncli_query_sync.md)	 - Querying commands for the sync module

###### Auto generated by spf13/cobra on 22-Jul-2020