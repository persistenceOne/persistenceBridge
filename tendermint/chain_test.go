package tendermint

import (
	"github.com/cosmos/relayer/relayer"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestInitializeAndStartChain(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	fileName := strings.Join([]string{dirname,"/.persistenceBridge/chain.json"},"")
	initAndStartChain, errInit := InitializeAndStartChain(fileName, "336h",dirname)
	if errInit != nil {
		t.Errorf("Error Initializing Chain: %v",errInit)
	}
	re := regexp.MustCompile(`^cosmos$`)
	require.Equal(t, true,re.MatchString(initAndStartChain.AccountPrefix))
	require.Equal(t, reflect.TypeOf(&relayer.Chain{}),reflect.TypeOf(initAndStartChain))
	require.NotNil(t,initAndStartChain )
}

func Test_fileInputAdd(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	fileName := strings.Join([]string{dirname,"/.persistenceBridge/chain.json"},"")
	fileIPadd, err := fileInputAdd(fileName)
	if err != nil {
		t.Errorf("Failed Reading config File")
	}
	require.NotNil(t, fileIPadd)
}
