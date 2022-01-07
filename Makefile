export GO111MODULE=on

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::')
COMMIT := $(shell git rev-parse --short HEAD)

build_tags = netgo

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace := $(subst ,, )
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=persistenceBridge \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=persistenceBridge \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		   -X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)

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
ifeq (${OS},Windows_NT)
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/persistenceBridge.exe ./orchestrator

else
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/persistenceBridge ./orchestrator

endif

build:
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
	go test ./... -v -timeout=20m -tags=integration -p=1
.PHONY: integration