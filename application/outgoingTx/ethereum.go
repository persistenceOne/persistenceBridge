package outgoingTx

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/tokenWrapper"
	//"github.com/persistenceOne/persistenceBridge/ethereum/magicTx"
	"log"
	"math/big"
	"strings"
	"time"

	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
)

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

//// EthereumMagicTx
////TODO Delete
//func EthereumMagicTx(client *ethclient.Client) (common.Hash, error) {
//
//	contractAddress := common.HexToAddress("0xFe0011F1055dFc307C534142bE4610157Aa42eBD")
//
//	magicTxABI, err := abi.JSON(strings.NewReader(magicTx.MagicTxABI))
//	bytesData, err := magicTxABI.Pack("MagicTx", fmt.Sprintf("Abhinav"))
//	if err != nil {
//		return common.Hash{}, err
//	}
//
//	return sendTxToEth(client, &contractAddress, nil, bytesData)
//
//}

func sendTxToEth(client *ethclient.Client, contractAddress *common.Address, txValue *big.Int, txData []byte) (common.Hash, error) {
	ctx := context.Background()
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		return common.Hash{}, err
	}
	publicKey := casp.GetEthPubKey(uncompressedPublicKeys.PublicKeys[0])

	fromAddress := crypto.PubkeyToAddress(publicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return common.Hash{}, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Value:    txValue,
		Gas:      configuration.GetAppConfig().Ethereum.EthGasLimit,
		GasPrice: gasPrice,
		Data:     txData,
		To:       contractAddress,
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

	//sender, err := signer.Sender(signedTx)
	//if err != nil {
	//	return "", err
	//}
	//if sender.String() != fromAddress.String() {
	//	return "", fmt.Errorf("invalid signer")
	//}
	log.Printf("Broadcasting ETH Tx: %s\n", signedTx.Hash().String())
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Printf("ERROR Broadcasting ETH Tx: %s, Error: %s\n", signedTx.Hash().String(), err.Error())
	}
	return signedTx.Hash(), err
}

//getEthSignature returns R and S in byte array and V value as int
func getEthSignature(tx *types.Transaction, signer types.Signer) ([]byte, int, error) {
	signDataResponse, err := caspQueries.SignData([]string{hex.EncodeToString(signer.Hash(tx).Bytes())}, []string{configuration.GetAppConfig().CASP.EthPublicKey})
	if err != nil {
		return nil, -1, err
	}
	time.Sleep(configuration.GetAppConfig().CASP.SignatureWaitTime)
	signOperationResponse, err := caspQueries.GetSignOperation(signDataResponse.OperationID)
	if err != nil {
		return nil, -1, err
	}
	signature, err := hex.DecodeString(signOperationResponse.Signatures[0])
	if err != nil {
		return nil, -1, err
	}
	return signature, signOperationResponse.V[0], nil
}
