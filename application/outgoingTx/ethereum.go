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
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	caspResponses "github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

var ethBridgeAdmin common.Address

type WrapTokenMsg struct {
	Address common.Address `json:"address"`
	Amount  *big.Int       `json:"amount"`
}

func EthereumWrapToken(client *ethclient.Client, msgs []WrapTokenMsg) (common.Hash, error) {
	if len(msgs) == 0 {
		return common.Hash{}, fmt.Errorf("no wrap token messages to broadcast")
	}

	contractAddress := common.HexToAddress(constants.TokenWrapperAddress)

	tokenWrapperABI, err := abi.JSON(strings.NewReader(constants.TokenWrapperABI))
	if err != nil {
		return common.Hash{}, err
	}

	addresses := make([]common.Address, len(msgs))
	amounts := make([]*big.Int, len(msgs))

	for i, msg := range msgs {
		addresses[i] = msg.Address
		amounts[i] = msg.Amount
	}

	var bytesData []byte

	bytesData, err = tokenWrapperABI.Pack("generateUTokensInBatch", addresses, amounts)
	if err != nil {
		return common.Hash{}, err
	}

	return sendTxToEth(client, &contractAddress, nil, bytesData)
}

func sendTxToEth(client *ethclient.Client, toAddress *common.Address, txValue *big.Int, txData []byte) (common.Hash, error) {
	ctx := context.Background()
	if ethBridgeAdmin.String() == "0x0000000000000000000000000000000000000000" {
		err := setEthBridgeAdmin()
		if err != nil {
			return common.Hash{}, err
		}
	}

	nonce, err := client.PendingNonceAt(ctx, ethBridgeAdmin)
	if err != nil {
		return common.Hash{}, err
	}

	var gasTipCap *big.Int

	gasTipCap, err = client.SuggestGasTipCap(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	var chainID *big.Int

	chainID, err = client.ChainID(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	// TODO set it as conf parameter
	gasFeeCap := big.NewInt(300000000000)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       configuration.GetAppConfig().Ethereum.GasLimit,
		To:        toAddress,
		Value:     txValue,
		Data:      txData,
	})

	signer := types.NewLondonSigner(chainID)

	var (
		caspSignature []byte
		v             int
	)

	caspSignature, v, err = getEthSignature(tx, signer) // Signature is of 64 bytes, need to append V value
	if err != nil {
		return common.Hash{}, err
	}

	var signedTx *types.Transaction

	signedTx, err = tx.WithSignature(signer, append(caspSignature, byte(v)))
	if err != nil {
		return common.Hash{}, err
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		logging.Error("Broadcasting ETH Tx:", signedTx.Hash().String(), "Error:", err.Error())
	}

	return signedTx.Hash(), err
}

// getEthSignature returns R and S in byte array and V value as int
func getEthSignature(tx *types.Transaction, signer types.Signer) (caspSignature []byte, v int, err error) {
	v = -1

	dataToSign := []string{hex.EncodeToString(signer.Hash(tx).Bytes())}

	var operationID string

	operationID, err = casp.GetCASPSigningOperationID(dataToSign, []string{configuration.GetAppConfig().CASP.EthereumPublicKey}, "eth")
	if err != nil {
		return
	}

	var signatureResponse caspResponses.SignOperationResponse

	signatureResponse, err = casp.GetCASPSignature(operationID)
	if err != nil {
		return
	}

	if len(signatureResponse.Signatures) == 0 {
		err = fmt.Errorf("ethereum signature not found from casp for operation %s", operationID)

		return
	}

	caspSignature, err = hex.DecodeString(signatureResponse.Signatures[0])
	if err != nil {
		return
	}

	v = signatureResponse.V[0]

	return
}

func setEthBridgeAdmin() error {
	if ethBridgeAdmin.String() != "0x0000000000000000000000000000000000000000" {
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
