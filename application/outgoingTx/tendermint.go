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
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"github.com/tendermint/tendermint/crypto"
)

var tmPublicKey cryptotypes.PubKey

// LogMessagesAndBroadcast filters msgs to check repeated withdraw reward message
func LogMessagesAndBroadcast(chain *relayer.Chain, msgs []sdk.Msg, timeoutHeight uint64) (*sdk.TxResponse, error) {
	msgsTypes := ""
	for _, msg := range msgs {
		if sdk.MsgTypeURL(msg) == constants.MsgSendTypeUrl {
			sendCoin := msg.(*bankTypes.MsgSend)
			msgsTypes = msgsTypes + sendCoin.Type() + " [to: " + sendCoin.ToAddress + " amount: " + sendCoin.Amount.String() + "] "
		} else {
			msgsTypes = msgsTypes + sdk.MsgTypeURL(msg) + " "
		}
	}
	logging.Info("Messages to tendermint:", msgsTypes)
	return tendermintSignAndBroadcastMsgs(chain, msgs, "pStake@PersistenceOne", timeoutHeight)
}

// Timeout height should be greater than current block height or set it 0 for none.
func tendermintSignAndBroadcastMsgs(chain *relayer.Chain, msgs []sdk.Msg, memo string, timeoutHeight uint64) (*sdk.TxResponse, error) {
	if tmPublicKey == nil {
		err := setTMPublicKey()
		if err != nil {
			return nil, err
		}
	}
	bytesToSign, txB, txF, err := getTMBytesToSign(chain, tmPublicKey, msgs, memo, timeoutHeight)
	if err != nil {
		return nil, err
	}
	signature, err := getTMSignature(bytesToSign)
	if err != nil {
		return nil, err
	}

	txRes, err := broadcastTMTx(chain, tmPublicKey, signature, txB, txF)
	if err != nil {
		return nil, err
	}
	return txRes, err
}

// Timeout height should be greater than current block height or set it 0 for none.
func getTMBytesToSign(chain *relayer.Chain, fromPublicKey cryptotypes.PubKey, msgs []sdk.Msg, memo string, timeoutHeight uint64) ([]byte, client.TxBuilder, tx.Factory, error) {

	from := sdk.AccAddress(fromPublicKey.Address())
	ctx := chain.CLIContext(0).WithFromAddress(from)

	txFactory, err := prepareFactory(ctx, chain.TxFactory(0))
	if err != nil {
		return []byte{}, nil, txFactory, err
	}

	_, adjusted, err := tx.CalculateGas(ctx, txFactory, msgs...)
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
	operationID, err := casp.GetCASPSigningOperationID(dataToSign, []string{configuration.GetAppConfig().CASP.TendermintPublicKey}, "tm")
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

func setTMPublicKey() error {
	if tmPublicKey != nil {
		logging.Warn("outgoingTx: casp tendermint public key already set to.", tmPublicKey.String(), "To change update config and restart.")
		return nil
	}
	logging.Info("outgoingTx: setting tendermint casp public key")
	uncompressedPublicKeys, err := caspQueries.GetUncompressedTMPublicKeys()
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

func prepareFactory(clientCtx client.Context, txf tx.Factory) (tx.Factory, error) {
	from := clientCtx.GetFromAddress()

	if err := txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
}
