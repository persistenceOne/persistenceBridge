package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

var emptyTx common.Hash

func main() {
	fmt.Println(emptyTx.String())
}
