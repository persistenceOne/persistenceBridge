package application

import (
	"fmt"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	"log"
	"time"
)

func Test(chain *relayer.Chain) {
	uncompressedPublicKeys, err := caspQueries.GetUncompressedPublicKeys()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("CASP PUBLIC KEY (should start with 04): " + uncompressedPublicKeys.PublicKeys[0])
	publicKey := casp.GetPubKey(uncompressedPublicKeys.PublicKeys[0])
	fmt.Printf("Address: %v\n", publicKey.Address())
	fmt.Println("Account Address: " + sdkTypes.AccAddress(publicKey.Address()).String())

	from := sdkTypes.AccAddress(publicKey.Address())
	to, _ := sdkTypes.AccAddressFromBech32("cosmos18gfnyqemvdv7dmqkcyctx2jacg7aswxu5layuq")
	sendMsg := bankTypes.NewMsgSend(from, to, sdkTypes.Coins{sdkTypes.NewCoin("stake", sdkTypes.NewInt(10))})

	bytesToSign, txB, txF, err := outgoingTx.GetTMBytesToSign(chain, publicKey, []sdkTypes.Msg{sendMsg}, "", 0)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Printf("BYTES TO SIGN %x\n", bytesToSign[:])

	signature, err := outgoingTx.GetTMSignature(bytesToSign, 8*time.Second)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("VERIFYING CASP SIGNATURE (expect true): %v\n", publicKey.VerifySignature(bytesToSign, signature))
	txRes, ok, err := outgoingTx.BroadcastTMTx(chain, publicKey, signature, txB, txF)
	if err != nil {
		log.Fatalln(err)
	}
	if !ok {
		log.Fatalf("TX HASH: %s CODE: %d Log: %s\n", txRes.TxHash, txRes.Code, txRes.RawLog)
	}
	log.Println("TX HASH: " + txRes.TxHash)
}
