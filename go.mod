module github.com/persistenceOne/persistenceBridge

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/Shopify/sarama v1.28.0
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/cosmos/cosmos-sdk v0.44.3
	github.com/cosmos/ibc-go/v2 v2.0.0
	github.com/cosmos/relayer v1.0.0
	github.com/dgraph-io/badger/v3 v3.2011.1
	github.com/ethereum/go-ethereum v1.10.8
	github.com/gorilla/mux v1.8.0
	github.com/gravity-devs/liquidity v1.4.2
	github.com/spf13/cobra v1.2.1
	github.com/strangelove-ventures/packet-forward-middleware v1.0.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/grpc v1.40.0
	gopkg.in/tucnak/telebot.v2 v2.3.5
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2
