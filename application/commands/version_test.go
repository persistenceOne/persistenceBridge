package commands

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetVersion(t *testing.T) {
	Version = "testVersion"
	err := GetVersion().Execute()
	require.Nil(t, err)
}
