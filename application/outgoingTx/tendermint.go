package outgoingTx

import (
	"encoding/hex"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authSigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	"github.com/tendermint/tendermint/crypto"
)

// Timeout height should be greater than current block height or set it 0 for none.
func SignAndBroadcastTM(chain *relayer.Chain, msgs []sdk.Msg, memo string, timeoutHeight uint64) (*sdk.TxResponse, bool, error) {
	uncompressedPublicKeys, err := caspQueries.GetUncompressedPublicKeys()
	if err != nil {
		return nil, false, err
	}
	publicKey := casp.GetPubKey(uncompressedPublicKeys.PublicKeys[0])

	bytesToSign, txB, txF, err := getTMBytesToSign(chain, publicKey, msgs, memo, timeoutHeight)
	if err != nil {
		return nil, false, err
	}

	signature, err := getTMSignature(bytesToSign, 8*time.Second)
	if err != nil {
		return nil, false, err
	}

	txRes, ok, err := broadcastTMTx(chain, publicKey, signature, txB, txF)
	if err != nil {
		return nil, false, err
	}
	return txRes, ok, err
}

// Timeout height should be greater than current block height or set it 0 for none.
func getTMBytesToSign(chain *relayer.Chain, fromPublicKey cryptotypes.PubKey, msgs []sdk.Msg, memo string, timeoutHeight uint64) ([]byte, client.TxBuilder, tx.Factory, error) {

	from := sdk.AccAddress(fromPublicKey.Address())
	ctx := chain.CLIContext(0).WithFromAddress(from)

	txFactory, err := tx.PrepareFactory(ctx, chain.TxFactory(0))
	if err != nil {
		return []byte{}, nil, txFactory, err
	}

	_, adjusted, err := tx.CalculateGas(ctx.QueryWithData, txFactory, msgs...)
	if err != nil {
		return []byte{}, nil, txFactory, err
	}

	txFactory = txFactory.WithGas(adjusted).WithMemo(memo).WithTimeoutHeight(timeoutHeight)

	txBuilder, err := tx.BuildUnsignedTx(txFactory, msgs...)
	if err != nil {
		return []byte{}, nil, txFactory, err
	}

	signMode := txFactory.SignMode()
	if signMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		signMode = ctx.TxConfig.SignModeHandler().DefaultMode()
	}

	signerData := authSigning.SignerData{
		ChainID:       txFactory.ChainID(),
		AccountNumber: txFactory.AccountNumber(),
		Sequence:      txFactory.Sequence(),
	}
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   fromPublicKey,
		Data:     &sigData,
		Sequence: txFactory.Sequence(),
	}
	if err := txBuilder.SetSignatures(sig); err != nil {
		return []byte{}, txBuilder, txFactory, err
	}

	bytesToSign, err := ctx.TxConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return []byte{}, txBuilder, txFactory, err
	}

	return bytesToSign, txBuilder, txFactory, nil
}

// broadcastTMTx chalk swarm motion broom chapter team guard bracket invest situate circle deny tuition park economy movie subway chase alert popular slogan emerge cricket category
func broadcastTMTx(chain *relayer.Chain, fromPublicKey cryptotypes.PubKey, sigBytes []byte, txBuilder client.TxBuilder, txFactory tx.Factory) (*sdk.TxResponse, bool, error) {

	from := sdk.AccAddress(fromPublicKey.Address())
	ctx := chain.CLIContext(0).WithFromAddress(from).WithBroadcastMode("async")

	signMode := txFactory.SignMode()
	if signMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		signMode = ctx.TxConfig.SignModeHandler().DefaultMode()
	}

	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
	}
	sig := signing.SignatureV2{
		PubKey:   fromPublicKey,
		Data:     &sigData,
		Sequence: txFactory.Sequence(),
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, false, err
	}

	txBytes, err := ctx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, false, err
	}

	res, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		return nil, false, err
	}
	if res.Code != 0 {
		return res, false, nil
	}

	log.Printf("TX HASH: %s, CODE: %d\n", res.TxHash, res.Code)

	return res, true, nil
}

func getTMSignature(bytesToSign []byte, signatureWaitTime time.Duration) ([]byte, error) {
	signDataResponse, err := caspQueries.SignData([]string{hex.EncodeToString(crypto.Sha256(bytesToSign))}, []string{configuration.GetAppConfig().CASP.PublicKey})
	if err != nil {
		return nil, err
	}
	time.Sleep(signatureWaitTime)
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
