#!/bin/bash

OPENZEPPELIN_DIR_NAME="@openzeppelin"
LIQUID_STAKING=LiquidStaking
TOKEN_WRAPPER=TokenWrapper
OPENZEPPELIN_VERSION=v3.3.0

rm -rf stkATOM-ERC20
rm -rf "$OPENZEPPELIN_DIR_NAME"
rm -rf ./ethereum/abi/liquidStaking/liquidStaking.go
rm -rf ./ethereum/abi/tokenWrapper/tokenWrapper.go

echo "Downloading stkATOM-ERC20..."
git clone git@github.com:persistenceOne/stkATOM-ERC20.git
# shellcheck disable=SC2164
cd stkATOM-ERC20
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
mv "$OPENZEPPELIN_DIR_NAME" ./stkATOM-ERC20/
# shellcheck disable=SC2164
cd stkATOM-ERC20

echo "Compiling smart contracts..."
# shellcheck disable=SC2046
solc --abi --bin contracts/"$LIQUID_STAKING".sol -o build @openzeppelin/=$(pwd)/@openzeppelin/
abigen --bin ./build/"$LIQUID_STAKING".bin --abi ./build/"$LIQUID_STAKING".abi --pkg=liquidStaking --out=liquidStaking.go

# shellcheck disable=SC2046
solc --abi --bin contracts/"$TOKEN_WRAPPER".sol -o build --overwrite @openzeppelin/=$(pwd)/@openzeppelin/
abigen --bin ./build/"$TOKEN_WRAPPER".bin --abi ./build/"$TOKEN_WRAPPER".abi --pkg=tokenWrapper --out=tokenWrapper.go

cd ..
mv stkATOM-ERC20/liquidStaking.go ./ethereum/abi/liquidStaking
mv stkATOM-ERC20/tokenWrapper.go ./ethereum/abi/tokenWrapper

rm -rf stkATOM-ERC20

echo "Compilation done."
