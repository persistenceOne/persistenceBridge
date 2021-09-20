#!/bin/bash

solc --abi --bin contracts/LiquidStaking.sol -o build
abigen --bin ./build/LiquidStaking.bin --abi ./build/LiquidStaking.abi --pkg=liquidStaking --out=liquidStaking.go

solc --abi --bin contracts/TokenWrapper.sol -o build --overwrite
abigen --bin ./build/TokenWrapper.bin --abi ./build/TokenWrapper.abi --pkg=tokenWrapper --out=TokenWrapper.go
