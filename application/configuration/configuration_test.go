package configuration

import (
	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/go-cmp/cmp"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestGetAppConfig(t *testing.T) {
	t.Logf("Testing GetAppConfig()")
	var newAppConfig *config
	var appConfig = GetAppConfig()
	cmp.Equal(newAppConfig, appConfig)
}

func TestInitConfig(t *testing.T) {
	t.Logf("Testing InitConfig()")
	newAppConfigAddress := InitConfig()
	newAppConfig := *newAppConfigAddress
	appConfig := newConfig()

	if !(unsafe.Sizeof(appConfig) == unsafe.Sizeof(newAppConfig)) {
		t.Errorf("Not initalized properly")
	}

	if !cmp.Equal(newAppConfig.Kafka, appConfig.Kafka) {
		t.Errorf("Kafka config not initalized properly")
	}

	if !cmp.Equal(newAppConfig.Ethereum, appConfig.Ethereum) {
		t.Errorf("Ethereum config not initalized properly")
	}

	if !cmp.Equal(newAppConfig.TelegramBot, appConfig.TelegramBot) {
		t.Errorf("TelegramBot config not initalized properly")
	}

	if !cmp.Equal(newAppConfig.seal, appConfig.seal) {
		t.Errorf("seal not initalized properly")
	}

	if !cmp.Equal(newAppConfig.RPCEndpoint, appConfig.RPCEndpoint) {
		t.Errorf("RPCEndpoints not initalized properly")
	}

	if !cmp.Equal(newAppConfig.Tendermint.PStakeDenom, appConfig.Tendermint.PStakeDenom) {
		t.Errorf("PStakeDenom not initalized properly")
	}

	if !cmp.Equal(newAppConfig.Tendermint.MinimumWrapAmount, appConfig.Tendermint.MinimumWrapAmount) {
		t.Errorf("MinimumWrapAmount not initalized properly")
	}

	if !cmp.Equal(newAppConfig.Tendermint.BroadcastMode, appConfig.Tendermint.BroadcastMode) {
		t.Errorf("BroadcastMode not initalized properly")
	}
}

