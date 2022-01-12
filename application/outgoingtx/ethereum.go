/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package outgoingtx

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
	caspResponses "github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

// nolint fixme: move into a config or a proper structure type
// nolint: gochecknoglobals
var ethBridgeAdmin common.Address

func EthereumWrapAndStakeToken(client *ethclient.Client, msgs []db.WrapTokenMsg) (common.Hash, error) {
	if len(msgs) == 0 {
		return common.Hash{}, ErrNoWrapTokenMessages
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
	// fixme: use a proper context with timeout
	ctx := context.Background()

	if ethBridgeAdmin == constants.EthereumZeroAddress() {
		err := setEthBridgeAdmin(ctx)
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

	var (
		caspSignature []byte
		v             int
	)

	caspSignature, v, err = getEthSignature(ctx, tx, signer) // Signature is of 64 bytes, need to append V value
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
func getEthSignature(ctx context.Context, tx *types.Transaction, signer types.Signer) (caspSignature []byte, v int, err error) {
	v = -1

	dataToSign := []string{hex.EncodeToString(signer.Hash(tx).Bytes())}

	var operationID string

	operationID, err = casp.SendDataToSign(ctx, dataToSign, []string{configuration.GetAppConfig().CASP.EthereumPublicKey}, true)
	if err != nil {
		return
	}

	var signatureResponse caspResponses.SignOperationResponse

	signatureResponse, err = casp.GetCASPSignature(ctx, operationID)
	if err != nil {
		return
	}

	if len(signatureResponse.Signatures) == 0 {
		err = fmt.Errorf("ethereum %w: ID %s", ErrNoSignature, operationID)

		return
	}

	caspSignature, err = hex.DecodeString(signatureResponse.Signatures[0])
	if err != nil {
		return
	}

	v = signatureResponse.V[0]

	return
}

func setEthBridgeAdmin(ctx context.Context) error {
	if ethBridgeAdmin != constants.EthereumZeroAddress() {
		logging.Warn("outgoingTx: casp ethereum bridge admin already set to", ethBridgeAdmin.String(), "To change update config and restart")

		return nil
	}

	logging.Info("outgoingtx: setting ethereum bridge admin from casp")

	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys(ctx)
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
