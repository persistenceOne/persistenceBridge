export GO111MODULE=on

export SMART_CONTRACT_VERSION=fbda08567aac3acd004e36dfe670a44253453907
export LIQUID_STAKING=LiquidStakingV3
export TOKEN_WRAPPER=TokenWrapperV3
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

build: .script/compileSC.sh build-only

build-only:
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

lintci:
	golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 --config .golangci.yaml
	#cosmossec -quiet -tests -nosec ./... # fails with panic
.PHONY: lintci

lintci-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOBIN} v1.43.0
	# git clone --depth 1 https://github.com/informalsystems/gosec.git cosmossec && cd cosmossec/cmd/gosec && go build -o ${GOBIN}/cosmossec && chmod 0544 ${GOBIN}/cosmossec && cd ../../.. && rm -Rdf ./cosmossec
.PHONY: lintci-install

lintci-remove:
	rm ${GOBIN}/golangci-lint
	# rm -Rdf ./gosec ./cosmossec ${GOBIN}/cosmossec
.PHONY: lintci-remove

lintci-update: lintci-remove lintci-install
.PHONY: lintci-update

goimports:
	goimports -local="github.com/persistenceOne/persistenceBridge" -w .
.PHONY: goimports

generate:
	go generate ./application/configuration/...
.PHONY: generate

deps:
	go install github.com/globusdigital/deep-copy@latest
.PHONY: deps

tests: units integration
.PHONY: tests

units:
	go test ./... -v -timeout=5m -tags=units
.PHONY: units

integration:
	go test ./... -v -timeout=20m -tags=integration -p=1 -parallel=1
.PHONY: integration