package tendermint

import (
	"fmt"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestGetWrapAddressAndAmounts(t *testing.T) {
	memo := "a,b,c"
	_, _, err := getWrapAddressAndStakingRatio(memo)
	require.Equal(t, fmt.Errorf("invalid memo for bridge"), err)

	memo = ""
	_, _, err = getWrapAddressAndStakingRatio(memo)
	require.Equal(t, fmt.Errorf("invalid memo for bridge"), err)

	memo = constants.DefaultEthZeroAddress
	_, _, err = getWrapAddressAndStakingRatio(memo)
	require.Equal(t, fmt.Errorf("invalid memo for bridge"), err)

	memo = "0x91aaE0aAfd9D2d730111b395c6871f248d7Bd728"
	address, ratio, err := getWrapAddressAndStakingRatio(memo)
	require.Nil(t, err)
	require.Equal(t, memo, address.String())
	require.Equal(t, sdkTypes.ZeroDec(), ratio)

	memo = "0x91aaE0aAfd9D2d730111b395c6871f248d7Bd728,0.5"
	dec, _ := sdkTypes.NewDecFromStr("0.5")
	address, ratio, err = getWrapAddressAndStakingRatio(memo)
	require.Nil(t, err)
	require.Equal(t, strings.Split(memo, ",")[0], address.String())
	require.Equal(t, dec, ratio)

	memo = "0x91aaE0aAfd9D2d730111b395c6871f248d7Bd728,-0.7"
	_, _, err = getWrapAddressAndStakingRatio(memo)
	require.Equal(t, fmt.Errorf("negative ratio: invalid memo for bridge"), err)

	memo = "0x91aaE0aAfd9D2d730111b395c6871f248d7Bd728,0.5,0.4"
	_, _, err = getWrapAddressAndStakingRatio(memo)
	require.Equal(t, fmt.Errorf("invalid memo for bridge"), err)

	memo = "0x91aaE0aAfd9D2d730111b395c6871f248d7Bd728,abcd"
	_, _, err = getWrapAddressAndStakingRatio(memo)
	require.Equal(t, fmt.Errorf("invalid memo for bridge"), err)

}
