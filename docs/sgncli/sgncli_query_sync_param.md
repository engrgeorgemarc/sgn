## sgncli query sync param

Query the parameters (voting|tallying) of the sync process

### Synopsis

Query the all the parameters for the sync process.

Example:
$ <appcli> query sync param voting
$ <appcli> query sync param tallying

```
sgncli query sync param [param-type] [flags]
```

### Options

```
      --height int    Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help          help for param
      --indent        Add indent to JSON response
      --ledger        Use a connected Ledger device
      --node string   <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --trust-node    Trust connected full node (don't verify proofs for responses)
```

### Options inherited from parent commands

```
      --chain-id string   Chain ID of tendermint node
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.sgncli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```

### SEE ALSO

* [sgncli query sync](sgncli_query_sync.md)	 - Querying commands for the sync module

###### Auto generated by spf13/cobra on 5-Aug-2020
