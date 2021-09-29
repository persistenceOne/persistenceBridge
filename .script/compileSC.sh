#!/bin/bash

solc --abi --bin contracts/LiquidStakingV2.sol -o build
abigen --bin ./build/LiquidStakingV2.bin --abi ./build/LiquidStakingV2.abi --pkg=liquidStaking --out=liquidStaking.go

solc --abi --bin contracts/TokenWrapperV2.sol -o build --overwrite
abigen --bin ./build/TokenWrapperV2.bin --abi ./build/TokenWrapperV2.abi --pkg=tokenWrapper --out=TokenWrapper.go
