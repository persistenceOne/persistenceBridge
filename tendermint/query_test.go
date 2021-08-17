package tendermint

import (
	"github.com/BurntSushi/toml"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddressIsDelegatorToValidator(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	fileName := strings.Join([]string{dirname,"/.persistenceBridge/chain.json"},"")
	fileInputAddress, err := fileInputAdd(fileName)
	if err != nil {
		t.Errorf("Error Varifying address delegator to validator")
	}
	pStakeConfig := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	tmAddress, err := casp.GetTendermintAddress()
	initAndStartChain, _ := InitializeAndStartChain(fileName, "336h",dirname)
	stakingClient := stakingTypes.NewQueryClient(initAndStartChain.CLIContext(0))
	println(stakingClient)
	configuration.SetPStakeAddress(tmAddress)
	//delegatorToValidator := AddressIsDelegatorToValidator(configuration.GetAppConfig().Tendermint.GetPStakeAddress(),"cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf",fileInputAddress)
	delegatorToValidator := AddressIsDelegatorToValidator("cosmos10snjt8dmpr5my0h76xj48ty80uzwhraqalu4eg","cosmosvaloper1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d76ktteq",fileInputAddress)
	require.Equal(t, true,delegatorToValidator)

}

func TestQueryDelegatorDelegations(t *testing.T) {

}

func TestQueryDelegatorValidatorDelegations(t *testing.T) {

}

func TestQueryValidatorDelegator(t *testing.T) {

}
