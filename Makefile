export GO111MODULE=on

export SMART_CONTRACT_VERSION=4aac303e15d4ad01ffdf44f1ce376f0f509c765a
export LIQUID_STAKING=LiquidStakingV2
export TOKEN_WRAPPER=TokenWrapperV2
export OPENZEPPELIN_VERSION=v3.4.2

VERSION := $(shell echo $(shell git describe --always) | sed 's/^v//')

build_tags = netgo

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace := $(subst ,, )
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

ldflags = -X github.com/persistenceOne/persistenceBridge/application/commands.Version=$(VERSION)

ifeq (cleveldb,$(findstring cleveldb,$(build_tags)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (badgerdb,$(findstring badgerdb,$(build_tags)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb
endif
ifeq (rocksdb,$(findstring rocksdb,$(build_tags)))
  CGO_ENABLED=1
endif

BUILD_FLAGS += -ldflags "${ldflags}" -tags "${build_tags}"

GOBIN = $(shell go env GOPATH)/bin
GOARCH = $(shell go env GOARCH)
GOOS = $(shell go env GOOS)

.PHONY: all install build verify

all: verify install

install:
	.script/compileSC.sh
ifeq (${OS},Windows_NT)
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/persistenceBridge.exe ./orchestrator

else
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/persistenceBridge ./orchestrator

endif

build:
	.script/compileSC.sh
ifeq (${OS},Windows_NT)
	go build  ${BUILD_FLAGS} -o build/${GOOS}/${GOARCH}/persistenceBridge.exe ./orchestrator

else
	go build  ${BUILD_FLAGS} -o build/${GOOS}/${GOARCH}/persistenceBridge ./orchestrator

endif

verify:
	@echo "verifying modules"
	@go mod verify

release: build
	mkdir -p release
ifeq (${OS},Windows_NT)
	tar -czvf release/persistenceBridge-${GOOS}-${GOARCH}.tar.gz --directory=build/${GOOS}/${GOARCH} persistenceBridge.exe
else
	tar -czvf release/persistenceBridge-${GOOS}-${GOARCH}.tar.gz --directory=build/${GOOS}/${GOARCH} persistenceBridge
endif


test:
	@go test ./application/casp
	@go test ./application/commands
	@go test ./application/configuration
	@go test ./application/db
	@go test ./application/outgoingTx
	@go test ./application/rest
	@go test ./application/rpc
	@go test ./application/shutdown

	go test ./ethereum/contracts
	go test ./ethereum

	go test ./tendermint