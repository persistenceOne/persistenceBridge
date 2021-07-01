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
	"time"

	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	"github.com/persistenceOne/persistenceBridge/ethereum/magicTx"
	//"math/big"
	"strings"
)

func SendTxToEth(client *ethclient.Client, gasLimit uint64) (string, error) {
	ctx := context.Background()
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		return "", err
	}
	publicKey := casp.GetEthPubKey(uncompressedPublicKeys.PublicKeys[0])
	//publicKey := configuration.GetAppConfig().Ethereum.EthAccountPrivateKey.Public()
	//publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	//if !ok {
	//	log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	//}

	fromAddress := crypto.PubkeyToAddress(publicKey)
	fmt.Println("ETH ADDRESS: " + fromAddress.String())
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return "", err
	}

	contractAddress := common.HexToAddress("0xFe0011F1055dFc307C534142bE4610157Aa42eBD")

	magicTxABI, err := abi.JSON(strings.NewReader(magicTx.MagicTxABI))
	bytesData, err := magicTxABI.Pack("MagicTx", fmt.Sprintf("Abhinav %d", nonce))
	if err != nil {
		return "", err
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Value:    nil,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     bytesData,
		To:       &contractAddress,
	})

	signature, err := getEthSignature(tx, types.HomesteadSigner{})
	if err != nil {
		return "", err
	}

	fmt.Printf("Chain ID: %v\n", chainID.Bytes())
	fmt.Println("ORIGINAL ETH SIGNATURE: " + hex.EncodeToString(signature))
	fmt.Println("ETH SIGNATURE: " + hex.EncodeToString(append(signature, chainID.Bytes()...)))
	fmt.Printf("ETH SIG Length: %d\n", len(append(signature, chainID.Bytes()...)))

	signedTx, err := tx.WithSignature(types.HomesteadSigner{}, append(signature, byte(1)))
	if err != nil {
		return "", err
	}
	//signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, configuration.GetAppConfig().Ethereum.EthAccountPrivateKey)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println(signedTx.)
	//
	err = client.SendTransaction(ctx, signedTx)

	return signedTx.Hash().String(), err
}

func getEthSignature(tx *types.Transaction, signer types.Signer) ([]byte, error) {
	publicKey := []string{"3056301006072A8648CE3D020106052B8104000A03420004B40777F842A9F8BB7ECB94785926D725EB1F96611DC2B2C424EBC8BD1A9B7651302DC7A55301D560D599B3F72D630353325FAED84514C4ECD58330B965A1946A"}

	hash := signer.Hash(tx)
	rlpEncodedTxBytes, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}

	fmt.Println(hash.String())
	signDataResponse, err := caspQueries.SignData([]string{hex.EncodeToString(rlpEncodedTxBytes)}, publicKey)
	if err != nil {
		return nil, err
	}
	fmt.Println("Sleeping for signing tx...")
	time.Sleep(configuration.GetAppConfig().CASP.SignatureWaitTime)
	fmt.Println("AwakeNow...")
	signOperationResponse, err := caspQueries.GetSignOperation(signDataResponse.OperationID)
	if err != nil {
		return nil, err
	}
	signature, err := hex.DecodeString(signOperationResponse.Signatures[0])
	if err != nil {
		return nil, err
	}
	return signature, nil
}
