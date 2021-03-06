## sgncli keys

Add or view local private keys

### Synopsis

Keys allows you to manage your local keystore for tendermint.

    These keys may be in any format supported by go-crypto and can be
    used by light-clients, full nodes, or any other application that
    needs to sign with a private key.

### Options

```
  -h, --help                     help for keys
      --keyring-backend string   Select keyring's backend (os|file|test) (default "os")
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

* [sgncli](sgncli.md)	 - SGN Client
* [sgncli keys ](sgncli_keys_.md)	 - 
* [sgncli keys add](sgncli_keys_add.md)	 - Add an encrypted private key (either newly generated or recovered), encrypt it, and save to disk
* [sgncli keys delete](sgncli_keys_delete.md)	 - Delete the given keys
* [sgncli keys export](sgncli_keys_export.md)	 - Export private keys
* [sgncli keys import](sgncli_keys_import.md)	 - Import private keys into the local keybase
* [sgncli keys list](sgncli_keys_list.md)	 - List all keys
* [sgncli keys migrate](sgncli_keys_migrate.md)	 - Migrate keys from the legacy (db-based) Keybase
* [sgncli keys mnemonic](sgncli_keys_mnemonic.md)	 - Compute the bip39 mnemonic for some input entropy
* [sgncli keys parse](sgncli_keys_parse.md)	 - Parse address from hex to bech32 and vice versa
* [sgncli keys show](sgncli_keys_show.md)	 - Show key info for the given name

###### Auto generated by spf13/cobra on 5-Aug-2020
