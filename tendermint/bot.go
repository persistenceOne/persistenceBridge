package tendermint

import (
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func handleSlashedOrAboutToBeSlashed(chain *relayer.Chain, validators []db.Validator, processHeight int64, missedBlockCounterForValidator map[string]int64) {
	for _, validator := range validators {
		validatorStatus, err := QueryValidator(validator.Address, chain)
		if err != nil {
			logging.Error("Could not query the validator", processHeight, "Validator Name:", validator.Name, "ERR:", err)
		}

		var pubKey cryptoTypes.PubKey
		err = chain.CLIContext(0).InterfaceRegistry.UnpackAny(validatorStatus.Validator.ConsensusPubkey, &pubKey)
		if err != nil {
			logging.Error("Could not unpack consensus pubKey", processHeight, "ERR:", err)
		}

		consAddress := types.ConsAddress(pubKey.Address())
		slashingInfoAboutValidator, err := QuerySlashingSigningInfo(consAddress, chain)
		if err != nil {
			logging.Error("Could not find the signing info about the validator", processHeight, "Validator Name:", validator.Name, "ERR:", err)
		}

		if slashingInfoAboutValidator.ValSigningInfo.MissedBlocksCounter > missedBlockCounterForValidator[validator.Name] && missedBlockCounterForValidator[validator.Name] != 0 {
			missedBlockCounterForValidator[validator.Name] = slashingInfoAboutValidator.ValSigningInfo.MissedBlocksCounter
			logging.Warn("Validator is about to be jailed", processHeight, "Validator Name:", validator.Name)
		}
		if validatorStatus.Validator.Jailed {
			logging.Warn("Validator is Jailed", processHeight, "Validator Name:", validator.Name)
		}
	}
}
