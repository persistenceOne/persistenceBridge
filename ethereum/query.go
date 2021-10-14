package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func QueryEtherBalance(address common.Address, client *ethclient.Client) (big.Int, error) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return big.Int{}, err
	}
	return *balance, nil
}
