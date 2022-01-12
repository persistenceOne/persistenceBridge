# persistenceBridge

This project implements an application for the Persistence Bridge Orchestrator that listens to, verifies, transforms and
relays transactions between a Cosmos-SDK chain and ethereum.

## Talk to us!

* [Twitter](https://twitter.com/PersistenceOne)
* [Telegram](https://t.me/PersistenceOneChat)
* [Discord](https://discord.com/channels/796174129077813248)

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

> Prerequisite: go1.16+ required. [ref](https://golang.org/doc/install)

> Prerequisite: git. [ref](https://github.com/git/git)

> Optional requirement: GNU make. [ref](https://www.gnu.org/software/make/manual/html_node/index.html)

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

When starting for first time `--tmStart` `--ethStart ` needs to be always given, after that not adding it will start
checking from last checked height + 1

* Starting the bridge

```shell 
persistenceBridge init
```
Start the `kafka` server. The starting script is available at `.script/startup.sh`

This generates a `config.toml` file in `$HOME/.persistenceBridge/` with empty and default values. Update this file as per configuration (Telegram configuration is not compulsory). 

Add participating validators:
```shell 
persistenceBridge add [validator_address] [validator_name]
```

Before starting the bridge, the `wrap address` (on cosmos chain) and `bridge admin` (on eth chain) should have some token balance otherwise bridge won't be able to do tx. Both the address can be queried (after all the configuration has been filled) by:
`persistenceBridge show`

Now when starting for the first time:

```shell 
persistenceBridge start --tmStart [START_HEIGHT_ON_COSMOS_CHAIN] --ethStart [START_HEIGHT_ON_ETHEREUM_CHAIN]
```

When starting next time
```shell 
persistenceBridge start
```
