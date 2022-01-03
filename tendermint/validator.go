package tendermint

import (
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/relayer"

	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func CheckValidators(chain *relayer.Chain, processHeight int64) {
	// Get validators list from db
	validators, err := db.GetValidators()
	if err != nil {
		logging.Error("Could not fetch validators from DB", processHeight, "ERR:", err)
	}

	missedBlockCounterForValidator := make(map[string]int64, len(validators))

	for _, validator := range validators {
		// Query details about validator
		var validatorDetails stakingTypes.QueryValidatorResponse

		validatorDetails, err = QueryValidator(validator.Address, chain)
		if err != nil {
			logging.Error("Could not query the validator", processHeight, "Validator Name:", validator.Name, "ERR:", err)
		}

		// Convert type.Any key to consensus pubKey
		var pubKey cryptoTypes.PubKey

		err = chain.CLIContext(0).InterfaceRegistry.UnpackAny(validatorDetails.Validator.ConsensusPubkey, &pubKey)
		if err != nil {
			logging.Error("Could not unpack consensus pubKey", processHeight, "ERR:", err)
		}

		// Convert PubKey to consAddress to query Slashing Info
		consAddress := types.ConsAddress(pubKey.Address())

		var validatorSlashingInfo slashingTypes.QuerySigningInfoResponse

		validatorSlashingInfo, err = QuerySlashingSigningInfo(consAddress, chain)
		if err != nil {
			logging.Error("Could not find the signing info about the validator", processHeight, "Validator Name:", validator.Name, "ERR:", err)
		}

		if validatorSlashingInfo.ValSigningInfo.MissedBlocksCounter > missedBlockCounterForValidator[validator.Name] {
			logging.Warn("Validator is about to be jailed", processHeight, "Validator Name:", validator.Name)
		}

		if validatorDetails.Validator.Jailed {
			logging.Warn("Validator is Jailed", processHeight, "Validator Name:", validator.Name)
		}

		missedBlockCounterForValidator[validator.Name] = validatorSlashingInfo.ValSigningInfo.MissedBlocksCounter
	}
}
