package commands

import (
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestInitCommand(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	cmd := InitCommand()
	cmd.SetArgs([]string{dirname + "/Documents/GitHub/persistenceBridge"})
	err := cmd.Execute()
	require.Equal(t, nil, err)
	config := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &config)
	require.Equal(t, nil, err)

}
