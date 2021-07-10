package queries

import (
	"fmt"
	tendermintResponses "github.com/persistenceOne/persistenceBridge/application/rest/responses/tendermint"
)

func GetABCI(rpcAddress string) (tendermintResponses.ABCIResponse, error) {
	var abci tendermintResponses.ABCIResponse
	url := rpcAddress + "/abci_info"
	err := get(url, &abci)
	if err != nil {
		return tendermintResponses.ABCIResponse{}, err
	}
	return abci, err
}

func GetTxsByHeight(rpcAddress, height string) (tendermintResponses.TxByHeightResponse, error) {
	var txByHeight tendermintResponses.TxByHeightResponse
	url := rpcAddress + fmt.Sprintf("/tx_search?query=\"tx.height=%s\"", height)
	err := get(url, &txByHeight)
	if err != nil {
		return tendermintResponses.TxByHeightResponse{}, err
	}
	return txByHeight, err
}

func GetTxHash(restAddress, txHash string) (tendermintResponses.TxHashResponse, error) {
	var txHashResponse tendermintResponses.TxHashResponse
	url := restAddress + "/cosmos/tx/v1beta1/txs/" + txHash
	err := get(url, &txHashResponse)
	if err != nil {
		return tendermintResponses.TxHashResponse{}, err
	}
	return txHashResponse, err
}

func GetDelegations(restAddress, accAddress string) (tendermintResponses.DelegationResponse, error) {
	var response tendermintResponses.DelegationResponse
	url := restAddress + "/cosmos/staking/v1beta1/delegations/" + accAddress
	err := get(url, &response)
	if err != nil {
		return tendermintResponses.DelegationResponse{}, err
	}
	return response, err
}
