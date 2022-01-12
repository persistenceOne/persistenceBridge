//go:build units

package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	Version = "testVersion"
	err := GetVersion().Execute()
	require.Nil(t, err)
}
