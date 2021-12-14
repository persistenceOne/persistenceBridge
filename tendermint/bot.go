package tendermint

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"time"
)

func allValidatorsMap(chain *relayer.Chain) (map[string]types.ConsAddress, error) {
	allValidators, _ := QueryAllValidators(chain)
	m := make(map[string]types.ConsAddress)
	for _, val := range allValidators {
		m[string(val.PubKey.Bytes())] = val.Address
	}
	return m, nil
}

func handleSlashedOrAboutToBeSlashed(chain *relayer.Chain, validators []db.Validator, processHeight int64, maxMissed int64) error {
	validatorsMap, err := allValidatorsMap(chain)
	if err != nil {
		return err
	}
	for _, validator := range validators {
		validatorStatus, err := QueryValidatorStatus(validator.Address, chain)
		if err != nil {
			return err
		}

		compareBytes := validatorStatus.Validator.ConsensusPubkey.Value[2:]
		if len(validatorsMap[string(compareBytes)]) > 0 {
			slashingInfoAboutValidator, err := QuerySlashingSigningInfo(validatorsMap[string(compareBytes)], chain)
			if err != nil {
				logging.Info("Could not find the signing info about the validator", processHeight, "Validator Name:", validator.Name, "ERR:", err)
			}

			if slashingInfoAboutValidator.ValSigningInfo.MissedBlocksCounter > maxMissed {
				logging.Error("Validator is about to be jailed", processHeight, "Validator Name:", validator.Name)
			}

			if time.Now().Before(slashingInfoAboutValidator.ValSigningInfo.JailedUntil) {
				logging.Error("Validator is already Jailed", processHeight, "Validator Name:", validator.Name)
			}
		} else {
			logging.Error("Validator might already be jailed", processHeight, "Validator Name:", validator.Name)
		}
	}
	return nil
}
