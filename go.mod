module github.com/persistenceOne/persistenceBridge

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Shopify/sarama v1.28.0
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/cosmos/cosmos-sdk v0.42.6
	github.com/cosmos/relayer v0.9.3
	github.com/dgraph-io/badger/v3 v3.2011.1
	github.com/ethereum/go-ethereum v1.10.8
	github.com/gorilla/mux v1.8.0
	github.com/gravity-devs/liquidity v1.2.9
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.11
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/grpc v1.37.0
	gopkg.in/tucnak/telebot.v2 v2.3.5
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2
