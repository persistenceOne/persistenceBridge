package queries

import (
	"fmt"
	"github.com/persistenceOne/persistenceCore/pStake/rest/responses/tendermint"
)

func GetABCI(rpcAddress string) (tendermint.ABCIResponse, error) {
	var abci tendermint.ABCIResponse
	url := rpcAddress + "/abci_info"
	err := get(url, &abci)
	if err != nil {
		return tendermint.ABCIResponse{}, err
	}
	return abci, err
}

func GetTxsByHeight(rpcAddress, height string) (tendermint.TxByHeightResponse, error) {
	var txByHeight tendermint.TxByHeightResponse
	url := rpcAddress + fmt.Sprintf("/tx_search?query=\"tx.height=%s\"", height)
	err := get(url, &txByHeight)
	if err != nil {
		return tendermint.TxByHeightResponse{}, err
	}
	return txByHeight, err
}

func GetTxHash(restAddress, txHash string) (tendermint.TxHashResponse, error) {
	var txHashResponse tendermint.TxHashResponse
	url := restAddress + "/cosmos/tx/v1beta1/txs/" + txHash
	err := get(url, &txHashResponse)
	if err != nil {
		return tendermint.TxHashResponse{}, err
	}
	return txHashResponse, err
}

func GetDelegations(restAddress, accAddress string) (tendermint.DelegationResponse, error) {
	var response tendermint.DelegationResponse
	url := restAddress + "/cosmos/staking/v1beta1/delegations/" + accAddress
	err := get(url, &response)
	if err != nil {
		return tendermint.DelegationResponse{}, err
	}
	return response, err
}
