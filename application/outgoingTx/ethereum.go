/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package outgoingTx

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	"github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

var ethBridgeAdmin common.Address

func EthereumWrapAndStakeToken(client *ethclient.Client, msgs []db.WrapTokenMsg) (common.Hash, error) {
	if len(msgs) == 0 {
		return common.Hash{}, fmt.Errorf("no wrap token messages to broadcast")
	}
	contractAddress := contracts.LiquidStaking.GetAddress()
	addresses := make([]common.Address, len(msgs))
	stakingAmounts := make([]*big.Int, len(msgs))
	wrappingAmounts := make([]*big.Int, len(msgs))
	for i, msg := range msgs {
		addresses[i] = msg.Address
		stakingAmounts[i] = msg.StakingAmount
		wrappingAmounts[i] = msg.WrapAmount
	}
	bytesData, err := contracts.LiquidStaking.GetABI().Pack("stakeDirectInBatch", addresses, stakingAmounts, wrappingAmounts)
	if err != nil {
		return common.Hash{}, err
	}
	return sendTxToEth(client, &contractAddress, nil, bytesData)
}

func sendTxToEth(client *ethclient.Client, toAddress *common.Address, txValue *big.Int, txData []byte) (common.Hash, error) {
	ctx := context.Background()
	if ethBridgeAdmin.String() == constants.EthereumZeroAddress {
		err := setEthBridgeAdmin()
		if err != nil {
			return common.Hash{}, err
		}
	}
	nonce, err := client.PendingNonceAt(ctx, ethBridgeAdmin)
	if err != nil {
		return common.Hash{}, err
	}

	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: big.NewInt(configuration.GetAppConfig().Ethereum.GasFeeCap),
		GasTipCap: gasTipCap,
		Gas:       configuration.GetAppConfig().Ethereum.GasLimit,
		To:        toAddress,
		Value:     txValue,
		Data:      txData,
	})

	signer := types.NewLondonSigner(chainID)
	caspSignature, v, err := getEthSignature(tx, signer) //Signature is of 64 bytes, need to append V value
	if err != nil {
		return common.Hash{}, err
	}

	signedTx, err := tx.WithSignature(signer, append(caspSignature, byte(v)))
	if err != nil {
		return common.Hash{}, err
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		logging.Error("Broadcasting ETH Tx:", signedTx.Hash().String(), "Error:", err.Error())
	}
	return signedTx.Hash(), err
}

//getEthSignature returns R and S in byte array and V value as int
func getEthSignature(tx *types.Transaction, signer types.Signer) ([]byte, int64, error) {
	dataToSign := []string{hex.EncodeToString(signer.Hash(tx).Bytes())}
	operationID, err := casp.SendDataToSign(dataToSign, []string{configuration.GetAppConfig().CASP.EthereumPublicKey}, true)
	if err != nil {
		return nil, -1, err
	}
	signatureResponse, err := casp.GetCASPSignature(operationID)
	if err != nil {
		return nil, -1, err
	}
	if len(signatureResponse.Signatures) == 0 {
		return nil, -1, fmt.Errorf("ethereum signature not found from casp for operation %s", operationID)
	}
	signature, err := hex.DecodeString(signatureResponse.Signatures[0])
	if err != nil {
		return nil, -1, err
	}
	return signature, signatureResponse.V[0], nil
}

func setEthBridgeAdmin() error {
	if ethBridgeAdmin.String() != constants.EthereumZeroAddress {
		logging.Warn("outgoingTx: casp ethereum bridge admin already set to", ethBridgeAdmin.String(), "To change update config and restart")
		return nil
	}
	logging.Info("outgoingTx: setting ethereum bridge admin from casp")
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		return err
	}
	if len(uncompressedPublicKeys.Items) == 0 {
		logging.Error("no eth public keys got from casp")
		return err
	}
	publicKey := casp.GetEthPubKey(uncompressedPublicKeys.Items[0])
	ethBridgeAdmin = crypto.PubkeyToAddress(publicKey)
	return nil
}
