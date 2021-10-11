#!/bin/bash

OPENZEPPELIN_VERSION="v3.4.2"
OPENZEPPELIN_DIR_NAME="@openzeppelin"

rm -rf pStake-smartContracts
rm -rf "$OPENZEPPELIN_DIR_NAME"

echo "Downloading pStake-smartContracts..."
git clone git@github.com:persistenceOne/pStake-smartContracts.git
# shellcheck disable=SC2164
cd pStake-smartContracts
git checkout "$SC_COMMIT_HASH"
# shellcheck disable=SC2103
cd ..

echo "Downloading openzeppelin..."
mkdir "$OPENZEPPELIN_DIR_NAME"
# shellcheck disable=SC2164
cd "$OPENZEPPELIN_DIR_NAME"

echo "Downloading openzeppelin-contracts..."
git clone https://github.com/OpenZeppelin/openzeppelin-contracts.git
# shellcheck disable=SC2164
cd openzeppelin-contracts
git checkout "$OPENZEPPELIN_VERSION"
cd ..
mv openzeppelin-contracts/contracts ./
rm -rf openzeppelin-contracts

echo "Downloading openzeppelin-contract-upgradeable..."
git clone https://github.com/OpenZeppelin/openzeppelin-contracts-upgradeable.git
# shellcheck disable=SC2164
cd openzeppelin-contracts-upgradeable
git checkout "$OPENZEPPELIN_VERSION"
cd ..
mv openzeppelin-contracts-upgradeable/contracts ./contracts-upgradeable
rm -rf openzeppelin-contracts-upgradeable

cd ..
mv "$OPENZEPPELIN_DIR_NAME" ./pStake-smartContracts/
# shellcheck disable=SC2164
cd pStake-smartContracts

echo "Compiling smart contracts..."
# shellcheck disable=SC2046
solc --abi --bin contracts/LiquidStakingV2.sol -o build @openzeppelin/=$(pwd)/@openzeppelin/
abigen --bin ./build/LiquidStakingV2.bin --abi ./build/LiquidStakingV2.abi --pkg=liquidStaking --out=liquidStaking.go

# shellcheck disable=SC2046
solc --abi --bin contracts/TokenWrapperV2.sol -o build --overwrite @openzeppelin/=$(pwd)/@openzeppelin/
abigen --bin ./build/TokenWrapperV2.bin --abi ./build/TokenWrapperV2.abi --pkg=tokenWrapper --out=tokenWrapper.go

cd ..
mv pStake-smartContracts/liquidStaking.go ./ethereum/abi/liquidStaking
mv pStake-smartContracts/tokenWrapper.go ./ethereum/abi/tokenWrapper

rm -rf pStake-smartContracts

echo "Compilation done."