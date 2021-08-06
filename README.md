# persistenceBridge

This project implements an application for the Persistence Bridge Orchestrator that listens to, verifies, transforms and relays transactions between a Cosmos-SDK chain and ethereum.

## Talk to us!
*   [Twitter](https://twitter.com/PersistenceOne)
*   [Telegram](https://t.me/PersistenceOneChat)
*   [Discord](https://discord.com/channels/796174129077813248)

## Hardware Requirements 
* **Minimal**
    * 1 GB RAM
    * 25 GB HDD
    * 1.4 GHz CPU
* **Recommended**
    * 2 GB RAM
    * 100 GB HDD
    * 2.0 GHz x2 CPU

> NOTE: SSDs have limited TBW before non-catastrophic data errors. Running a full node requires a TB+ writes per day, causing rapid deterioration of SSDs over HDDs of comparable quality.

## Operating System
* Linux/Windows/MacOS(x86)
* **Recommended**
    * Linux(x86_64)

## Installation Steps
>Prerequisite: go1.15+ required. [ref](https://golang.org/doc/install)

>Prerequisite: git. [ref](https://github.com/git/git)

>Optional requirement: GNU make. [ref](https://www.gnu.org/software/make/manual/html_node/index.html)


* Clone git repository
```shell
git clone https://github.com/persistenceOne/persistenceBridge.git
```

> Note: If running go the latest version (tested on `1.16.3`), do `export CGO_ENABLED="0"` before make install

* Make the binary  
Might require you to run `export CGO_ENABLED="0"` before make.
```shell
make all
```

* Start the bridge
```shell
./build/..../persistenceBridge start {arguments}
```

When starting for first time `--tmStart` `--ethStart ` needs to be always given,
after that not adding it will start checking from last checked height + 1

`path_to_chain_json` json file for tendermint chain, same as the relayer format -
`{"key":"acc_0","chain-id":"test","rpc-addr":"http://192.168.1.4:26657","account-prefix":"cosmos","gas-adjustment":1.5,"gas-prices":"0.025stake","trusting-period":"336h"`

* First time start
`persistenceBridge init`

`persistenceBridge add [validator_address] [validator_name]`

Update the config file generated (CASP server, vault and public keys) (Telegram configuration is not compulsory).

`persistenceBridge start chain.json --tmStart 1 --ethStart 4772131` 

>IMPORTANT NOTE: Remove `--tmStart  --ethStart` when starting next time
