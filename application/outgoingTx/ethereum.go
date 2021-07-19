package outgoingTx

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
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
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/tokenWrapper"
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
	tokenWrapperABI, err := abi.JSON(strings.NewReader(tokenWrapper.TokenWrapperABI))
	if err != nil {
		return common.Hash{}, err
	}
	addresses := make([]common.Address, len(msgs))
	amounts := make([]*big.Int, len(msgs))
	for i, msg := range msgs {
		addresses[i] = msg.Address
		amounts[i] = msg.Amount
	}
	bytesData, err := tokenWrapperABI.Pack("generateUTokensInBatch", addresses, amounts)
	if err != nil {
		return common.Hash{}, err
	}
	return sendTxToEth(client, &contractAddress, nil, bytesData)
}

func sendTxToEth(client *ethclient.Client, toAddress *common.Address, txValue *big.Int, txData []byte) (common.Hash, error) {
	ctx := context.Background()
	if ethBridgeAdmin.String() == "0x0000000000000000000000000000000000000000" {
		setEthBridgeAdmin()
	}
	nonce, err := client.PendingNonceAt(ctx, ethBridgeAdmin)
	if err != nil {
		return common.Hash{}, err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Value:    txValue,
		Gas:      configuration.GetAppConfig().Ethereum.GasLimit,
		GasPrice: gasPrice.Add(gasPrice, big.NewInt(4000000000)),
		Data:     txData,
		To:       toAddress,
	})

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	signer := types.NewEIP155Signer(chainID)
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
func getEthSignature(tx *types.Transaction, signer types.Signer) ([]byte, int, error) {
	dataToSign := []string{hex.EncodeToString(signer.Hash(tx).Bytes())}
	operationID, err := casp.GetCASPSigningOperationID(dataToSign, []string{configuration.GetAppConfig().CASP.EthereumPublicKey}, "eth")
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

func setEthBridgeAdmin() {
	if ethBridgeAdmin.String() != "0x0000000000000000000000000000000000000000" {
		logging.Warn("outgoingTx: casp ethereum bridge admin already set to", ethBridgeAdmin.String(), "To change update config and restart")
		return
	}
	logging.Info("outgoingTx: setting ethereum bridge admin from casp")
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		logging.Fatal(err)
	}
	if len(uncompressedPublicKeys.PublicKeys) == 0 {
		logging.Fatal("no eth public keys got from casp")
	}
	publicKey := casp.GetEthPubKey(uncompressedPublicKeys.PublicKeys[0])
	ethBridgeAdmin = crypto.PubkeyToAddress(publicKey)
}