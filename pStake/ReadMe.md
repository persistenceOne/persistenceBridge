If running go the latest version (tested on `1.16.3`), do `export CGO_ENABLED="0"` before make install

When starting for first time `--tmStart` `--ethStart ` needs to be always given.
After that not adding it will start checking from last checked height + 1

`path_to_chain_json` : json file for tendermint chain, same as relayer format -
`{"key":"acc_0","chain-id":"test","rpc-addr":"http://192.168.1.4:26657","account-prefix":"cosmos","gas-adjustment":1.5,"gas-prices":"0.025stake","trusting-period":"336h"`

`--ethPrivateKey` private key of account which will do txs to eth

First time start:
`persistenceCore pStake chain.json "wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please" --tmStart 1 --ethStart 4772131 --ethPrivateKey [ETH_ACC_PRIVATE_KEY]`

Remove `--tmStart  --ethStart` when starting next time