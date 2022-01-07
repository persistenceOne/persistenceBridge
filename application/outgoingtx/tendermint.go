/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package outgoingtx

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authSigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/tendermint/tendermint/crypto"

	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	caspResponses "github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

// nolint fixme: move into a config or a proper structure type
// nolint: gochecknoglobals
var tmPublicKey cryptotypes.PubKey

// LogMessagesAndBroadcast filters msgs to check repeated withdraw reward message
func LogMessagesAndBroadcast(ctx context.Context, chain *relayer.Chain, msgs []sdk.Msg, timeoutHeight uint64) (*sdk.TxResponse, error) {
	msgsTypes := ""

	for _, msg := range msgs {
		if msg.Type() == bankTypes.TypeMsgSend {
			sendCoin := msg.(*bankTypes.MsgSend)
			msgsTypes = msgsTypes + msg.Type() + " [to: " + sendCoin.ToAddress + " amount: " + sendCoin.Amount.String() + "] "
		} else {
			msgsTypes = msgsTypes + msg.Type() + " "
		}
	}

	logging.Info("Messages to tendermint:", msgsTypes)

	return tendermintSignAndBroadcastMsgs(ctx, chain, msgs, "pStake@PersistenceOne", timeoutHeight)
}

// Timeout height should be greater than current block height or set it 0 for none.
func tendermintSignAndBroadcastMsgs(ctx context.Context, chain *relayer.Chain, msgs []sdk.Msg, memo string, timeoutHeight uint64) (*sdk.TxResponse, error) {
	if tmPublicKey == nil {
		err := setTMPublicKey(ctx)
		if err != nil {
			return nil, err
		}
	}

	bytesToSign, txB, txF, err := getTMBytesToSign(chain, tmPublicKey, msgs, memo, timeoutHeight)
	if err != nil {
		return nil, err
	}

	var signature []byte

	signature, err = getTMSignature(ctx, bytesToSign)
	if err != nil {
		return nil, err
	}

	var txRes *sdk.TxResponse

	txRes, err = broadcastTMTx(chain, tmPublicKey, signature, txB, txF)
	if err != nil {
		return nil, err
	}

	return txRes, err
}

// Timeout height should be greater than current block height or set it 0 for none.
func getTMBytesToSign(chain *relayer.Chain, fromPublicKey cryptotypes.PubKey, msgs []sdk.Msg, memo string, timeoutHeight uint64) ([]byte, client.TxBuilder, *tx.Factory, error) {
	from := sdk.AccAddress(fromPublicKey.Address())
	ctx := chain.CLIContext(0).WithFromAddress(from)

	txFactory, err := tx.PrepareFactory(ctx, chain.TxFactory(0))
	if err != nil {
		return []byte{}, nil, &txFactory, err
	}

	var adjusted uint64

	_, adjusted, err = tx.CalculateGas(ctx.QueryWithData, txFactory, msgs...)
	if err != nil {
		return []byte{}, nil, &txFactory, err
	}

	txFactory = txFactory.WithGas(adjusted).WithMemo(memo).WithTimeoutHeight(timeoutHeight)

	var txBuilder client.TxBuilder

	txBuilder, err = tx.BuildUnsignedTx(txFactory, msgs...)
	if err != nil {
		return []byte{}, nil, &txFactory, err
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

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return []byte{}, txBuilder, &txFactory, err
	}

	var bytesToSign []byte

	bytesToSign, err = ctx.TxConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return []byte{}, txBuilder, &txFactory, err
	}

	return bytesToSign, txBuilder, &txFactory, nil
}

// broadcastTMTx chalk swarm motion broom chapter team guard bracket invest situate circle deny tuition park economy movie subway chase alert popular slogan emerge cricket category
func broadcastTMTx(chain *relayer.Chain, fromPublicKey cryptotypes.PubKey, sigBytes []byte, txBuilder client.TxBuilder, txFactory *tx.Factory) (*sdk.TxResponse, error) {
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

	err := txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	var txBytes []byte

	txBytes, err = ctx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return ctx.BroadcastTx(txBytes)
}

func getTMSignature(ctx context.Context, bytesToSign []byte) ([]byte, error) {
	dataToSign := []string{hex.EncodeToString(crypto.Sha256(bytesToSign))}

	operationID, err := casp.GetCASPSigningOperationID(ctx, dataToSign, []string{configuration.GetAppConfig().CASP.TendermintPublicKey}, "tm")
	if err != nil {
		return nil, err
	}

	var signatureResponse caspResponses.SignOperationResponse

	signatureResponse, err = casp.GetCASPSignature(ctx, operationID)
	if err != nil {
		return nil, err
	}

	if len(signatureResponse.Signatures) == 0 {
		return nil, fmt.Errorf("tendermint %w: ID %s", ErrNoSignature, operationID)
	}

	return hex.DecodeString(signatureResponse.Signatures[0])
}

func setTMPublicKey(ctx context.Context) error {
	if tmPublicKey != nil {
		logging.Warn("outgoingtx: casp tendermint public key already set to.", tmPublicKey.String(), "To change update config and restart.")

		return nil
	}

	logging.Info("outgoingtx: setting tendermint casp public key")

	uncompressedPublicKeys, err := caspQueries.GetUncompressedTMPublicKeys(ctx)
	if err != nil {
		return err
	}

	if len(uncompressedPublicKeys.Items) == 0 {
		logging.Error("no tendermint public keys got from casp")

		return err
	}

	tmPublicKey = casp.GetTMPubKey(uncompressedPublicKeys.Items[0])

	return nil
}