func TestSetConfig(t *testing.T) {
	appConfig = InitConfig()
	var cmd cobra.Command

	cmd.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	cmd.Flags().String(constants2.FlagEthereumEndPoint, constants2.DefaultEthereumEndPoint, "ethereum orchestrator to connect")
	cmd.Flags().String(constants2.FlagKafkaPorts, constants2.DefaultKafkaPorts, "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	cmd.Flags().String(constants2.FlagDenom, constants2.DefaultDenom, "denom name")
	cmd.Flags().Uint64(constants2.FlagEthGasLimit, constants2.DefaultEthGasLimit, "Gas limit for eth txs")
	cmd.Flags().String(constants2.FlagBroadcastMode, constants2.DefaultBroadcastMode, "broadcast mode for tendermint")
	cmd.Flags().String(constants2.FlagCASPURL, "https://65.2.149.241:443", "casp api url (with http)")
	cmd.Flags().String(constants2.FlagCASPVaultID, "4ec017bf-4af8-41b3-9527-a466e05971cb", "casp vault id")
	cmd.Flags().String(constants2.FlagCASPApiToken, "Bearer cHVuZWV0TmV3QXBpa2V5MTI6OWM1NDBhMzAtNTQ5NC00ZDdhLTljODktODA3MDZiNWNhYzQ1", "casp api token (in format: Bearer ...)")
	cmd.Flags().String(constants2.FlagCASPTMPublicKey, "3056301006072A8648CE3D020106052B8104000A0342000413109ECEADCBF6122EF44184B207F8C6820E509497792DDFB166BC090A0FB4447CFFCE16BAAF9EC7F57D14C02641B3A6A698614D973ED744E725A85E62535DA4", "casp tendermint public key")
	cmd.Flags().String(constants2.FlagCASPEthPublicKey, "3056301006072A8648CE3D020106052B8104000A034200049D8BB9DC3E37511273286F60C989BFFC3E28909F426AF7D4A7899FACC4E3DB00413E2DA7A8CF33F367D8C4D8FC2BFA791DD4389CC1E75154CD38429FD9525946", "casp ethereum public key")
	cmd.Flags().Int(constants2.FlagCASPSignatureWaitTime, int(constants2.DefaultCASPSignatureWaitTime.Seconds()), "casp signature wait time")
	cmd.Flags().Bool(constants2.FlagCASPConcurrentKey, true, "allows starting multiple sign operations that specify the same key")
	cmd.Flags().String(constants2.FlagRPCEndpoint, constants2.DefaultRPCEndpoint, "rpc Endpoint for server")
	cmd.Flags().Int64(constants2.FlagMinimumWrapAmount, constants2.DefaultMinimumWrapAmount, "minimum amount in send coin tx to wrap onto eth")
	cmd.Flags().String(constants2.FlagTelegramBotToken, "xyz", "telegram bot token")
	cmd.Flags().Int64(constants2.FlagTelegramChatID, 0, "telegram chat id")

	appConfig := SetConfig(&cmd)
	if appConfig.Tendermint.PStakeDenom == "" {
		t.Errorf("FlagDenom not set")
	}
	if appConfig.Ethereum.EthereumEndPoint == "" {
		t.Errorf("EthereumEndPoint not set")
	}
	if appConfig.Ethereum.GasLimit == 0 {
		t.Errorf("GasLimit not set")
	}
	if appConfig.Kafka.Brokers[0] == "" {
		t.Errorf("Brokers not set")
	}
	if appConfig.Tendermint.BroadcastMode == "" {
		t.Errorf("BroadcastMode not set")
	}
	if appConfig.CASP.URL == "" {
		t.Errorf("FlagDenom not set")
	}
	if appConfig.Tendermint.PStakeDenom == "" {
		t.Errorf("CASP URL not set")
	}
	if appConfig.CASP.VaultID == "" {
		t.Errorf("CASP VaultID not set")
	}
	if appConfig.CASP.APIToken == "" {
		t.Errorf("CASP APItoken not set")
	}
	if appConfig.CASP.TendermintPublicKey == "" {
		t.Errorf("TendermintPublicKey not set")
	}
	if appConfig.CASP.EthereumPublicKey == "" {
		t.Errorf("EthereumPublicKey not set")
	}
	if appConfig.CASP.SignatureWaitTime < 0 {
		t.Errorf("CASP signature wait time not set")
	}
	if appConfig.CASP.EthereumPublicKey == "" {
		t.Errorf("EthereumPublicKey not set")
	}
	if appConfig.CASP.AllowConcurrentKeyUsage == false {
		t.Errorf("AllowConcurrentKeyUsage not set")
	}
	if appConfig.RPCEndpoint == "" {
		t.Errorf("RPCEndpoint not set")
	}
	if appConfig.Tendermint.MinimumWrapAmount < 0 {
		t.Errorf("MinimumWrapAmount not set")
	}
	if appConfig.TelegramBot.Token == "" {
		t.Errorf("RPCEndpoint not set")
	}
	if appConfig.TelegramBot.ChatID == 1 {
		t.Errorf("RPCEndpoint not set")
	}
}

func TestSetPStakeAddress(t *testing.T) {
	pStakeConfig := InitConfig()
	dirname, _ := os.UserHomeDir()
	_, err := toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	opaddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetPStakeAddress(opaddress)
	if appConfig.Tendermint.pStakeAddress != opaddress.String() {
		t.Errorf("PStakeAddress not set")
	}
}

func TestValidateAndSeal(t *testing.T) {
	pStakeConfig := InitConfig()
	dirname, _ := os.UserHomeDir()
	_, err := toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	opaddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetPStakeAddress(opaddress)
	ValidateAndSeal()
	if appConfig.seal != true {
		t.Errorf("appConfig did not get validated")
	}
}

func TestGetPStakeAddress(t *testing.T) {
	pStakeConfig := InitConfig()
	dirname, _ := os.UserHomeDir()
	_, err := toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	opaddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetPStakeAddress(opaddress)
	PstakeAddress := pStakeConfig.Tendermint.GetPStakeAddress()
	if PstakeAddress == "" {
		log.Fatalln("pStakeAddress not set")
	}
}
