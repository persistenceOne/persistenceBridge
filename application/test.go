package application

import (
	"encoding/hex"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	"github.com/persistenceOne/persistenceBridge/application/transaction"
	"log"
)

func Test(chain *relayer.Chain) {
	uncompressedPublicKeys, err := caspQueries.GetUncompressedPublicKeys()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("CASP PUBLIC KEY (should start with 04): " + uncompressedPublicKeys.PublicKeys[0])
	publicKey := casp.GetPubKey(uncompressedPublicKeys.PublicKeys[0])
	log.Printf("Address: %v\n", publicKey.Address())
	log.Println("Account Address: " + sdkTypes.AccAddress(publicKey.Address()).String())

	from := sdkTypes.AccAddress(publicKey.Address())
	to, _ := sdkTypes.AccAddressFromBech32("cosmos18gfnyqemvdv7dmqkcyctx2jacg7aswxu5layuq")
	sendMsg := bankTypes.NewMsgSend(from, to, sdkTypes.Coins{sdkTypes.NewCoin("stake", sdkTypes.NewInt(10))})

	bytesToSign, txB, txF, err := transaction.GetBytesToSign(chain, publicKey, []sdkTypes.Msg{sendMsg}, "", 0)
	if err != nil {
		log.Fatalln(err)
	}

	signDataRes, err := caspQueries.SignData([]string{hex.EncodeToString(bytesToSign)}, []string{constants.CASP_PUBLIC_KEY})
	if err != nil {
		log.Fatalln(err)
	}
	signOperationRes, err := caspQueries.GetSignOperation(signDataRes.OperationID)
	if err != nil {
		log.Fatalln(err)
	}
	signature, err := hex.DecodeString(signOperationRes.Signatures[0])
	txRes, ok, err := transaction.BroadcastMsgs(chain, publicKey, signature, txB, txF)
	if err != nil {
		log.Fatalln(err)
	}
	if !ok {
		log.Printf("TX HASH: %s, CODE: %d\n", txRes.TxHash, txRes.Code)
	}
	log.Println("TX HASH: " + txRes.TxHash)
}
