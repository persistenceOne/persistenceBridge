package blockchain

import (
	"fmt"
	"github.com/persistenceOne/persistenceBridge/application/rest"
	"github.com/persistenceOne/persistenceBridge/application/rest/responses/blockchain"
)

func GetABCI(rpcAddress string) (blockchain.ABCIResponse, error) {
	var abci blockchain.ABCIResponse
	url := rpcAddress + "/abci_info"
	err := rest.Get(url, &abci)
	if err != nil {
		return blockchain.ABCIResponse{}, err
	}
	return abci, err
}

func GetTxsByHeight(rpcAddress, height string) (blockchain.TxByHeightResponse, error) {
	var txByHeight blockchain.TxByHeightResponse
	url := rpcAddress + fmt.Sprintf("/tx_search?query=\"tx.height=%s\"", height)
	err := rest.Get(url, &txByHeight)
	if err != nil {
		return blockchain.TxByHeightResponse{}, err
	}
	return txByHeight, err
}

func GetTxHash(restAddress, txHash string) (blockchain.TxHashResponse, error) {
	var txHashResponse blockchain.TxHashResponse
	url := restAddress + "/cosmos/tx/v1beta1/txs/" + txHash
	err := rest.Get(url, &txHashResponse)
	if err != nil {
		return blockchain.TxHashResponse{}, err
	}
	return txHashResponse, err
}

func GetDelegations(restAddress, accAddress string) (blockchain.DelegationResponse, error) {
	var response blockchain.DelegationResponse
	url := restAddress + "/cosmos/staking/v1beta1/delegations/" + accAddress
	err := rest.Get(url, &response)
	if err != nil {
		return blockchain.DelegationResponse{}, err
	}
	return response, err
}
