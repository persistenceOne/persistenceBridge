package outgoingTx

import (
	"encoding/hex"
	"fmt"
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
	"log"
)

// FilterMessagesAndBroadcast filters msgs to check repeated withdraw reward message
func FilterMessagesAndBroadcast(chain *relayer.Chain, msgs []sdk.Msg, timeoutHeight uint64) (*sdk.TxResponse, error) {
	var filteredMsgs []sdk.Msg
	msgsTypes := ""
	messageHash := make(map[string]bool)
	for _, msg := range msgs {
		if !messageHash[hex.EncodeToString(crypto.Sha256(msg.GetSignBytes()))] {
			filteredMsgs = append(filteredMsgs, msg)
			messageHash[hex.EncodeToString(crypto.Sha256(msg.GetSignBytes()))] = true
			msgsTypes = msgsTypes + msg.Type() + " "
		}
	}
	log.Println("Messages to tendermint: " + msgsTypes)
	return tendermintSignAndBroadcastMsgs(chain, filteredMsgs, "pStake@PersistenceOne", timeoutHeight)
}

// Timeout height should be greater than current block height or set it 0 for none.
func tendermintSignAndBroadcastMsgs(chain *relayer.Chain, msgs []sdk.Msg, memo string, timeoutHeight uint64) (*sdk.TxResponse, error) {
	uncompressedPublicKeys, err := caspQueries.GetUncompressedTMPublicKeys()
	if err != nil {
		return nil, err
	}
	if len(uncompressedPublicKeys.PublicKeys) == 0 {
		return nil, fmt.Errorf("no public keys got from casp")
	}
	publicKey := casp.GetTMPubKey(uncompressedPublicKeys.PublicKeys[0])
	bytesToSign, txB, txF, err := getTMBytesToSign(chain, publicKey, msgs, memo, timeoutHeight)
	if err != nil {
		return nil, err
	}

	signature, err := getTMSignature(bytesToSign)
	if err != nil {
		return nil, err
	}

	txRes, err := broadcastTMTx(chain, publicKey, signature, txB, txF)
	if err != nil {
		return nil, err
	}
	return txRes, err
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
func broadcastTMTx(chain *relayer.Chain, fromPublicKey cryptotypes.PubKey, sigBytes []byte, txBuilder client.TxBuilder, txFactory tx.Factory) (*sdk.TxResponse, error) {

	from := sdk.AccAddress(fromPublicKey.Address())
	ctx := chain.CLIContext(0).WithFromAddress(from).WithBroadcastMode(configuration.GetAppConfig().Tendermint.BroadcastMode)

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
		return nil, err
	}

	txBytes, err := ctx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	res, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		return nil, err
	}

	return res, err
}

func getTMSignature(bytesToSign []byte) ([]byte, error) {
	dataToSign := []string{hex.EncodeToString(crypto.Sha256(bytesToSign))}
	operationID, err := casp.GetCASPSigningOperationID(dataToSign, []string{configuration.GetAppConfig().CASP.TendermintPublicKey})
	if err != nil {
		return nil, err
	}
	signatureResponse, err := casp.GetCASPSignature(operationID)
	if err != nil {
		return nil, err
	}
	if len(signatureResponse.Signatures) == 0 {
		return nil, fmt.Errorf("tendermint signature not found from casp for operation %s", operationID)
	}
	signature, err := hex.DecodeString(signatureResponse.Signatures[0])
	if err != nil {
		return nil, err
	}
	return signature, nil
}
