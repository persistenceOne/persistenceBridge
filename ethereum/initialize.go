package ethereum

import (
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/liquidStaking"
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/tokenWrapper"
	contracts2 "github.com/persistenceOne/persistenceBridge/ethereum/contracts"
)

func init() {
	contracts2.LiquidStaking.SetABI(liquidStaking.LiquidStakingMetaData.ABI)
	contracts2.TokenWrapper.SetABI(tokenWrapper.TokenWrapperMetaData.ABI)
}
