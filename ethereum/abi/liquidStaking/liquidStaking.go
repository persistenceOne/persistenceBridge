// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package liquidStaking

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// LiquidStakingMetaData contains all meta data concerning the LiquidStaking contract.
var LiquidStakingMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"batchingLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"SetBatchingLimit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakeFee\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"unstakeFee\",\"type\":\"uint256\"}],\"name\":\"SetFees\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"minStake\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"minUnstake\",\"type\":\"uint256\"}],\"name\":\"SetMinimumValues\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"SetSTokensContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"SetTokenWrapperContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"SetUTokensContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"unstakeEpoch\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"unstakeEpochPrevious\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"epochInterval\",\"type\":\"uint256\"}],\"name\":\"SetUnstakeEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"unstakingLockTime\",\"type\":\"uint256\"}],\"name\":\"SetUnstakingLockTime\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakingAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountStaked\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wrappingAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountWrapped\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"StakeDirect\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"toAddressses\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"amountsStaked\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"wrappingAmounts\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"amountsWrapped\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"StakeDirectInBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"finalTokens\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"StakeTokens\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"finalTokens\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"UnstakeTokens\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"WithdrawUnstakeTokens\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BRIDGE_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_batchingLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_sTokens\",\"outputs\":[{\"internalType\":\"contractISTokensV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_tokenWrapperContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_uTokens\",\"outputs\":[{\"internalType\":\"contractIUTokensV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_unstakingAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_unstakingExpiration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_valueDivisor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_withdrawCounters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBatchingLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"batchingLimit\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStakeUnstakeProps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"stakeFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minUnstake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"valueDivisor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochInterval\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeEpochPrevious\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakingLockTime\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokens\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"uTokenAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"getTotalUnbondedTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"unbondingTokens\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"getTotalUnbondingTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"unbondingTokens\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_unstakeTimestamp\",\"type\":\"uint256\"}],\"name\":\"getUnstakeEpochMilestone\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"unstakeEpochMilestone\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_unstakeTimestamp\",\"type\":\"uint256\"}],\"name\":\"getUnstakeTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"unstakeTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeEpochPrevious\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"uAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pauserAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"unstakingLockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochInterval\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"valueDivisor\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchingLimit\",\"type\":\"uint256\"}],\"name\":\"setBatchingLimit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"stakeFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeFee\",\"type\":\"uint256\"}],\"name\":\"setFees\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minUnstake\",\"type\":\"uint256\"}],\"name\":\"setMinimumValues\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sAddress\",\"type\":\"address\"}],\"name\":\"setSTokensContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenWrapperContract\",\"type\":\"address\"}],\"name\":\"setTokenWrapperContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"uAddress\",\"type\":\"address\"}],\"name\":\"setUTokensContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"unstakeEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeEpochPrevious\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochInterval\",\"type\":\"uint256\"}],\"name\":\"setUnstakeEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"unstakingLockTime\",\"type\":\"uint256\"}],\"name\":\"setUnstakingLockTime\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountStaked\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"stakingAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"wrappingAmount\",\"type\":\"uint256\"}],\"name\":\"stakeDirect\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountStaked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountWrapped\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"toAddressses\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"wrappingAmounts\",\"type\":\"uint256[]\"}],\"name\":\"stakeDirectInBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amountsStaked\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amountsWrapped\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountUnstaked\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"withdrawUnstakedTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"withdrawBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50614f1f806100206000396000f3fe608060405234801561001057600080fd5b50600436106102695760003560e01c80638456cb5911610151578063b3deb34f116100c3578063d79cf55a11610087578063d79cf55a14610dc3578063d939381414610df7578063e321fca614610e59578063e63ab1e914610ecc578063f3e6f94e14610eea578063f75f0cf614610f4c57610269565b8063b3deb34f14610c9d578063bef1307c14610cbb578063ca15c87314610cff578063d061adab14610d41578063d547741f14610d7557610269565b80639b9eea46116101155780639b9eea4614610ada5780639dc77cab14610b1c578063a217fddf14610b74578063aa6ca80814610b92578063adc9772e14610be3578063b17e316514610c4557610269565b80638456cb591461094c57806385e328401461096c578063898a9664146109bc5780639010d07c14610a1457806391d1485414610a7657610269565b806331b188cc116101ea57806355f2323f116101ae57806355f2323f1461078a5780635c975abb146107e0578063728a99ee14610800578063728cdbca1461081e57806378a2f074146108c0578063819020a0146108f457610269565b806331b188cc1461066c57806336568abe146106b05780633e118dbe146106fe5780633f4ba83a1461071c57806351db00e71461073c57610269565b806316027ed01161023157806316027ed01461053657806321bdf9e514610554578063248a9ca3146105985780632e31fab8146105da5780632f2ff15d1461061e57610269565b80630273b70d1461026e5780630b78f9c0146102c65780630e4b29431461031457806310ef62f614610358578063118c38c714610518575b600080fd5b6102b06004803603602081101561028457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610fae565b6040518082815260200191505060405180910390f35b6102fc600480360360408110156102dc57600080fd5b810190808035906020019092919080359060200190929190505050611140565b60405180821515815260200191505060405180910390f35b6103406004803603602081101561032a57600080fd5b81019080803590602001909291905050506112d0565b60405180821515815260200191505060405180910390f35b6104796004803603606081101561036e57600080fd5b810190808035906020019064010000000081111561038b57600080fd5b82018360208201111561039d57600080fd5b803590602001918460208302840111640100000000831117156103bf57600080fd5b9091929391929390803590602001906401000000008111156103e057600080fd5b8201836020820111156103f257600080fd5b8035906020019184602083028401116401000000008311171561041457600080fd5b90919293919293908035906020019064010000000081111561043557600080fd5b82018360208201111561044757600080fd5b8035906020019184602083028401116401000000008311171561046957600080fd5b9091929391929390505050611395565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b838110156104c05780820151818401526020810190506104a5565b50505050905001838103825284818151815260200191508051906020019060200280838360005b838110156105025780820151818401526020810190506104e7565b5050505090500194505050505060405180910390f35b61052061178d565b6040518082815260200191505060405180910390f35b61053e6117b1565b6040518082815260200191505060405180910390f35b6105966004803603602081101561056a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506117b7565b005b6105c4600480360360208110156105ae57600080fd5b81019080803590602001909291905050506118c4565b6040518082815260200191505060405180910390f35b61061c600480360360208110156105f057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506118e4565b005b61066a6004803603604081101561063457600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506119f1565b005b6106ae6004803603602081101561068257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611a7b565b005b6106fc600480360360408110156106c657600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611b88565b005b610706611c21565b6040518082815260200191505060405180910390f35b610724611c27565b60405180821515815260200191505060405180910390f35b6107726004803603604081101561075257600080fd5b810190808035906020019092919080359060200190929190505050611cdb565b60405180821515815260200191505060405180910390f35b610792611e97565b604051808a8152602001898152602001888152602001878152602001868152602001858152602001848152602001838152602001828152602001995050505050505050505060405180910390f35b6107e8611edd565b60405180821515815260200191505060405180910390f35b610808611ef4565b6040518082815260200191505060405180910390f35b6108be600480360360c081101561083457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019092919080359060200190929190505050611efa565b005b6108c8612087565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6109366004803603602081101561090a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506120ad565b6040518082815260200191505060405180910390f35b6109546126b0565b60405180821515815260200191505060405180910390f35b6109986004803603602081101561098257600080fd5b8101908080359060200190929190505050612764565b60405180848152602001838152602001828152602001935050505060405180910390f35b6109fe600480360360208110156109d257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506127b2565b6040518082815260200191505060405180910390f35b610a4a60048036036040811015610a2a57600080fd5b8101908080359060200190929190803590602001909291905050506127ca565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610ac260048036036040811015610a8c57600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506127fc565b60405180821515815260200191505060405180910390f35b610b0660048036036020811015610af057600080fd5b810190808035906020019092919050505061282e565b6040518082815260200191505060405180910390f35b610b5c60048036036060811015610b3257600080fd5b810190808035906020019092919080359060200190929190803590602001909291905050506128fa565b60405180821515815260200191505060405180910390f35b610b7c612acb565b6040518082815260200191505060405180910390f35b610b9a612ad2565b604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390f35b610c2f60048036036040811015610bf957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050612b23565b6040518082815260200191505060405180910390f35b610c8760048036036020811015610c5b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050613039565b6040518082815260200191505060405180910390f35b610ca56131cc565b6040518082815260200191505060405180910390f35b610ce760048036036020811015610cd157600080fd5b81019080803590602001909291905050506131d6565b60405180821515815260200191505060405180910390f35b610d2b60048036036020811015610d1557600080fd5b81019080803590602001909291905050506132a6565b6040518082815260200191505060405180910390f35b610d496132cd565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610dc160048036036040811015610d8b57600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506132f3565b005b610dcb61337d565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610e4360048036036040811015610e0d57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506133a3565b6040518082815260200191505060405180910390f35b610eaf60048036036060811015610e6f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190803590602001909291905050506138b6565b604051808381526020018281526020019250505060405180910390f35b610ed46139ec565b6040518082815260200191505060405180910390f35b610f3660048036036040811015610f0057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050613a10565b6040518082815260200191505060405180910390f35b610f9860048036036040811015610f6257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050613a41565b6040518082815260200191505060405180910390f35b60008060a260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020805490509050600060a460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905060008190505b828110156111385760006110a760a260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020838154811061109757fe5b9060005260206000200154612764565b5050905080421061111c5761111960a360008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002083815481106110ff57fe5b906000526020600020015486613a7290919063ffffffff16565b94505b50611131600182613a7290919063ffffffff16565b9050611040565b505050919050565b60006111566000801b611151613afa565b6127fc565b6111c8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260038152602001807f4c5131000000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b6111de6064609d54613b0290919063ffffffff16565b831115806111ec5750600083145b801561121857506112096064609d54613b0290919063ffffffff16565b821115806112175750600082145b5b61128a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260038152602001807f4c5132000000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b82609b8190555081609c8190555081837f3be6e637f54f081a2dd12982cfd58481304fc98b604b18fa8ed4aa4e4a9d153260405160405180910390a36001905092915050565b60006112e66000801b6112e1613afa565b6127fc565b611358576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260038152602001807f4c5133000000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b81609e81905550817ffced4450741ad45114dbc966176eb1d16dc4f2259bc864aa395736a7e8355bbb60405160405180910390a260019050919050565b6060806113c97f751b795d24b92e3d92d1d0d8f2885f4e9c9c269da350af36ae6b49069babf4bf6113c4613afa565b6127fc565b61143b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5132380000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b858590508888905014801561145557508383905086869050145b6114c7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5132390000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b6000806000808b8b905090508067ffffffffffffffff811180156114ea57600080fd5b506040519080825280602002602001820160405280156115195781602001602082028036833780820191505090505b5095508067ffffffffffffffff8111801561153357600080fd5b506040519080825280602002602001820160405280156115625781602001602082028036833780820191505090505b509450600091505b80821015611610576115c98c8c8481811061158157fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff168b8b858181106115aa57fe5b905060200201358a8a868181106115bd57fe5b90506020020135613b88565b8094508195505050838683815181106115de57fe5b602002602001018181525050828583815181106115f757fe5b602002602001018181525050818060010192505061156a565b7fae1ff5caf707784d6d12e6dd90a042ce2994fc6f3d32634ea155dd6e352ed6cc8c8c8c8c8a8d8d8c4260405180806020018060200180602001806020018060200187815260200186810386528f8f82818152602001925060200280828437600081840152601f19601f82011690508083019250505086810385528d8d82818152602001925060200280828437600081840152601f19601f82011690508083019250505086810384528b818151815260200191508051906020019060200280838360005b838110156116ef5780820151818401526020810190506116d4565b5050505090500186810383528a8a82818152602001925060200280828437600081840152601f19601f820116905080830192505050868103825288818151815260200191508051906020019060200280838360005b8381101561175f578082015181840152602081019050611744565b505050509050019e50505050505050505050505050505060405180910390a150505050965096945050505050565b7f751b795d24b92e3d92d1d0d8f2885f4e9c9c269da350af36ae6b49069babf4bf81565b60a65481565b6117cb6000801b6117c6613afa565b6127fc565b61183d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5131300000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b80609760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff167ff745c285f36f88bea6af14d4d0f33dd9350cef7895216c1615b1caaee7857e0c60405160405180910390a250565b600060656000838152602001908152602001600020600201549050919050565b6118f86000801b6118f3613afa565b6127fc565b61196a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5131310000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b80609860006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff167fd0b9dbe26f676f4cdb02fc249c203baa0314cec79e18b4194dc862795a5de03f60405160405180910390a250565b611a186065600084815260200190815260200160002060020154611a13613afa565b6127fc565b611a6d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f815260200180614e0d602f913960400191505060405180910390fd5b611a778282613f55565b5050565b611a8f6000801b611a8a613afa565b6127fc565b611b01576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5132370000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b8060a760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff167f2f9f272420570ea07d7da789fdc59a50c36792d6ea4f49a4d218b6f2431d049260405160405180910390a250565b611b90613afa565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614611c13576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f815260200180614ebb602f913960400191505060405180910390fd5b611c1d8282613fe9565b5050565b60a55481565b6000611c5a7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a611c55613afa565b6127fc565b611ccc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5132330000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b611cd461407d565b6001905090565b6000611cf16000801b611cec613afa565b6127fc565b611d63576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260038152602001807f4c5134000000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b6001831015611dda576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260038152602001807f4c5135000000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b6001821015611e51576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260038152602001807f4c5136000000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b8260998190555081609a8190555081837fa94d40347449ea19b873d64ed55c99772194c49791e82be5c2dc24e4c0044fe360405160405180910390a36001905092915050565b6000806000806000806000806000609b549850609c5497506099549650609a549550609d549450609f54935060a054925060a1549150609e549050909192939495969798565b6000603360009054906101000a900460ff16905090565b609d5481565b600060019054906101000a900460ff1680611f195750611f18614168565b5b80611f2f575060008054906101000a900460ff16155b611f84576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180614e6c602e913960400191505060405180910390fd5b60008060019054906101000a900460ff161590508015611fd4576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b611fdc614179565b611fe4614287565b611ff86000801b611ff3613afa565b614395565b6120227f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a86614395565b61202b876117b7565b612034866118e4565b61203d846112d0565b50612049600180611cdb565b5081609d8190555061205c4242856128fa565b50801561207e5760008060016101000a81548160ff0219169083151502179055505b50505050505050565b609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60006120b7611edd565b1561212a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f5061757361626c653a207061757365640000000000000000000000000000000081525060200191505060405180910390fd5b612132613afa565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16146121d2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5132300000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b600060a460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050600060a460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905060006122718360a654613a7290919063ffffffff16565b60a260008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002080549050116123015760a260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002080549050612317565b6123168360a654613a7290919063ffffffff16565b5b905060008390505b818110156124cf57600061238660a260008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020838154811061237657fe5b9060005260206000200154612764565b505090508042106124b3576123f860a360008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002083815481106123de57fe5b906000526020600020015487613a7290919063ffffffff16565b955060a260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020828154811061244457fe5b906000526020600020016000905560a360008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020828154811061249c57fe5b906000526020600020016000905583806001019450505b506124c8600182613a7290919063ffffffff16565b905061231f565b5060008411612546576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5132310000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b838573ffffffffffffffffffffffffffffffffffffffff167f31453b59c0eb931edc765cc2400876a524254457b5534d660c9da1d74ab4666b426040518082815260200191505060405180910390a38160a460008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340c10f1986866040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561266c57600080fd5b505af1158015612680573d6000803e3d6000fd5b505050506040513d602081101561269657600080fd5b810190808051906020019092919050505050505050919050565b60006126e37f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a6126de613afa565b6127fc565b612755576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5132320000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b61275d6143a3565b6001905090565b6000806000806127738561282e565b905060008114156127885760009350506127ab565b60a054925060a15491506127a7609e5482613a7290919063ffffffff16565b9350505b9193909250565b60a46020528060005260406000206000915090505481565b60006127f4826065600086815260200190815260200160002060000161448f90919063ffffffff16565b905092915050565b600061282682606560008681526020019081526020016000206000016144a990919063ffffffff16565b905092915050565b60008082141561284157600090506128f5565b600060a0541480156128555750600060a154145b8061286257506000609f54145b1561286f578190506128f5565b8160a05411156128835760a05490506128f5565b60006128ae61289f60a15460a0546144d990919063ffffffff16565b84613a7290919063ffffffff16565b905060006128c760a054836144d990919063ffffffff16565b90506128f0826128e2609f548461455c90919063ffffffff16565b613a7290919063ffffffff16565b925050505b919050565b60006129106000801b61290b613afa565b6127fc565b612982576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260038152602001807f4c5137000000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b838311156129f8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260038152602001807f4c5138000000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b600084148015612a09575060008214155b15612a7c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260038152602001807f4c5139000000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b8360a0819055508260a18190555081609f819055508183857f342db083975cc0697d0f81d0876299ecbf95d909e920426089062e540904497c60405160405180910390a4600190509392505050565b6000801b81565b600080609860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169150609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690509091565b6000612b2d611edd565b15612ba0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f5061757361626c653a207061757365640000000000000000000000000000000081525060200191505060405180910390fd5b612ba8613afa565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614612c48576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5131320000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b609954821015612cc0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5131330000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b6000609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231856040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b158015612d4b57600080fd5b505afa158015612d5f573d6000803e3d6000fd5b505050506040513d6020811015612d7557600080fd5b8101908080519060200190929190505050905082811015612dfe576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5131340000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b6000612e2c6064612e1e609b54609d54886145e59092919063ffffffff16565b6146c590919063ffffffff16565b9050612e4181856144d990919063ffffffff16565b925083838673ffffffffffffffffffffffffffffffffffffffff167f5fe79871cd2431c06447cbcf2557091da5d2ed5bc640f1028f42665913786e42426040518082815260200191505060405180910390a4609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639dc29fac86866040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015612f2657600080fd5b505af1158015612f3a573d6000803e3d6000fd5b505050506040513d6020811015612f5057600080fd5b810190808051906020019092919050505050609860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340c10f1986856040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015612ff557600080fd5b505af1158015613009573d6000803e3d6000fd5b505050506040513d602081101561301f57600080fd5b810190808051906020019092919050505050505092915050565b60008060a260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020805490509050600060a460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905060008190505b828110156131c457600061313260a260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020838154811061312257fe5b9060005260206000200154612764565b50509050804210156131a8576131a560a360008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020838154811061318b57fe5b906000526020600020015486613a7290919063ffffffff16565b94505b506131bd600182613a7290919063ffffffff16565b90506130cb565b505050919050565b600060a654905090565b60006131ec6000801b6131e7613afa565b6127fc565b61325e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5132340000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b8160a681905550817f3e0887a905af0c5764f2c25c201fa1b9595065cb7008f54fe7734fefdd8e88dc426040518082815260200191505060405180910390a260019050919050565b60006132c66065600084815260200190815260200160002060000161474e565b9050919050565b60a760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b61331a6065600084815260200190815260200160002060020154613315613afa565b6127fc565b61336f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526030815260200180614e3c6030913960400191505060405180910390fd5b6133798282613fe9565b5050565b609860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60006133ad611edd565b15613420576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f5061757361626c653a207061757365640000000000000000000000000000000081525060200191505060405180910390fd5b613428613afa565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16146134c8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5131350000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b609a54821015613540576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5131380000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b6000609860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231856040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156135cb57600080fd5b505afa1580156135df573d6000803e3d6000fd5b505050506040513d60208110156135f557600080fd5b810190808051906020019092919050505090508281101561367e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5131390000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b60006136ac606461369e609c54609d54886145e59092919063ffffffff16565b6146c590919063ffffffff16565b90506136c181856144d990919063ffffffff16565b925083838673ffffffffffffffffffffffffffffffffffffffff167f98850af53deeaf736beb3ea22091c900e8e5cfdf8153342daa366dc1d83a39fd426040518082815260200191505060405180910390a4609860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639dc29fac86866040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b1580156137a657600080fd5b505af11580156137ba573d6000803e3d6000fd5b505050506040513d60208110156137d057600080fd5b81019080805190602001909291905050505060a260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002042908060018154018082558091505060019003906000526020600020016000909190919091505560a360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020839080600181540180825580915050600190039060005260206000200160009091909190915055505092915050565b6000806138ea7f751b795d24b92e3d92d1d0d8f2885f4e9c9c269da350af36ae6b49069babf4bf6138e5613afa565b6127fc565b61395c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5132350000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b613967858585613b88565b80925081935050507f7fd980aad29d3d6256218a13a2d1d8dac5cf5d7411c1c122e2503923952eddd5858584868542604051808773ffffffffffffffffffffffffffffffffffffffff168152602001868152602001858152602001848152602001838152602001828152602001965050505050505060405180910390a1935093915050565b7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a81565b60a26020528160005260406000208181548110613a2c57600080fd5b90600052602060002001600091509150505481565b60a36020528160005260406000208181548110613a5d57600080fd5b90600052602060002001600091509150505481565b600080828401905083811015613af0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b600033905090565b600080831415613b155760009050613b82565b6000828402905082848281613b2657fe5b0414613b7d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526021815260200180614e9a6021913960400191505060405180910390fd5b809150505b92915050565b600080600073ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1614158015613bd55750600084148015613bd35750600083145b155b613c47576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260048152602001807f4c5132360000000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b600080600080613c796064613c6b609b54609d548c6145e59092919063ffffffff16565b6146c590919063ffffffff16565b9350613c8e84896144d990919063ffffffff16565b955060a760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ec1893b46040518163ffffffff1660e01b815260040160a06040518083038186803b158015613cf857600080fd5b505afa158015613d0c573d6000803e3d6000fd5b505050506040513d60a0811015613d2257600080fd5b8101908080519060200190929190805190602001909291908051906020019092919080519060200190929190805190602001909291905050509091925090915090508092508193505050613d946064613d8684848b6145e59092919063ffffffff16565b6146c590919063ffffffff16565b9250613da983886144d990919063ffffffff16565b9450609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340c10f198a876040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015613e3e57600080fd5b505af1158015613e52573d6000803e3d6000fd5b505050506040513d6020811015613e6857600080fd5b810190808051906020019092919050505050609860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340c10f198a886040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015613f0d57600080fd5b505af1158015613f21573d6000803e3d6000fd5b505050506040513d6020811015613f3757600080fd5b81019080805190602001909291905050505050505050935093915050565b613f7d816065600085815260200190815260200160002060000161476390919063ffffffff16565b15613fe557613f8a613afa565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b614011816065600085815260200190815260200160002060000161479390919063ffffffff16565b156140795761401e613afa565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b60405160405180910390a45b5050565b614085611edd565b6140f7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f5061757361626c653a206e6f742070617573656400000000000000000000000081525060200191505060405180910390fd5b6000603360006101000a81548160ff0219169083151502179055507f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa61413b613afa565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1565b6000614173306147c3565b15905090565b600060019054906101000a900460ff16806141985750614197614168565b5b806141ae575060008054906101000a900460ff16155b614203576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180614e6c602e913960400191505060405180910390fd5b60008060019054906101000a900460ff161590508015614253576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b61425b6147d6565b6142636148d4565b80156142845760008060016101000a81548160ff0219169083151502179055505b50565b600060019054906101000a900460ff16806142a657506142a5614168565b5b806142bc575060008054906101000a900460ff16155b614311576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180614e6c602e913960400191505060405180910390fd5b60008060019054906101000a900460ff161590508015614361576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b6143696147d6565b6143716149d2565b80156143925760008060016101000a81548160ff0219169083151502179055505b50565b61439f8282613f55565b5050565b6143ab611edd565b1561441e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f5061757361626c653a207061757365640000000000000000000000000000000081525060200191505060405180910390fd5b6001603360006101000a81548160ff0219169083151502179055507f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258614462613afa565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1565b600061449e8360000183614aeb565b60001c905092915050565b60006144d1836000018373ffffffffffffffffffffffffffffffffffffffff1660001b614b6e565b905092915050565b600082821115614551576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b818303905092915050565b60008082116145d3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f536166654d6174683a206d6f64756c6f206279207a65726f000000000000000081525060200191505060405180910390fd5b8183816145dc57fe5b06905092915050565b60008060006145f48686614b91565b915091506000848061460257fe5b868809905082811115614616576001820391505b808303925060008214156146385784838161462d57fe5b0493505050506146be565b8482106146ad576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f46756c6c4d6174683a2046554c4c4449565f4f564552464c4f5700000000000081525060200191505060405180910390fd5b6146b8838387614be4565b93505050505b9392505050565b600080821161473c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f536166654d6174683a206469766973696f6e206279207a65726f00000000000081525060200191505060405180910390fd5b81838161474557fe5b04905092915050565b600061475c82600001614c81565b9050919050565b600061478b836000018373ffffffffffffffffffffffffffffffffffffffff1660001b614c92565b905092915050565b60006147bb836000018373ffffffffffffffffffffffffffffffffffffffff1660001b614d02565b905092915050565b600080823b905060008111915050919050565b600060019054906101000a900460ff16806147f557506147f4614168565b5b8061480b575060008054906101000a900460ff16155b614860576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180614e6c602e913960400191505060405180910390fd5b60008060019054906101000a900460ff1615905080156148b0576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b80156148d15760008060016101000a81548160ff0219169083151502179055505b50565b600060019054906101000a900460ff16806148f357506148f2614168565b5b80614909575060008054906101000a900460ff16155b61495e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180614e6c602e913960400191505060405180910390fd5b60008060019054906101000a900460ff1615905080156149ae576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b80156149cf5760008060016101000a81548160ff0219169083151502179055505b50565b600060019054906101000a900460ff16806149f157506149f0614168565b5b80614a07575060008054906101000a900460ff16155b614a5c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180614e6c602e913960400191505060405180910390fd5b60008060019054906101000a900460ff161590508015614aac576001600060016101000a81548160ff02191690831515021790555060016000806101000a81548160ff0219169083151502179055505b6000603360006101000a81548160ff0219169083151502179055508015614ae85760008060016101000a81548160ff0219169083151502179055505b50565b600081836000018054905011614b4c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180614deb6022913960400191505060405180910390fd5b826000018281548110614b5b57fe5b9060005260206000200154905092915050565b600080836001016000848152602001908152602001600020541415905092915050565b60008060007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80614bbe57fe5b84860990508385029250828103915082811015614bdc576001820391505b509250929050565b6000808260000383169050808381614bf857fe5b049250808581614c0457fe5b0494506001818260000381614c1557fe5b04018402850194506000600190508084026002038102905080840260020381029050808402600203810290508084026002038102905080840260020381029050808402600203810290508084026002038102905080840260020381029050808602925050509392505050565b600081600001805490509050919050565b6000614c9e8383614b6e565b614cf7578260000182908060018154018082558091505060019003906000526020600020016000909190919091505582600001805490508360010160008481526020019081526020016000208190555060019050614cfc565b600090505b92915050565b60008083600101600084815260200190815260200160002054905060008114614dde5760006001820390506000600186600001805490500390506000866000018281548110614d4d57fe5b9060005260206000200154905080876000018481548110614d6a57fe5b9060005260206000200181905550600183018760010160008381526020019081526020016000208190555086600001805480614da257fe5b60019003818190600052602060002001600090559055866001016000878152602001908152602001600020600090556001945050505050614de4565b60009150505b9291505056fe456e756d657261626c655365743a20696e646578206f7574206f6620626f756e6473416363657373436f6e74726f6c3a2073656e646572206d75737420626520616e2061646d696e20746f206772616e74416363657373436f6e74726f6c3a2073656e646572206d75737420626520616e2061646d696e20746f207265766f6b65496e697469616c697a61626c653a20636f6e747261637420697320616c726561647920696e697469616c697a6564536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636520726f6c657320666f722073656c66a2646970667358221220d3517cf7ae32b27ec973d508c1ed77edf0717c838a4bdff6515b41829117f5ea64736f6c63430007060033",
}

// LiquidStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use LiquidStakingMetaData.ABI instead.
var LiquidStakingABI = LiquidStakingMetaData.ABI

// LiquidStakingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LiquidStakingMetaData.Bin instead.
var LiquidStakingBin = LiquidStakingMetaData.Bin

// DeployLiquidStaking deploys a new Ethereum contract, binding an instance of LiquidStaking to it.
func DeployLiquidStaking(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LiquidStaking, error) {
	parsed, err := LiquidStakingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LiquidStakingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LiquidStaking{LiquidStakingCaller: LiquidStakingCaller{contract: contract}, LiquidStakingTransactor: LiquidStakingTransactor{contract: contract}, LiquidStakingFilterer: LiquidStakingFilterer{contract: contract}}, nil
}

// LiquidStaking is an auto generated Go binding around an Ethereum contract.
type LiquidStaking struct {
	LiquidStakingCaller     // Read-only binding to the contract
	LiquidStakingTransactor // Write-only binding to the contract
	LiquidStakingFilterer   // Log filterer for contract events
}

// LiquidStakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type LiquidStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidStakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LiquidStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidStakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LiquidStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidStakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LiquidStakingSession struct {
	Contract     *LiquidStaking    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LiquidStakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LiquidStakingCallerSession struct {
	Contract *LiquidStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// LiquidStakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LiquidStakingTransactorSession struct {
	Contract     *LiquidStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// LiquidStakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type LiquidStakingRaw struct {
	Contract *LiquidStaking // Generic contract binding to access the raw methods on
}

// LiquidStakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LiquidStakingCallerRaw struct {
	Contract *LiquidStakingCaller // Generic read-only contract binding to access the raw methods on
}

// LiquidStakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LiquidStakingTransactorRaw struct {
	Contract *LiquidStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLiquidStaking creates a new instance of LiquidStaking, bound to a specific deployed contract.
func NewLiquidStaking(address common.Address, backend bind.ContractBackend) (*LiquidStaking, error) {
	contract, err := bindLiquidStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LiquidStaking{LiquidStakingCaller: LiquidStakingCaller{contract: contract}, LiquidStakingTransactor: LiquidStakingTransactor{contract: contract}, LiquidStakingFilterer: LiquidStakingFilterer{contract: contract}}, nil
}

// NewLiquidStakingCaller creates a new read-only instance of LiquidStaking, bound to a specific deployed contract.
func NewLiquidStakingCaller(address common.Address, caller bind.ContractCaller) (*LiquidStakingCaller, error) {
	contract, err := bindLiquidStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingCaller{contract: contract}, nil
}

// NewLiquidStakingTransactor creates a new write-only instance of LiquidStaking, bound to a specific deployed contract.
func NewLiquidStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*LiquidStakingTransactor, error) {
	contract, err := bindLiquidStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingTransactor{contract: contract}, nil
}

// NewLiquidStakingFilterer creates a new log filterer instance of LiquidStaking, bound to a specific deployed contract.
func NewLiquidStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*LiquidStakingFilterer, error) {
	contract, err := bindLiquidStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingFilterer{contract: contract}, nil
}

// bindLiquidStaking binds a generic wrapper to an already deployed contract.
func bindLiquidStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LiquidStakingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LiquidStaking *LiquidStakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidStaking.Contract.LiquidStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LiquidStaking *LiquidStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidStaking.Contract.LiquidStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LiquidStaking *LiquidStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidStaking.Contract.LiquidStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LiquidStaking *LiquidStakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LiquidStaking *LiquidStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LiquidStaking *LiquidStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidStaking.Contract.contract.Transact(opts, method, params...)
}

// BRIDGEADMINROLE is a free data retrieval call binding the contract method 0x118c38c7.
//
// Solidity: function BRIDGE_ADMIN_ROLE() view returns(bytes32)
func (_LiquidStaking *LiquidStakingCaller) BRIDGEADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "BRIDGE_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BRIDGEADMINROLE is a free data retrieval call binding the contract method 0x118c38c7.
//
// Solidity: function BRIDGE_ADMIN_ROLE() view returns(bytes32)
func (_LiquidStaking *LiquidStakingSession) BRIDGEADMINROLE() ([32]byte, error) {
	return _LiquidStaking.Contract.BRIDGEADMINROLE(&_LiquidStaking.CallOpts)
}

// BRIDGEADMINROLE is a free data retrieval call binding the contract method 0x118c38c7.
//
// Solidity: function BRIDGE_ADMIN_ROLE() view returns(bytes32)
func (_LiquidStaking *LiquidStakingCallerSession) BRIDGEADMINROLE() ([32]byte, error) {
	return _LiquidStaking.Contract.BRIDGEADMINROLE(&_LiquidStaking.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_LiquidStaking *LiquidStakingCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_LiquidStaking *LiquidStakingSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _LiquidStaking.Contract.DEFAULTADMINROLE(&_LiquidStaking.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_LiquidStaking *LiquidStakingCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _LiquidStaking.Contract.DEFAULTADMINROLE(&_LiquidStaking.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_LiquidStaking *LiquidStakingCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_LiquidStaking *LiquidStakingSession) PAUSERROLE() ([32]byte, error) {
	return _LiquidStaking.Contract.PAUSERROLE(&_LiquidStaking.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_LiquidStaking *LiquidStakingCallerSession) PAUSERROLE() ([32]byte, error) {
	return _LiquidStaking.Contract.PAUSERROLE(&_LiquidStaking.CallOpts)
}

// BatchingLimit is a free data retrieval call binding the contract method 0x16027ed0.
//
// Solidity: function _batchingLimit() view returns(uint256)
func (_LiquidStaking *LiquidStakingCaller) BatchingLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "_batchingLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BatchingLimit is a free data retrieval call binding the contract method 0x16027ed0.
//
// Solidity: function _batchingLimit() view returns(uint256)
func (_LiquidStaking *LiquidStakingSession) BatchingLimit() (*big.Int, error) {
	return _LiquidStaking.Contract.BatchingLimit(&_LiquidStaking.CallOpts)
}

// BatchingLimit is a free data retrieval call binding the contract method 0x16027ed0.
//
// Solidity: function _batchingLimit() view returns(uint256)
func (_LiquidStaking *LiquidStakingCallerSession) BatchingLimit() (*big.Int, error) {
	return _LiquidStaking.Contract.BatchingLimit(&_LiquidStaking.CallOpts)
}

// STokens is a free data retrieval call binding the contract method 0xd79cf55a.
//
// Solidity: function _sTokens() view returns(address)
func (_LiquidStaking *LiquidStakingCaller) STokens(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "_sTokens")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// STokens is a free data retrieval call binding the contract method 0xd79cf55a.
//
// Solidity: function _sTokens() view returns(address)
func (_LiquidStaking *LiquidStakingSession) STokens() (common.Address, error) {
	return _LiquidStaking.Contract.STokens(&_LiquidStaking.CallOpts)
}

// STokens is a free data retrieval call binding the contract method 0xd79cf55a.
//
// Solidity: function _sTokens() view returns(address)
func (_LiquidStaking *LiquidStakingCallerSession) STokens() (common.Address, error) {
	return _LiquidStaking.Contract.STokens(&_LiquidStaking.CallOpts)
}

// TokenWrapperContract is a free data retrieval call binding the contract method 0xd061adab.
//
// Solidity: function _tokenWrapperContract() view returns(address)
func (_LiquidStaking *LiquidStakingCaller) TokenWrapperContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "_tokenWrapperContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenWrapperContract is a free data retrieval call binding the contract method 0xd061adab.
//
// Solidity: function _tokenWrapperContract() view returns(address)
func (_LiquidStaking *LiquidStakingSession) TokenWrapperContract() (common.Address, error) {
	return _LiquidStaking.Contract.TokenWrapperContract(&_LiquidStaking.CallOpts)
}

// TokenWrapperContract is a free data retrieval call binding the contract method 0xd061adab.
//
// Solidity: function _tokenWrapperContract() view returns(address)
func (_LiquidStaking *LiquidStakingCallerSession) TokenWrapperContract() (common.Address, error) {
	return _LiquidStaking.Contract.TokenWrapperContract(&_LiquidStaking.CallOpts)
}

// UTokens is a free data retrieval call binding the contract method 0x78a2f074.
//
// Solidity: function _uTokens() view returns(address)
func (_LiquidStaking *LiquidStakingCaller) UTokens(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "_uTokens")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UTokens is a free data retrieval call binding the contract method 0x78a2f074.
//
// Solidity: function _uTokens() view returns(address)
func (_LiquidStaking *LiquidStakingSession) UTokens() (common.Address, error) {
	return _LiquidStaking.Contract.UTokens(&_LiquidStaking.CallOpts)
}

// UTokens is a free data retrieval call binding the contract method 0x78a2f074.
//
// Solidity: function _uTokens() view returns(address)
func (_LiquidStaking *LiquidStakingCallerSession) UTokens() (common.Address, error) {
	return _LiquidStaking.Contract.UTokens(&_LiquidStaking.CallOpts)
}

// UnstakingAmount is a free data retrieval call binding the contract method 0xf75f0cf6.
//
// Solidity: function _unstakingAmount(address , uint256 ) view returns(uint256)
func (_LiquidStaking *LiquidStakingCaller) UnstakingAmount(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "_unstakingAmount", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakingAmount is a free data retrieval call binding the contract method 0xf75f0cf6.
//
// Solidity: function _unstakingAmount(address , uint256 ) view returns(uint256)
func (_LiquidStaking *LiquidStakingSession) UnstakingAmount(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _LiquidStaking.Contract.UnstakingAmount(&_LiquidStaking.CallOpts, arg0, arg1)
}

// UnstakingAmount is a free data retrieval call binding the contract method 0xf75f0cf6.
//
// Solidity: function _unstakingAmount(address , uint256 ) view returns(uint256)
func (_LiquidStaking *LiquidStakingCallerSession) UnstakingAmount(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _LiquidStaking.Contract.UnstakingAmount(&_LiquidStaking.CallOpts, arg0, arg1)
}

// UnstakingExpiration is a free data retrieval call binding the contract method 0xf3e6f94e.
//
// Solidity: function _unstakingExpiration(address , uint256 ) view returns(uint256)
func (_LiquidStaking *LiquidStakingCaller) UnstakingExpiration(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "_unstakingExpiration", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakingExpiration is a free data retrieval call binding the contract method 0xf3e6f94e.
//
// Solidity: function _unstakingExpiration(address , uint256 ) view returns(uint256)
func (_LiquidStaking *LiquidStakingSession) UnstakingExpiration(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _LiquidStaking.Contract.UnstakingExpiration(&_LiquidStaking.CallOpts, arg0, arg1)
}

// UnstakingExpiration is a free data retrieval call binding the contract method 0xf3e6f94e.
//
// Solidity: function _unstakingExpiration(address , uint256 ) view returns(uint256)
func (_LiquidStaking *LiquidStakingCallerSession) UnstakingExpiration(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _LiquidStaking.Contract.UnstakingExpiration(&_LiquidStaking.CallOpts, arg0, arg1)
}

// ValueDivisor is a free data retrieval call binding the contract method 0x728a99ee.
//
// Solidity: function _valueDivisor() view returns(uint256)
func (_LiquidStaking *LiquidStakingCaller) ValueDivisor(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "_valueDivisor")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValueDivisor is a free data retrieval call binding the contract method 0x728a99ee.
//
// Solidity: function _valueDivisor() view returns(uint256)
func (_LiquidStaking *LiquidStakingSession) ValueDivisor() (*big.Int, error) {
	return _LiquidStaking.Contract.ValueDivisor(&_LiquidStaking.CallOpts)
}

// ValueDivisor is a free data retrieval call binding the contract method 0x728a99ee.
//
// Solidity: function _valueDivisor() view returns(uint256)
func (_LiquidStaking *LiquidStakingCallerSession) ValueDivisor() (*big.Int, error) {
	return _LiquidStaking.Contract.ValueDivisor(&_LiquidStaking.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x3e118dbe.
//
// Solidity: function _version() view returns(uint256)
func (_LiquidStaking *LiquidStakingCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "_version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x3e118dbe.
//
// Solidity: function _version() view returns(uint256)
func (_LiquidStaking *LiquidStakingSession) Version() (*big.Int, error) {
	return _LiquidStaking.Contract.Version(&_LiquidStaking.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x3e118dbe.
//
// Solidity: function _version() view returns(uint256)
func (_LiquidStaking *LiquidStakingCallerSession) Version() (*big.Int, error) {
	return _LiquidStaking.Contract.Version(&_LiquidStaking.CallOpts)
}

// WithdrawCounters is a free data retrieval call binding the contract method 0x898a9664.
//
// Solidity: function _withdrawCounters(address ) view returns(uint256)
func (_LiquidStaking *LiquidStakingCaller) WithdrawCounters(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "_withdrawCounters", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawCounters is a free data retrieval call binding the contract method 0x898a9664.
//
// Solidity: function _withdrawCounters(address ) view returns(uint256)
func (_LiquidStaking *LiquidStakingSession) WithdrawCounters(arg0 common.Address) (*big.Int, error) {
	return _LiquidStaking.Contract.WithdrawCounters(&_LiquidStaking.CallOpts, arg0)
}

// WithdrawCounters is a free data retrieval call binding the contract method 0x898a9664.
//
// Solidity: function _withdrawCounters(address ) view returns(uint256)
func (_LiquidStaking *LiquidStakingCallerSession) WithdrawCounters(arg0 common.Address) (*big.Int, error) {
	return _LiquidStaking.Contract.WithdrawCounters(&_LiquidStaking.CallOpts, arg0)
}

// GetBatchingLimit is a free data retrieval call binding the contract method 0xb3deb34f.
//
// Solidity: function getBatchingLimit() view returns(uint256 batchingLimit)
func (_LiquidStaking *LiquidStakingCaller) GetBatchingLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "getBatchingLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBatchingLimit is a free data retrieval call binding the contract method 0xb3deb34f.
//
// Solidity: function getBatchingLimit() view returns(uint256 batchingLimit)
func (_LiquidStaking *LiquidStakingSession) GetBatchingLimit() (*big.Int, error) {
	return _LiquidStaking.Contract.GetBatchingLimit(&_LiquidStaking.CallOpts)
}

// GetBatchingLimit is a free data retrieval call binding the contract method 0xb3deb34f.
//
// Solidity: function getBatchingLimit() view returns(uint256 batchingLimit)
func (_LiquidStaking *LiquidStakingCallerSession) GetBatchingLimit() (*big.Int, error) {
	return _LiquidStaking.Contract.GetBatchingLimit(&_LiquidStaking.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_LiquidStaking *LiquidStakingCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_LiquidStaking *LiquidStakingSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _LiquidStaking.Contract.GetRoleAdmin(&_LiquidStaking.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_LiquidStaking *LiquidStakingCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _LiquidStaking.Contract.GetRoleAdmin(&_LiquidStaking.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_LiquidStaking *LiquidStakingCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_LiquidStaking *LiquidStakingSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _LiquidStaking.Contract.GetRoleMember(&_LiquidStaking.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_LiquidStaking *LiquidStakingCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _LiquidStaking.Contract.GetRoleMember(&_LiquidStaking.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_LiquidStaking *LiquidStakingCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_LiquidStaking *LiquidStakingSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _LiquidStaking.Contract.GetRoleMemberCount(&_LiquidStaking.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_LiquidStaking *LiquidStakingCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _LiquidStaking.Contract.GetRoleMemberCount(&_LiquidStaking.CallOpts, role)
}

// GetStakeUnstakeProps is a free data retrieval call binding the contract method 0x55f2323f.
//
// Solidity: function getStakeUnstakeProps() view returns(uint256 stakeFee, uint256 unstakeFee, uint256 minStake, uint256 minUnstake, uint256 valueDivisor, uint256 epochInterval, uint256 unstakeEpoch, uint256 unstakeEpochPrevious, uint256 unstakingLockTime)
func (_LiquidStaking *LiquidStakingCaller) GetStakeUnstakeProps(opts *bind.CallOpts) (struct {
	StakeFee             *big.Int
	UnstakeFee           *big.Int
	MinStake             *big.Int
	MinUnstake           *big.Int
	ValueDivisor         *big.Int
	EpochInterval        *big.Int
	UnstakeEpoch         *big.Int
	UnstakeEpochPrevious *big.Int
	UnstakingLockTime    *big.Int
}, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "getStakeUnstakeProps")

	outstruct := new(struct {
		StakeFee             *big.Int
		UnstakeFee           *big.Int
		MinStake             *big.Int
		MinUnstake           *big.Int
		ValueDivisor         *big.Int
		EpochInterval        *big.Int
		UnstakeEpoch         *big.Int
		UnstakeEpochPrevious *big.Int
		UnstakingLockTime    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StakeFee = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.UnstakeFee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.MinStake = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.MinUnstake = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.ValueDivisor = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.EpochInterval = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.UnstakeEpoch = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.UnstakeEpochPrevious = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.UnstakingLockTime = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetStakeUnstakeProps is a free data retrieval call binding the contract method 0x55f2323f.
//
// Solidity: function getStakeUnstakeProps() view returns(uint256 stakeFee, uint256 unstakeFee, uint256 minStake, uint256 minUnstake, uint256 valueDivisor, uint256 epochInterval, uint256 unstakeEpoch, uint256 unstakeEpochPrevious, uint256 unstakingLockTime)
func (_LiquidStaking *LiquidStakingSession) GetStakeUnstakeProps() (struct {
	StakeFee             *big.Int
	UnstakeFee           *big.Int
	MinStake             *big.Int
	MinUnstake           *big.Int
	ValueDivisor         *big.Int
	EpochInterval        *big.Int
	UnstakeEpoch         *big.Int
	UnstakeEpochPrevious *big.Int
	UnstakingLockTime    *big.Int
}, error) {
	return _LiquidStaking.Contract.GetStakeUnstakeProps(&_LiquidStaking.CallOpts)
}

// GetStakeUnstakeProps is a free data retrieval call binding the contract method 0x55f2323f.
//
// Solidity: function getStakeUnstakeProps() view returns(uint256 stakeFee, uint256 unstakeFee, uint256 minStake, uint256 minUnstake, uint256 valueDivisor, uint256 epochInterval, uint256 unstakeEpoch, uint256 unstakeEpochPrevious, uint256 unstakingLockTime)
func (_LiquidStaking *LiquidStakingCallerSession) GetStakeUnstakeProps() (struct {
	StakeFee             *big.Int
	UnstakeFee           *big.Int
	MinStake             *big.Int
	MinUnstake           *big.Int
	ValueDivisor         *big.Int
	EpochInterval        *big.Int
	UnstakeEpoch         *big.Int
	UnstakeEpochPrevious *big.Int
	UnstakingLockTime    *big.Int
}, error) {
	return _LiquidStaking.Contract.GetStakeUnstakeProps(&_LiquidStaking.CallOpts)
}

// GetTokens is a free data retrieval call binding the contract method 0xaa6ca808.
//
// Solidity: function getTokens() view returns(address sTokenAddress, address uTokenAddress)
func (_LiquidStaking *LiquidStakingCaller) GetTokens(opts *bind.CallOpts) (struct {
	STokenAddress common.Address
	UTokenAddress common.Address
}, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "getTokens")

	outstruct := new(struct {
		STokenAddress common.Address
		UTokenAddress common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.STokenAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.UTokenAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// GetTokens is a free data retrieval call binding the contract method 0xaa6ca808.
//
// Solidity: function getTokens() view returns(address sTokenAddress, address uTokenAddress)
func (_LiquidStaking *LiquidStakingSession) GetTokens() (struct {
	STokenAddress common.Address
	UTokenAddress common.Address
}, error) {
	return _LiquidStaking.Contract.GetTokens(&_LiquidStaking.CallOpts)
}

// GetTokens is a free data retrieval call binding the contract method 0xaa6ca808.
//
// Solidity: function getTokens() view returns(address sTokenAddress, address uTokenAddress)
func (_LiquidStaking *LiquidStakingCallerSession) GetTokens() (struct {
	STokenAddress common.Address
	UTokenAddress common.Address
}, error) {
	return _LiquidStaking.Contract.GetTokens(&_LiquidStaking.CallOpts)
}

// GetTotalUnbondedTokens is a free data retrieval call binding the contract method 0x0273b70d.
//
// Solidity: function getTotalUnbondedTokens(address staker) view returns(uint256 unbondingTokens)
func (_LiquidStaking *LiquidStakingCaller) GetTotalUnbondedTokens(opts *bind.CallOpts, staker common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "getTotalUnbondedTokens", staker)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalUnbondedTokens is a free data retrieval call binding the contract method 0x0273b70d.
//
// Solidity: function getTotalUnbondedTokens(address staker) view returns(uint256 unbondingTokens)
func (_LiquidStaking *LiquidStakingSession) GetTotalUnbondedTokens(staker common.Address) (*big.Int, error) {
	return _LiquidStaking.Contract.GetTotalUnbondedTokens(&_LiquidStaking.CallOpts, staker)
}

// GetTotalUnbondedTokens is a free data retrieval call binding the contract method 0x0273b70d.
//
// Solidity: function getTotalUnbondedTokens(address staker) view returns(uint256 unbondingTokens)
func (_LiquidStaking *LiquidStakingCallerSession) GetTotalUnbondedTokens(staker common.Address) (*big.Int, error) {
	return _LiquidStaking.Contract.GetTotalUnbondedTokens(&_LiquidStaking.CallOpts, staker)
}

// GetTotalUnbondingTokens is a free data retrieval call binding the contract method 0xb17e3165.
//
// Solidity: function getTotalUnbondingTokens(address staker) view returns(uint256 unbondingTokens)
func (_LiquidStaking *LiquidStakingCaller) GetTotalUnbondingTokens(opts *bind.CallOpts, staker common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "getTotalUnbondingTokens", staker)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalUnbondingTokens is a free data retrieval call binding the contract method 0xb17e3165.
//
// Solidity: function getTotalUnbondingTokens(address staker) view returns(uint256 unbondingTokens)
func (_LiquidStaking *LiquidStakingSession) GetTotalUnbondingTokens(staker common.Address) (*big.Int, error) {
	return _LiquidStaking.Contract.GetTotalUnbondingTokens(&_LiquidStaking.CallOpts, staker)
}

// GetTotalUnbondingTokens is a free data retrieval call binding the contract method 0xb17e3165.
//
// Solidity: function getTotalUnbondingTokens(address staker) view returns(uint256 unbondingTokens)
func (_LiquidStaking *LiquidStakingCallerSession) GetTotalUnbondingTokens(staker common.Address) (*big.Int, error) {
	return _LiquidStaking.Contract.GetTotalUnbondingTokens(&_LiquidStaking.CallOpts, staker)
}

// GetUnstakeEpochMilestone is a free data retrieval call binding the contract method 0x9b9eea46.
//
// Solidity: function getUnstakeEpochMilestone(uint256 _unstakeTimestamp) view returns(uint256 unstakeEpochMilestone)
func (_LiquidStaking *LiquidStakingCaller) GetUnstakeEpochMilestone(opts *bind.CallOpts, _unstakeTimestamp *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "getUnstakeEpochMilestone", _unstakeTimestamp)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUnstakeEpochMilestone is a free data retrieval call binding the contract method 0x9b9eea46.
//
// Solidity: function getUnstakeEpochMilestone(uint256 _unstakeTimestamp) view returns(uint256 unstakeEpochMilestone)
func (_LiquidStaking *LiquidStakingSession) GetUnstakeEpochMilestone(_unstakeTimestamp *big.Int) (*big.Int, error) {
	return _LiquidStaking.Contract.GetUnstakeEpochMilestone(&_LiquidStaking.CallOpts, _unstakeTimestamp)
}

// GetUnstakeEpochMilestone is a free data retrieval call binding the contract method 0x9b9eea46.
//
// Solidity: function getUnstakeEpochMilestone(uint256 _unstakeTimestamp) view returns(uint256 unstakeEpochMilestone)
func (_LiquidStaking *LiquidStakingCallerSession) GetUnstakeEpochMilestone(_unstakeTimestamp *big.Int) (*big.Int, error) {
	return _LiquidStaking.Contract.GetUnstakeEpochMilestone(&_LiquidStaking.CallOpts, _unstakeTimestamp)
}

// GetUnstakeTime is a free data retrieval call binding the contract method 0x85e32840.
//
// Solidity: function getUnstakeTime(uint256 _unstakeTimestamp) view returns(uint256 unstakeTime, uint256 unstakeEpoch, uint256 unstakeEpochPrevious)
func (_LiquidStaking *LiquidStakingCaller) GetUnstakeTime(opts *bind.CallOpts, _unstakeTimestamp *big.Int) (struct {
	UnstakeTime          *big.Int
	UnstakeEpoch         *big.Int
	UnstakeEpochPrevious *big.Int
}, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "getUnstakeTime", _unstakeTimestamp)

	outstruct := new(struct {
		UnstakeTime          *big.Int
		UnstakeEpoch         *big.Int
		UnstakeEpochPrevious *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UnstakeTime = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.UnstakeEpoch = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.UnstakeEpochPrevious = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetUnstakeTime is a free data retrieval call binding the contract method 0x85e32840.
//
// Solidity: function getUnstakeTime(uint256 _unstakeTimestamp) view returns(uint256 unstakeTime, uint256 unstakeEpoch, uint256 unstakeEpochPrevious)
func (_LiquidStaking *LiquidStakingSession) GetUnstakeTime(_unstakeTimestamp *big.Int) (struct {
	UnstakeTime          *big.Int
	UnstakeEpoch         *big.Int
	UnstakeEpochPrevious *big.Int
}, error) {
	return _LiquidStaking.Contract.GetUnstakeTime(&_LiquidStaking.CallOpts, _unstakeTimestamp)
}

// GetUnstakeTime is a free data retrieval call binding the contract method 0x85e32840.
//
// Solidity: function getUnstakeTime(uint256 _unstakeTimestamp) view returns(uint256 unstakeTime, uint256 unstakeEpoch, uint256 unstakeEpochPrevious)
func (_LiquidStaking *LiquidStakingCallerSession) GetUnstakeTime(_unstakeTimestamp *big.Int) (struct {
	UnstakeTime          *big.Int
	UnstakeEpoch         *big.Int
	UnstakeEpochPrevious *big.Int
}, error) {
	return _LiquidStaking.Contract.GetUnstakeTime(&_LiquidStaking.CallOpts, _unstakeTimestamp)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_LiquidStaking *LiquidStakingCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_LiquidStaking *LiquidStakingSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _LiquidStaking.Contract.HasRole(&_LiquidStaking.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_LiquidStaking *LiquidStakingCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _LiquidStaking.Contract.HasRole(&_LiquidStaking.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_LiquidStaking *LiquidStakingCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LiquidStaking.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_LiquidStaking *LiquidStakingSession) Paused() (bool, error) {
	return _LiquidStaking.Contract.Paused(&_LiquidStaking.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_LiquidStaking *LiquidStakingCallerSession) Paused() (bool, error) {
	return _LiquidStaking.Contract.Paused(&_LiquidStaking.CallOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_LiquidStaking *LiquidStakingTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_LiquidStaking *LiquidStakingSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.GrantRole(&_LiquidStaking.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_LiquidStaking *LiquidStakingTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.GrantRole(&_LiquidStaking.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x728cdbca.
//
// Solidity: function initialize(address uAddress, address sAddress, address pauserAddress, uint256 unstakingLockTime, uint256 epochInterval, uint256 valueDivisor) returns()
func (_LiquidStaking *LiquidStakingTransactor) Initialize(opts *bind.TransactOpts, uAddress common.Address, sAddress common.Address, pauserAddress common.Address, unstakingLockTime *big.Int, epochInterval *big.Int, valueDivisor *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "initialize", uAddress, sAddress, pauserAddress, unstakingLockTime, epochInterval, valueDivisor)
}

// Initialize is a paid mutator transaction binding the contract method 0x728cdbca.
//
// Solidity: function initialize(address uAddress, address sAddress, address pauserAddress, uint256 unstakingLockTime, uint256 epochInterval, uint256 valueDivisor) returns()
func (_LiquidStaking *LiquidStakingSession) Initialize(uAddress common.Address, sAddress common.Address, pauserAddress common.Address, unstakingLockTime *big.Int, epochInterval *big.Int, valueDivisor *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.Initialize(&_LiquidStaking.TransactOpts, uAddress, sAddress, pauserAddress, unstakingLockTime, epochInterval, valueDivisor)
}

// Initialize is a paid mutator transaction binding the contract method 0x728cdbca.
//
// Solidity: function initialize(address uAddress, address sAddress, address pauserAddress, uint256 unstakingLockTime, uint256 epochInterval, uint256 valueDivisor) returns()
func (_LiquidStaking *LiquidStakingTransactorSession) Initialize(uAddress common.Address, sAddress common.Address, pauserAddress common.Address, unstakingLockTime *big.Int, epochInterval *big.Int, valueDivisor *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.Initialize(&_LiquidStaking.TransactOpts, uAddress, sAddress, pauserAddress, unstakingLockTime, epochInterval, valueDivisor)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool success)
func (_LiquidStaking *LiquidStakingTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool success)
func (_LiquidStaking *LiquidStakingSession) Pause() (*types.Transaction, error) {
	return _LiquidStaking.Contract.Pause(&_LiquidStaking.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool success)
func (_LiquidStaking *LiquidStakingTransactorSession) Pause() (*types.Transaction, error) {
	return _LiquidStaking.Contract.Pause(&_LiquidStaking.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_LiquidStaking *LiquidStakingTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_LiquidStaking *LiquidStakingSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.RenounceRole(&_LiquidStaking.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_LiquidStaking *LiquidStakingTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.RenounceRole(&_LiquidStaking.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_LiquidStaking *LiquidStakingTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_LiquidStaking *LiquidStakingSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.RevokeRole(&_LiquidStaking.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_LiquidStaking *LiquidStakingTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.RevokeRole(&_LiquidStaking.TransactOpts, role, account)
}

// SetBatchingLimit is a paid mutator transaction binding the contract method 0xbef1307c.
//
// Solidity: function setBatchingLimit(uint256 batchingLimit) returns(bool success)
func (_LiquidStaking *LiquidStakingTransactor) SetBatchingLimit(opts *bind.TransactOpts, batchingLimit *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "setBatchingLimit", batchingLimit)
}

// SetBatchingLimit is a paid mutator transaction binding the contract method 0xbef1307c.
//
// Solidity: function setBatchingLimit(uint256 batchingLimit) returns(bool success)
func (_LiquidStaking *LiquidStakingSession) SetBatchingLimit(batchingLimit *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetBatchingLimit(&_LiquidStaking.TransactOpts, batchingLimit)
}

// SetBatchingLimit is a paid mutator transaction binding the contract method 0xbef1307c.
//
// Solidity: function setBatchingLimit(uint256 batchingLimit) returns(bool success)
func (_LiquidStaking *LiquidStakingTransactorSession) SetBatchingLimit(batchingLimit *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetBatchingLimit(&_LiquidStaking.TransactOpts, batchingLimit)
}

// SetFees is a paid mutator transaction binding the contract method 0x0b78f9c0.
//
// Solidity: function setFees(uint256 stakeFee, uint256 unstakeFee) returns(bool success)
func (_LiquidStaking *LiquidStakingTransactor) SetFees(opts *bind.TransactOpts, stakeFee *big.Int, unstakeFee *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "setFees", stakeFee, unstakeFee)
}

// SetFees is a paid mutator transaction binding the contract method 0x0b78f9c0.
//
// Solidity: function setFees(uint256 stakeFee, uint256 unstakeFee) returns(bool success)
func (_LiquidStaking *LiquidStakingSession) SetFees(stakeFee *big.Int, unstakeFee *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetFees(&_LiquidStaking.TransactOpts, stakeFee, unstakeFee)
}

// SetFees is a paid mutator transaction binding the contract method 0x0b78f9c0.
//
// Solidity: function setFees(uint256 stakeFee, uint256 unstakeFee) returns(bool success)
func (_LiquidStaking *LiquidStakingTransactorSession) SetFees(stakeFee *big.Int, unstakeFee *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetFees(&_LiquidStaking.TransactOpts, stakeFee, unstakeFee)
}

// SetMinimumValues is a paid mutator transaction binding the contract method 0x51db00e7.
//
// Solidity: function setMinimumValues(uint256 minStake, uint256 minUnstake) returns(bool success)
func (_LiquidStaking *LiquidStakingTransactor) SetMinimumValues(opts *bind.TransactOpts, minStake *big.Int, minUnstake *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "setMinimumValues", minStake, minUnstake)
}

// SetMinimumValues is a paid mutator transaction binding the contract method 0x51db00e7.
//
// Solidity: function setMinimumValues(uint256 minStake, uint256 minUnstake) returns(bool success)
func (_LiquidStaking *LiquidStakingSession) SetMinimumValues(minStake *big.Int, minUnstake *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetMinimumValues(&_LiquidStaking.TransactOpts, minStake, minUnstake)
}

// SetMinimumValues is a paid mutator transaction binding the contract method 0x51db00e7.
//
// Solidity: function setMinimumValues(uint256 minStake, uint256 minUnstake) returns(bool success)
func (_LiquidStaking *LiquidStakingTransactorSession) SetMinimumValues(minStake *big.Int, minUnstake *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetMinimumValues(&_LiquidStaking.TransactOpts, minStake, minUnstake)
}

// SetSTokensContract is a paid mutator transaction binding the contract method 0x2e31fab8.
//
// Solidity: function setSTokensContract(address sAddress) returns()
func (_LiquidStaking *LiquidStakingTransactor) SetSTokensContract(opts *bind.TransactOpts, sAddress common.Address) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "setSTokensContract", sAddress)
}

// SetSTokensContract is a paid mutator transaction binding the contract method 0x2e31fab8.
//
// Solidity: function setSTokensContract(address sAddress) returns()
func (_LiquidStaking *LiquidStakingSession) SetSTokensContract(sAddress common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetSTokensContract(&_LiquidStaking.TransactOpts, sAddress)
}

// SetSTokensContract is a paid mutator transaction binding the contract method 0x2e31fab8.
//
// Solidity: function setSTokensContract(address sAddress) returns()
func (_LiquidStaking *LiquidStakingTransactorSession) SetSTokensContract(sAddress common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetSTokensContract(&_LiquidStaking.TransactOpts, sAddress)
}

// SetTokenWrapperContract is a paid mutator transaction binding the contract method 0x31b188cc.
//
// Solidity: function setTokenWrapperContract(address tokenWrapperContract) returns()
func (_LiquidStaking *LiquidStakingTransactor) SetTokenWrapperContract(opts *bind.TransactOpts, tokenWrapperContract common.Address) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "setTokenWrapperContract", tokenWrapperContract)
}

// SetTokenWrapperContract is a paid mutator transaction binding the contract method 0x31b188cc.
//
// Solidity: function setTokenWrapperContract(address tokenWrapperContract) returns()
func (_LiquidStaking *LiquidStakingSession) SetTokenWrapperContract(tokenWrapperContract common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetTokenWrapperContract(&_LiquidStaking.TransactOpts, tokenWrapperContract)
}

// SetTokenWrapperContract is a paid mutator transaction binding the contract method 0x31b188cc.
//
// Solidity: function setTokenWrapperContract(address tokenWrapperContract) returns()
func (_LiquidStaking *LiquidStakingTransactorSession) SetTokenWrapperContract(tokenWrapperContract common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetTokenWrapperContract(&_LiquidStaking.TransactOpts, tokenWrapperContract)
}

// SetUTokensContract is a paid mutator transaction binding the contract method 0x21bdf9e5.
//
// Solidity: function setUTokensContract(address uAddress) returns()
func (_LiquidStaking *LiquidStakingTransactor) SetUTokensContract(opts *bind.TransactOpts, uAddress common.Address) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "setUTokensContract", uAddress)
}

// SetUTokensContract is a paid mutator transaction binding the contract method 0x21bdf9e5.
//
// Solidity: function setUTokensContract(address uAddress) returns()
func (_LiquidStaking *LiquidStakingSession) SetUTokensContract(uAddress common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetUTokensContract(&_LiquidStaking.TransactOpts, uAddress)
}

// SetUTokensContract is a paid mutator transaction binding the contract method 0x21bdf9e5.
//
// Solidity: function setUTokensContract(address uAddress) returns()
func (_LiquidStaking *LiquidStakingTransactorSession) SetUTokensContract(uAddress common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetUTokensContract(&_LiquidStaking.TransactOpts, uAddress)
}

// SetUnstakeEpoch is a paid mutator transaction binding the contract method 0x9dc77cab.
//
// Solidity: function setUnstakeEpoch(uint256 unstakeEpoch, uint256 unstakeEpochPrevious, uint256 epochInterval) returns(bool success)
func (_LiquidStaking *LiquidStakingTransactor) SetUnstakeEpoch(opts *bind.TransactOpts, unstakeEpoch *big.Int, unstakeEpochPrevious *big.Int, epochInterval *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "setUnstakeEpoch", unstakeEpoch, unstakeEpochPrevious, epochInterval)
}

// SetUnstakeEpoch is a paid mutator transaction binding the contract method 0x9dc77cab.
//
// Solidity: function setUnstakeEpoch(uint256 unstakeEpoch, uint256 unstakeEpochPrevious, uint256 epochInterval) returns(bool success)
func (_LiquidStaking *LiquidStakingSession) SetUnstakeEpoch(unstakeEpoch *big.Int, unstakeEpochPrevious *big.Int, epochInterval *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetUnstakeEpoch(&_LiquidStaking.TransactOpts, unstakeEpoch, unstakeEpochPrevious, epochInterval)
}

// SetUnstakeEpoch is a paid mutator transaction binding the contract method 0x9dc77cab.
//
// Solidity: function setUnstakeEpoch(uint256 unstakeEpoch, uint256 unstakeEpochPrevious, uint256 epochInterval) returns(bool success)
func (_LiquidStaking *LiquidStakingTransactorSession) SetUnstakeEpoch(unstakeEpoch *big.Int, unstakeEpochPrevious *big.Int, epochInterval *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetUnstakeEpoch(&_LiquidStaking.TransactOpts, unstakeEpoch, unstakeEpochPrevious, epochInterval)
}

// SetUnstakingLockTime is a paid mutator transaction binding the contract method 0x0e4b2943.
//
// Solidity: function setUnstakingLockTime(uint256 unstakingLockTime) returns(bool success)
func (_LiquidStaking *LiquidStakingTransactor) SetUnstakingLockTime(opts *bind.TransactOpts, unstakingLockTime *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "setUnstakingLockTime", unstakingLockTime)
}

// SetUnstakingLockTime is a paid mutator transaction binding the contract method 0x0e4b2943.
//
// Solidity: function setUnstakingLockTime(uint256 unstakingLockTime) returns(bool success)
func (_LiquidStaking *LiquidStakingSession) SetUnstakingLockTime(unstakingLockTime *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetUnstakingLockTime(&_LiquidStaking.TransactOpts, unstakingLockTime)
}

// SetUnstakingLockTime is a paid mutator transaction binding the contract method 0x0e4b2943.
//
// Solidity: function setUnstakingLockTime(uint256 unstakingLockTime) returns(bool success)
func (_LiquidStaking *LiquidStakingTransactorSession) SetUnstakingLockTime(unstakingLockTime *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.SetUnstakingLockTime(&_LiquidStaking.TransactOpts, unstakingLockTime)
}

// Stake is a paid mutator transaction binding the contract method 0xadc9772e.
//
// Solidity: function stake(address to, uint256 amount) returns(uint256 amountStaked)
func (_LiquidStaking *LiquidStakingTransactor) Stake(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "stake", to, amount)
}

// Stake is a paid mutator transaction binding the contract method 0xadc9772e.
//
// Solidity: function stake(address to, uint256 amount) returns(uint256 amountStaked)
func (_LiquidStaking *LiquidStakingSession) Stake(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.Stake(&_LiquidStaking.TransactOpts, to, amount)
}

// Stake is a paid mutator transaction binding the contract method 0xadc9772e.
//
// Solidity: function stake(address to, uint256 amount) returns(uint256 amountStaked)
func (_LiquidStaking *LiquidStakingTransactorSession) Stake(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.Stake(&_LiquidStaking.TransactOpts, to, amount)
}

// StakeDirect is a paid mutator transaction binding the contract method 0xe321fca6.
//
// Solidity: function stakeDirect(address to, uint256 stakingAmount, uint256 wrappingAmount) returns(uint256 amountStaked, uint256 amountWrapped)
func (_LiquidStaking *LiquidStakingTransactor) StakeDirect(opts *bind.TransactOpts, to common.Address, stakingAmount *big.Int, wrappingAmount *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "stakeDirect", to, stakingAmount, wrappingAmount)
}

// StakeDirect is a paid mutator transaction binding the contract method 0xe321fca6.
//
// Solidity: function stakeDirect(address to, uint256 stakingAmount, uint256 wrappingAmount) returns(uint256 amountStaked, uint256 amountWrapped)
func (_LiquidStaking *LiquidStakingSession) StakeDirect(to common.Address, stakingAmount *big.Int, wrappingAmount *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.StakeDirect(&_LiquidStaking.TransactOpts, to, stakingAmount, wrappingAmount)
}

// StakeDirect is a paid mutator transaction binding the contract method 0xe321fca6.
//
// Solidity: function stakeDirect(address to, uint256 stakingAmount, uint256 wrappingAmount) returns(uint256 amountStaked, uint256 amountWrapped)
func (_LiquidStaking *LiquidStakingTransactorSession) StakeDirect(to common.Address, stakingAmount *big.Int, wrappingAmount *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.StakeDirect(&_LiquidStaking.TransactOpts, to, stakingAmount, wrappingAmount)
}

// StakeDirectInBatch is a paid mutator transaction binding the contract method 0x10ef62f6.
//
// Solidity: function stakeDirectInBatch(address[] toAddressses, uint256[] stakingAmounts, uint256[] wrappingAmounts) returns(uint256[] amountsStaked, uint256[] amountsWrapped)
func (_LiquidStaking *LiquidStakingTransactor) StakeDirectInBatch(opts *bind.TransactOpts, toAddressses []common.Address, stakingAmounts []*big.Int, wrappingAmounts []*big.Int) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "stakeDirectInBatch", toAddressses, stakingAmounts, wrappingAmounts)
}

// StakeDirectInBatch is a paid mutator transaction binding the contract method 0x10ef62f6.
//
// Solidity: function stakeDirectInBatch(address[] toAddressses, uint256[] stakingAmounts, uint256[] wrappingAmounts) returns(uint256[] amountsStaked, uint256[] amountsWrapped)
func (_LiquidStaking *LiquidStakingSession) StakeDirectInBatch(toAddressses []common.Address, stakingAmounts []*big.Int, wrappingAmounts []*big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.StakeDirectInBatch(&_LiquidStaking.TransactOpts, toAddressses, stakingAmounts, wrappingAmounts)
}

// StakeDirectInBatch is a paid mutator transaction binding the contract method 0x10ef62f6.
//
// Solidity: function stakeDirectInBatch(address[] toAddressses, uint256[] stakingAmounts, uint256[] wrappingAmounts) returns(uint256[] amountsStaked, uint256[] amountsWrapped)
func (_LiquidStaking *LiquidStakingTransactorSession) StakeDirectInBatch(toAddressses []common.Address, stakingAmounts []*big.Int, wrappingAmounts []*big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.StakeDirectInBatch(&_LiquidStaking.TransactOpts, toAddressses, stakingAmounts, wrappingAmounts)
}

// UnStake is a paid mutator transaction binding the contract method 0xd9393814.
//
// Solidity: function unStake(address to, uint256 amount) returns(uint256 amountUnstaked)
func (_LiquidStaking *LiquidStakingTransactor) UnStake(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "unStake", to, amount)
}

// UnStake is a paid mutator transaction binding the contract method 0xd9393814.
//
// Solidity: function unStake(address to, uint256 amount) returns(uint256 amountUnstaked)
func (_LiquidStaking *LiquidStakingSession) UnStake(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.UnStake(&_LiquidStaking.TransactOpts, to, amount)
}

// UnStake is a paid mutator transaction binding the contract method 0xd9393814.
//
// Solidity: function unStake(address to, uint256 amount) returns(uint256 amountUnstaked)
func (_LiquidStaking *LiquidStakingTransactorSession) UnStake(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LiquidStaking.Contract.UnStake(&_LiquidStaking.TransactOpts, to, amount)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool success)
func (_LiquidStaking *LiquidStakingTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool success)
func (_LiquidStaking *LiquidStakingSession) Unpause() (*types.Transaction, error) {
	return _LiquidStaking.Contract.Unpause(&_LiquidStaking.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool success)
func (_LiquidStaking *LiquidStakingTransactorSession) Unpause() (*types.Transaction, error) {
	return _LiquidStaking.Contract.Unpause(&_LiquidStaking.TransactOpts)
}

// WithdrawUnstakedTokens is a paid mutator transaction binding the contract method 0x819020a0.
//
// Solidity: function withdrawUnstakedTokens(address staker) returns(uint256 withdrawBalance)
func (_LiquidStaking *LiquidStakingTransactor) WithdrawUnstakedTokens(opts *bind.TransactOpts, staker common.Address) (*types.Transaction, error) {
	return _LiquidStaking.contract.Transact(opts, "withdrawUnstakedTokens", staker)
}

// WithdrawUnstakedTokens is a paid mutator transaction binding the contract method 0x819020a0.
//
// Solidity: function withdrawUnstakedTokens(address staker) returns(uint256 withdrawBalance)
func (_LiquidStaking *LiquidStakingSession) WithdrawUnstakedTokens(staker common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.WithdrawUnstakedTokens(&_LiquidStaking.TransactOpts, staker)
}

// WithdrawUnstakedTokens is a paid mutator transaction binding the contract method 0x819020a0.
//
// Solidity: function withdrawUnstakedTokens(address staker) returns(uint256 withdrawBalance)
func (_LiquidStaking *LiquidStakingTransactorSession) WithdrawUnstakedTokens(staker common.Address) (*types.Transaction, error) {
	return _LiquidStaking.Contract.WithdrawUnstakedTokens(&_LiquidStaking.TransactOpts, staker)
}

// LiquidStakingPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the LiquidStaking contract.
type LiquidStakingPausedIterator struct {
	Event *LiquidStakingPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingPaused represents a Paused event raised by the LiquidStaking contract.
type LiquidStakingPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_LiquidStaking *LiquidStakingFilterer) FilterPaused(opts *bind.FilterOpts) (*LiquidStakingPausedIterator, error) {

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &LiquidStakingPausedIterator{contract: _LiquidStaking.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_LiquidStaking *LiquidStakingFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *LiquidStakingPaused) (event.Subscription, error) {

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingPaused)
				if err := _LiquidStaking.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_LiquidStaking *LiquidStakingFilterer) ParsePaused(log types.Log) (*LiquidStakingPaused, error) {
	event := new(LiquidStakingPaused)
	if err := _LiquidStaking.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the LiquidStaking contract.
type LiquidStakingRoleAdminChangedIterator struct {
	Event *LiquidStakingRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingRoleAdminChanged represents a RoleAdminChanged event raised by the LiquidStaking contract.
type LiquidStakingRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_LiquidStaking *LiquidStakingFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*LiquidStakingRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingRoleAdminChangedIterator{contract: _LiquidStaking.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_LiquidStaking *LiquidStakingFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *LiquidStakingRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingRoleAdminChanged)
				if err := _LiquidStaking.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_LiquidStaking *LiquidStakingFilterer) ParseRoleAdminChanged(log types.Log) (*LiquidStakingRoleAdminChanged, error) {
	event := new(LiquidStakingRoleAdminChanged)
	if err := _LiquidStaking.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the LiquidStaking contract.
type LiquidStakingRoleGrantedIterator struct {
	Event *LiquidStakingRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingRoleGranted represents a RoleGranted event raised by the LiquidStaking contract.
type LiquidStakingRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_LiquidStaking *LiquidStakingFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*LiquidStakingRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingRoleGrantedIterator{contract: _LiquidStaking.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_LiquidStaking *LiquidStakingFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *LiquidStakingRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingRoleGranted)
				if err := _LiquidStaking.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_LiquidStaking *LiquidStakingFilterer) ParseRoleGranted(log types.Log) (*LiquidStakingRoleGranted, error) {
	event := new(LiquidStakingRoleGranted)
	if err := _LiquidStaking.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the LiquidStaking contract.
type LiquidStakingRoleRevokedIterator struct {
	Event *LiquidStakingRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingRoleRevoked represents a RoleRevoked event raised by the LiquidStaking contract.
type LiquidStakingRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_LiquidStaking *LiquidStakingFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*LiquidStakingRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingRoleRevokedIterator{contract: _LiquidStaking.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_LiquidStaking *LiquidStakingFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *LiquidStakingRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingRoleRevoked)
				if err := _LiquidStaking.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_LiquidStaking *LiquidStakingFilterer) ParseRoleRevoked(log types.Log) (*LiquidStakingRoleRevoked, error) {
	event := new(LiquidStakingRoleRevoked)
	if err := _LiquidStaking.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingSetBatchingLimitIterator is returned from FilterSetBatchingLimit and is used to iterate over the raw logs and unpacked data for SetBatchingLimit events raised by the LiquidStaking contract.
type LiquidStakingSetBatchingLimitIterator struct {
	Event *LiquidStakingSetBatchingLimit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingSetBatchingLimitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingSetBatchingLimit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingSetBatchingLimit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingSetBatchingLimitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingSetBatchingLimitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingSetBatchingLimit represents a SetBatchingLimit event raised by the LiquidStaking contract.
type LiquidStakingSetBatchingLimit struct {
	BatchingLimit *big.Int
	Timestamp     *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSetBatchingLimit is a free log retrieval operation binding the contract event 0x3e0887a905af0c5764f2c25c201fa1b9595065cb7008f54fe7734fefdd8e88dc.
//
// Solidity: event SetBatchingLimit(uint256 indexed batchingLimit, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) FilterSetBatchingLimit(opts *bind.FilterOpts, batchingLimit []*big.Int) (*LiquidStakingSetBatchingLimitIterator, error) {

	var batchingLimitRule []interface{}
	for _, batchingLimitItem := range batchingLimit {
		batchingLimitRule = append(batchingLimitRule, batchingLimitItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "SetBatchingLimit", batchingLimitRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingSetBatchingLimitIterator{contract: _LiquidStaking.contract, event: "SetBatchingLimit", logs: logs, sub: sub}, nil
}

// WatchSetBatchingLimit is a free log subscription operation binding the contract event 0x3e0887a905af0c5764f2c25c201fa1b9595065cb7008f54fe7734fefdd8e88dc.
//
// Solidity: event SetBatchingLimit(uint256 indexed batchingLimit, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) WatchSetBatchingLimit(opts *bind.WatchOpts, sink chan<- *LiquidStakingSetBatchingLimit, batchingLimit []*big.Int) (event.Subscription, error) {

	var batchingLimitRule []interface{}
	for _, batchingLimitItem := range batchingLimit {
		batchingLimitRule = append(batchingLimitRule, batchingLimitItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "SetBatchingLimit", batchingLimitRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingSetBatchingLimit)
				if err := _LiquidStaking.contract.UnpackLog(event, "SetBatchingLimit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetBatchingLimit is a log parse operation binding the contract event 0x3e0887a905af0c5764f2c25c201fa1b9595065cb7008f54fe7734fefdd8e88dc.
//
// Solidity: event SetBatchingLimit(uint256 indexed batchingLimit, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) ParseSetBatchingLimit(log types.Log) (*LiquidStakingSetBatchingLimit, error) {
	event := new(LiquidStakingSetBatchingLimit)
	if err := _LiquidStaking.contract.UnpackLog(event, "SetBatchingLimit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingSetFeesIterator is returned from FilterSetFees and is used to iterate over the raw logs and unpacked data for SetFees events raised by the LiquidStaking contract.
type LiquidStakingSetFeesIterator struct {
	Event *LiquidStakingSetFees // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingSetFeesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingSetFees)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingSetFees)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingSetFeesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingSetFeesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingSetFees represents a SetFees event raised by the LiquidStaking contract.
type LiquidStakingSetFees struct {
	StakeFee   *big.Int
	UnstakeFee *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSetFees is a free log retrieval operation binding the contract event 0x3be6e637f54f081a2dd12982cfd58481304fc98b604b18fa8ed4aa4e4a9d1532.
//
// Solidity: event SetFees(uint256 indexed stakeFee, uint256 indexed unstakeFee)
func (_LiquidStaking *LiquidStakingFilterer) FilterSetFees(opts *bind.FilterOpts, stakeFee []*big.Int, unstakeFee []*big.Int) (*LiquidStakingSetFeesIterator, error) {

	var stakeFeeRule []interface{}
	for _, stakeFeeItem := range stakeFee {
		stakeFeeRule = append(stakeFeeRule, stakeFeeItem)
	}
	var unstakeFeeRule []interface{}
	for _, unstakeFeeItem := range unstakeFee {
		unstakeFeeRule = append(unstakeFeeRule, unstakeFeeItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "SetFees", stakeFeeRule, unstakeFeeRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingSetFeesIterator{contract: _LiquidStaking.contract, event: "SetFees", logs: logs, sub: sub}, nil
}

// WatchSetFees is a free log subscription operation binding the contract event 0x3be6e637f54f081a2dd12982cfd58481304fc98b604b18fa8ed4aa4e4a9d1532.
//
// Solidity: event SetFees(uint256 indexed stakeFee, uint256 indexed unstakeFee)
func (_LiquidStaking *LiquidStakingFilterer) WatchSetFees(opts *bind.WatchOpts, sink chan<- *LiquidStakingSetFees, stakeFee []*big.Int, unstakeFee []*big.Int) (event.Subscription, error) {

	var stakeFeeRule []interface{}
	for _, stakeFeeItem := range stakeFee {
		stakeFeeRule = append(stakeFeeRule, stakeFeeItem)
	}
	var unstakeFeeRule []interface{}
	for _, unstakeFeeItem := range unstakeFee {
		unstakeFeeRule = append(unstakeFeeRule, unstakeFeeItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "SetFees", stakeFeeRule, unstakeFeeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingSetFees)
				if err := _LiquidStaking.contract.UnpackLog(event, "SetFees", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetFees is a log parse operation binding the contract event 0x3be6e637f54f081a2dd12982cfd58481304fc98b604b18fa8ed4aa4e4a9d1532.
//
// Solidity: event SetFees(uint256 indexed stakeFee, uint256 indexed unstakeFee)
func (_LiquidStaking *LiquidStakingFilterer) ParseSetFees(log types.Log) (*LiquidStakingSetFees, error) {
	event := new(LiquidStakingSetFees)
	if err := _LiquidStaking.contract.UnpackLog(event, "SetFees", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingSetMinimumValuesIterator is returned from FilterSetMinimumValues and is used to iterate over the raw logs and unpacked data for SetMinimumValues events raised by the LiquidStaking contract.
type LiquidStakingSetMinimumValuesIterator struct {
	Event *LiquidStakingSetMinimumValues // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingSetMinimumValuesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingSetMinimumValues)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingSetMinimumValues)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingSetMinimumValuesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingSetMinimumValuesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingSetMinimumValues represents a SetMinimumValues event raised by the LiquidStaking contract.
type LiquidStakingSetMinimumValues struct {
	MinStake   *big.Int
	MinUnstake *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSetMinimumValues is a free log retrieval operation binding the contract event 0xa94d40347449ea19b873d64ed55c99772194c49791e82be5c2dc24e4c0044fe3.
//
// Solidity: event SetMinimumValues(uint256 indexed minStake, uint256 indexed minUnstake)
func (_LiquidStaking *LiquidStakingFilterer) FilterSetMinimumValues(opts *bind.FilterOpts, minStake []*big.Int, minUnstake []*big.Int) (*LiquidStakingSetMinimumValuesIterator, error) {

	var minStakeRule []interface{}
	for _, minStakeItem := range minStake {
		minStakeRule = append(minStakeRule, minStakeItem)
	}
	var minUnstakeRule []interface{}
	for _, minUnstakeItem := range minUnstake {
		minUnstakeRule = append(minUnstakeRule, minUnstakeItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "SetMinimumValues", minStakeRule, minUnstakeRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingSetMinimumValuesIterator{contract: _LiquidStaking.contract, event: "SetMinimumValues", logs: logs, sub: sub}, nil
}

// WatchSetMinimumValues is a free log subscription operation binding the contract event 0xa94d40347449ea19b873d64ed55c99772194c49791e82be5c2dc24e4c0044fe3.
//
// Solidity: event SetMinimumValues(uint256 indexed minStake, uint256 indexed minUnstake)
func (_LiquidStaking *LiquidStakingFilterer) WatchSetMinimumValues(opts *bind.WatchOpts, sink chan<- *LiquidStakingSetMinimumValues, minStake []*big.Int, minUnstake []*big.Int) (event.Subscription, error) {

	var minStakeRule []interface{}
	for _, minStakeItem := range minStake {
		minStakeRule = append(minStakeRule, minStakeItem)
	}
	var minUnstakeRule []interface{}
	for _, minUnstakeItem := range minUnstake {
		minUnstakeRule = append(minUnstakeRule, minUnstakeItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "SetMinimumValues", minStakeRule, minUnstakeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingSetMinimumValues)
				if err := _LiquidStaking.contract.UnpackLog(event, "SetMinimumValues", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetMinimumValues is a log parse operation binding the contract event 0xa94d40347449ea19b873d64ed55c99772194c49791e82be5c2dc24e4c0044fe3.
//
// Solidity: event SetMinimumValues(uint256 indexed minStake, uint256 indexed minUnstake)
func (_LiquidStaking *LiquidStakingFilterer) ParseSetMinimumValues(log types.Log) (*LiquidStakingSetMinimumValues, error) {
	event := new(LiquidStakingSetMinimumValues)
	if err := _LiquidStaking.contract.UnpackLog(event, "SetMinimumValues", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingSetSTokensContractIterator is returned from FilterSetSTokensContract and is used to iterate over the raw logs and unpacked data for SetSTokensContract events raised by the LiquidStaking contract.
type LiquidStakingSetSTokensContractIterator struct {
	Event *LiquidStakingSetSTokensContract // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingSetSTokensContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingSetSTokensContract)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingSetSTokensContract)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingSetSTokensContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingSetSTokensContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingSetSTokensContract represents a SetSTokensContract event raised by the LiquidStaking contract.
type LiquidStakingSetSTokensContract struct {
	Contract common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetSTokensContract is a free log retrieval operation binding the contract event 0xd0b9dbe26f676f4cdb02fc249c203baa0314cec79e18b4194dc862795a5de03f.
//
// Solidity: event SetSTokensContract(address indexed _contract)
func (_LiquidStaking *LiquidStakingFilterer) FilterSetSTokensContract(opts *bind.FilterOpts, _contract []common.Address) (*LiquidStakingSetSTokensContractIterator, error) {

	var _contractRule []interface{}
	for _, _contractItem := range _contract {
		_contractRule = append(_contractRule, _contractItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "SetSTokensContract", _contractRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingSetSTokensContractIterator{contract: _LiquidStaking.contract, event: "SetSTokensContract", logs: logs, sub: sub}, nil
}

// WatchSetSTokensContract is a free log subscription operation binding the contract event 0xd0b9dbe26f676f4cdb02fc249c203baa0314cec79e18b4194dc862795a5de03f.
//
// Solidity: event SetSTokensContract(address indexed _contract)
func (_LiquidStaking *LiquidStakingFilterer) WatchSetSTokensContract(opts *bind.WatchOpts, sink chan<- *LiquidStakingSetSTokensContract, _contract []common.Address) (event.Subscription, error) {

	var _contractRule []interface{}
	for _, _contractItem := range _contract {
		_contractRule = append(_contractRule, _contractItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "SetSTokensContract", _contractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingSetSTokensContract)
				if err := _LiquidStaking.contract.UnpackLog(event, "SetSTokensContract", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetSTokensContract is a log parse operation binding the contract event 0xd0b9dbe26f676f4cdb02fc249c203baa0314cec79e18b4194dc862795a5de03f.
//
// Solidity: event SetSTokensContract(address indexed _contract)
func (_LiquidStaking *LiquidStakingFilterer) ParseSetSTokensContract(log types.Log) (*LiquidStakingSetSTokensContract, error) {
	event := new(LiquidStakingSetSTokensContract)
	if err := _LiquidStaking.contract.UnpackLog(event, "SetSTokensContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingSetTokenWrapperContractIterator is returned from FilterSetTokenWrapperContract and is used to iterate over the raw logs and unpacked data for SetTokenWrapperContract events raised by the LiquidStaking contract.
type LiquidStakingSetTokenWrapperContractIterator struct {
	Event *LiquidStakingSetTokenWrapperContract // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingSetTokenWrapperContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingSetTokenWrapperContract)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingSetTokenWrapperContract)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingSetTokenWrapperContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingSetTokenWrapperContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingSetTokenWrapperContract represents a SetTokenWrapperContract event raised by the LiquidStaking contract.
type LiquidStakingSetTokenWrapperContract struct {
	Contract common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetTokenWrapperContract is a free log retrieval operation binding the contract event 0x2f9f272420570ea07d7da789fdc59a50c36792d6ea4f49a4d218b6f2431d0492.
//
// Solidity: event SetTokenWrapperContract(address indexed _contract)
func (_LiquidStaking *LiquidStakingFilterer) FilterSetTokenWrapperContract(opts *bind.FilterOpts, _contract []common.Address) (*LiquidStakingSetTokenWrapperContractIterator, error) {

	var _contractRule []interface{}
	for _, _contractItem := range _contract {
		_contractRule = append(_contractRule, _contractItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "SetTokenWrapperContract", _contractRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingSetTokenWrapperContractIterator{contract: _LiquidStaking.contract, event: "SetTokenWrapperContract", logs: logs, sub: sub}, nil
}

// WatchSetTokenWrapperContract is a free log subscription operation binding the contract event 0x2f9f272420570ea07d7da789fdc59a50c36792d6ea4f49a4d218b6f2431d0492.
//
// Solidity: event SetTokenWrapperContract(address indexed _contract)
func (_LiquidStaking *LiquidStakingFilterer) WatchSetTokenWrapperContract(opts *bind.WatchOpts, sink chan<- *LiquidStakingSetTokenWrapperContract, _contract []common.Address) (event.Subscription, error) {

	var _contractRule []interface{}
	for _, _contractItem := range _contract {
		_contractRule = append(_contractRule, _contractItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "SetTokenWrapperContract", _contractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingSetTokenWrapperContract)
				if err := _LiquidStaking.contract.UnpackLog(event, "SetTokenWrapperContract", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetTokenWrapperContract is a log parse operation binding the contract event 0x2f9f272420570ea07d7da789fdc59a50c36792d6ea4f49a4d218b6f2431d0492.
//
// Solidity: event SetTokenWrapperContract(address indexed _contract)
func (_LiquidStaking *LiquidStakingFilterer) ParseSetTokenWrapperContract(log types.Log) (*LiquidStakingSetTokenWrapperContract, error) {
	event := new(LiquidStakingSetTokenWrapperContract)
	if err := _LiquidStaking.contract.UnpackLog(event, "SetTokenWrapperContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingSetUTokensContractIterator is returned from FilterSetUTokensContract and is used to iterate over the raw logs and unpacked data for SetUTokensContract events raised by the LiquidStaking contract.
type LiquidStakingSetUTokensContractIterator struct {
	Event *LiquidStakingSetUTokensContract // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingSetUTokensContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingSetUTokensContract)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingSetUTokensContract)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingSetUTokensContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingSetUTokensContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingSetUTokensContract represents a SetUTokensContract event raised by the LiquidStaking contract.
type LiquidStakingSetUTokensContract struct {
	Contract common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetUTokensContract is a free log retrieval operation binding the contract event 0xf745c285f36f88bea6af14d4d0f33dd9350cef7895216c1615b1caaee7857e0c.
//
// Solidity: event SetUTokensContract(address indexed _contract)
func (_LiquidStaking *LiquidStakingFilterer) FilterSetUTokensContract(opts *bind.FilterOpts, _contract []common.Address) (*LiquidStakingSetUTokensContractIterator, error) {

	var _contractRule []interface{}
	for _, _contractItem := range _contract {
		_contractRule = append(_contractRule, _contractItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "SetUTokensContract", _contractRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingSetUTokensContractIterator{contract: _LiquidStaking.contract, event: "SetUTokensContract", logs: logs, sub: sub}, nil
}

// WatchSetUTokensContract is a free log subscription operation binding the contract event 0xf745c285f36f88bea6af14d4d0f33dd9350cef7895216c1615b1caaee7857e0c.
//
// Solidity: event SetUTokensContract(address indexed _contract)
func (_LiquidStaking *LiquidStakingFilterer) WatchSetUTokensContract(opts *bind.WatchOpts, sink chan<- *LiquidStakingSetUTokensContract, _contract []common.Address) (event.Subscription, error) {

	var _contractRule []interface{}
	for _, _contractItem := range _contract {
		_contractRule = append(_contractRule, _contractItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "SetUTokensContract", _contractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingSetUTokensContract)
				if err := _LiquidStaking.contract.UnpackLog(event, "SetUTokensContract", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetUTokensContract is a log parse operation binding the contract event 0xf745c285f36f88bea6af14d4d0f33dd9350cef7895216c1615b1caaee7857e0c.
//
// Solidity: event SetUTokensContract(address indexed _contract)
func (_LiquidStaking *LiquidStakingFilterer) ParseSetUTokensContract(log types.Log) (*LiquidStakingSetUTokensContract, error) {
	event := new(LiquidStakingSetUTokensContract)
	if err := _LiquidStaking.contract.UnpackLog(event, "SetUTokensContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingSetUnstakeEpochIterator is returned from FilterSetUnstakeEpoch and is used to iterate over the raw logs and unpacked data for SetUnstakeEpoch events raised by the LiquidStaking contract.
type LiquidStakingSetUnstakeEpochIterator struct {
	Event *LiquidStakingSetUnstakeEpoch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingSetUnstakeEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingSetUnstakeEpoch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingSetUnstakeEpoch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingSetUnstakeEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingSetUnstakeEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingSetUnstakeEpoch represents a SetUnstakeEpoch event raised by the LiquidStaking contract.
type LiquidStakingSetUnstakeEpoch struct {
	UnstakeEpoch         *big.Int
	UnstakeEpochPrevious *big.Int
	EpochInterval        *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetUnstakeEpoch is a free log retrieval operation binding the contract event 0x342db083975cc0697d0f81d0876299ecbf95d909e920426089062e540904497c.
//
// Solidity: event SetUnstakeEpoch(uint256 indexed unstakeEpoch, uint256 indexed unstakeEpochPrevious, uint256 indexed epochInterval)
func (_LiquidStaking *LiquidStakingFilterer) FilterSetUnstakeEpoch(opts *bind.FilterOpts, unstakeEpoch []*big.Int, unstakeEpochPrevious []*big.Int, epochInterval []*big.Int) (*LiquidStakingSetUnstakeEpochIterator, error) {

	var unstakeEpochRule []interface{}
	for _, unstakeEpochItem := range unstakeEpoch {
		unstakeEpochRule = append(unstakeEpochRule, unstakeEpochItem)
	}
	var unstakeEpochPreviousRule []interface{}
	for _, unstakeEpochPreviousItem := range unstakeEpochPrevious {
		unstakeEpochPreviousRule = append(unstakeEpochPreviousRule, unstakeEpochPreviousItem)
	}
	var epochIntervalRule []interface{}
	for _, epochIntervalItem := range epochInterval {
		epochIntervalRule = append(epochIntervalRule, epochIntervalItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "SetUnstakeEpoch", unstakeEpochRule, unstakeEpochPreviousRule, epochIntervalRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingSetUnstakeEpochIterator{contract: _LiquidStaking.contract, event: "SetUnstakeEpoch", logs: logs, sub: sub}, nil
}

// WatchSetUnstakeEpoch is a free log subscription operation binding the contract event 0x342db083975cc0697d0f81d0876299ecbf95d909e920426089062e540904497c.
//
// Solidity: event SetUnstakeEpoch(uint256 indexed unstakeEpoch, uint256 indexed unstakeEpochPrevious, uint256 indexed epochInterval)
func (_LiquidStaking *LiquidStakingFilterer) WatchSetUnstakeEpoch(opts *bind.WatchOpts, sink chan<- *LiquidStakingSetUnstakeEpoch, unstakeEpoch []*big.Int, unstakeEpochPrevious []*big.Int, epochInterval []*big.Int) (event.Subscription, error) {

	var unstakeEpochRule []interface{}
	for _, unstakeEpochItem := range unstakeEpoch {
		unstakeEpochRule = append(unstakeEpochRule, unstakeEpochItem)
	}
	var unstakeEpochPreviousRule []interface{}
	for _, unstakeEpochPreviousItem := range unstakeEpochPrevious {
		unstakeEpochPreviousRule = append(unstakeEpochPreviousRule, unstakeEpochPreviousItem)
	}
	var epochIntervalRule []interface{}
	for _, epochIntervalItem := range epochInterval {
		epochIntervalRule = append(epochIntervalRule, epochIntervalItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "SetUnstakeEpoch", unstakeEpochRule, unstakeEpochPreviousRule, epochIntervalRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingSetUnstakeEpoch)
				if err := _LiquidStaking.contract.UnpackLog(event, "SetUnstakeEpoch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetUnstakeEpoch is a log parse operation binding the contract event 0x342db083975cc0697d0f81d0876299ecbf95d909e920426089062e540904497c.
//
// Solidity: event SetUnstakeEpoch(uint256 indexed unstakeEpoch, uint256 indexed unstakeEpochPrevious, uint256 indexed epochInterval)
func (_LiquidStaking *LiquidStakingFilterer) ParseSetUnstakeEpoch(log types.Log) (*LiquidStakingSetUnstakeEpoch, error) {
	event := new(LiquidStakingSetUnstakeEpoch)
	if err := _LiquidStaking.contract.UnpackLog(event, "SetUnstakeEpoch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingSetUnstakingLockTimeIterator is returned from FilterSetUnstakingLockTime and is used to iterate over the raw logs and unpacked data for SetUnstakingLockTime events raised by the LiquidStaking contract.
type LiquidStakingSetUnstakingLockTimeIterator struct {
	Event *LiquidStakingSetUnstakingLockTime // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingSetUnstakingLockTimeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingSetUnstakingLockTime)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingSetUnstakingLockTime)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingSetUnstakingLockTimeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingSetUnstakingLockTimeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingSetUnstakingLockTime represents a SetUnstakingLockTime event raised by the LiquidStaking contract.
type LiquidStakingSetUnstakingLockTime struct {
	UnstakingLockTime *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterSetUnstakingLockTime is a free log retrieval operation binding the contract event 0xfced4450741ad45114dbc966176eb1d16dc4f2259bc864aa395736a7e8355bbb.
//
// Solidity: event SetUnstakingLockTime(uint256 indexed unstakingLockTime)
func (_LiquidStaking *LiquidStakingFilterer) FilterSetUnstakingLockTime(opts *bind.FilterOpts, unstakingLockTime []*big.Int) (*LiquidStakingSetUnstakingLockTimeIterator, error) {

	var unstakingLockTimeRule []interface{}
	for _, unstakingLockTimeItem := range unstakingLockTime {
		unstakingLockTimeRule = append(unstakingLockTimeRule, unstakingLockTimeItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "SetUnstakingLockTime", unstakingLockTimeRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingSetUnstakingLockTimeIterator{contract: _LiquidStaking.contract, event: "SetUnstakingLockTime", logs: logs, sub: sub}, nil
}

// WatchSetUnstakingLockTime is a free log subscription operation binding the contract event 0xfced4450741ad45114dbc966176eb1d16dc4f2259bc864aa395736a7e8355bbb.
//
// Solidity: event SetUnstakingLockTime(uint256 indexed unstakingLockTime)
func (_LiquidStaking *LiquidStakingFilterer) WatchSetUnstakingLockTime(opts *bind.WatchOpts, sink chan<- *LiquidStakingSetUnstakingLockTime, unstakingLockTime []*big.Int) (event.Subscription, error) {

	var unstakingLockTimeRule []interface{}
	for _, unstakingLockTimeItem := range unstakingLockTime {
		unstakingLockTimeRule = append(unstakingLockTimeRule, unstakingLockTimeItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "SetUnstakingLockTime", unstakingLockTimeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingSetUnstakingLockTime)
				if err := _LiquidStaking.contract.UnpackLog(event, "SetUnstakingLockTime", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetUnstakingLockTime is a log parse operation binding the contract event 0xfced4450741ad45114dbc966176eb1d16dc4f2259bc864aa395736a7e8355bbb.
//
// Solidity: event SetUnstakingLockTime(uint256 indexed unstakingLockTime)
func (_LiquidStaking *LiquidStakingFilterer) ParseSetUnstakingLockTime(log types.Log) (*LiquidStakingSetUnstakingLockTime, error) {
	event := new(LiquidStakingSetUnstakingLockTime)
	if err := _LiquidStaking.contract.UnpackLog(event, "SetUnstakingLockTime", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingStakeDirectIterator is returned from FilterStakeDirect and is used to iterate over the raw logs and unpacked data for StakeDirect events raised by the LiquidStaking contract.
type LiquidStakingStakeDirectIterator struct {
	Event *LiquidStakingStakeDirect // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingStakeDirectIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingStakeDirect)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingStakeDirect)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingStakeDirectIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingStakeDirectIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingStakeDirect represents a StakeDirect event raised by the LiquidStaking contract.
type LiquidStakingStakeDirect struct {
	To             common.Address
	StakingAmount  *big.Int
	AmountStaked   *big.Int
	WrappingAmount *big.Int
	AmountWrapped  *big.Int
	Timestamp      *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterStakeDirect is a free log retrieval operation binding the contract event 0x7fd980aad29d3d6256218a13a2d1d8dac5cf5d7411c1c122e2503923952eddd5.
//
// Solidity: event StakeDirect(address to, uint256 stakingAmount, uint256 amountStaked, uint256 wrappingAmount, uint256 amountWrapped, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) FilterStakeDirect(opts *bind.FilterOpts) (*LiquidStakingStakeDirectIterator, error) {

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "StakeDirect")
	if err != nil {
		return nil, err
	}
	return &LiquidStakingStakeDirectIterator{contract: _LiquidStaking.contract, event: "StakeDirect", logs: logs, sub: sub}, nil
}

// WatchStakeDirect is a free log subscription operation binding the contract event 0x7fd980aad29d3d6256218a13a2d1d8dac5cf5d7411c1c122e2503923952eddd5.
//
// Solidity: event StakeDirect(address to, uint256 stakingAmount, uint256 amountStaked, uint256 wrappingAmount, uint256 amountWrapped, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) WatchStakeDirect(opts *bind.WatchOpts, sink chan<- *LiquidStakingStakeDirect) (event.Subscription, error) {

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "StakeDirect")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingStakeDirect)
				if err := _LiquidStaking.contract.UnpackLog(event, "StakeDirect", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStakeDirect is a log parse operation binding the contract event 0x7fd980aad29d3d6256218a13a2d1d8dac5cf5d7411c1c122e2503923952eddd5.
//
// Solidity: event StakeDirect(address to, uint256 stakingAmount, uint256 amountStaked, uint256 wrappingAmount, uint256 amountWrapped, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) ParseStakeDirect(log types.Log) (*LiquidStakingStakeDirect, error) {
	event := new(LiquidStakingStakeDirect)
	if err := _LiquidStaking.contract.UnpackLog(event, "StakeDirect", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingStakeDirectInBatchIterator is returned from FilterStakeDirectInBatch and is used to iterate over the raw logs and unpacked data for StakeDirectInBatch events raised by the LiquidStaking contract.
type LiquidStakingStakeDirectInBatchIterator struct {
	Event *LiquidStakingStakeDirectInBatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingStakeDirectInBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingStakeDirectInBatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingStakeDirectInBatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingStakeDirectInBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingStakeDirectInBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingStakeDirectInBatch represents a StakeDirectInBatch event raised by the LiquidStaking contract.
type LiquidStakingStakeDirectInBatch struct {
	ToAddressses    []common.Address
	StakingAmounts  []*big.Int
	AmountsStaked   []*big.Int
	WrappingAmounts []*big.Int
	AmountsWrapped  []*big.Int
	Timestamp       *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterStakeDirectInBatch is a free log retrieval operation binding the contract event 0xae1ff5caf707784d6d12e6dd90a042ce2994fc6f3d32634ea155dd6e352ed6cc.
//
// Solidity: event StakeDirectInBatch(address[] toAddressses, uint256[] stakingAmounts, uint256[] amountsStaked, uint256[] wrappingAmounts, uint256[] amountsWrapped, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) FilterStakeDirectInBatch(opts *bind.FilterOpts) (*LiquidStakingStakeDirectInBatchIterator, error) {

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "StakeDirectInBatch")
	if err != nil {
		return nil, err
	}
	return &LiquidStakingStakeDirectInBatchIterator{contract: _LiquidStaking.contract, event: "StakeDirectInBatch", logs: logs, sub: sub}, nil
}

// WatchStakeDirectInBatch is a free log subscription operation binding the contract event 0xae1ff5caf707784d6d12e6dd90a042ce2994fc6f3d32634ea155dd6e352ed6cc.
//
// Solidity: event StakeDirectInBatch(address[] toAddressses, uint256[] stakingAmounts, uint256[] amountsStaked, uint256[] wrappingAmounts, uint256[] amountsWrapped, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) WatchStakeDirectInBatch(opts *bind.WatchOpts, sink chan<- *LiquidStakingStakeDirectInBatch) (event.Subscription, error) {

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "StakeDirectInBatch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingStakeDirectInBatch)
				if err := _LiquidStaking.contract.UnpackLog(event, "StakeDirectInBatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStakeDirectInBatch is a log parse operation binding the contract event 0xae1ff5caf707784d6d12e6dd90a042ce2994fc6f3d32634ea155dd6e352ed6cc.
//
// Solidity: event StakeDirectInBatch(address[] toAddressses, uint256[] stakingAmounts, uint256[] amountsStaked, uint256[] wrappingAmounts, uint256[] amountsWrapped, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) ParseStakeDirectInBatch(log types.Log) (*LiquidStakingStakeDirectInBatch, error) {
	event := new(LiquidStakingStakeDirectInBatch)
	if err := _LiquidStaking.contract.UnpackLog(event, "StakeDirectInBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingStakeTokensIterator is returned from FilterStakeTokens and is used to iterate over the raw logs and unpacked data for StakeTokens events raised by the LiquidStaking contract.
type LiquidStakingStakeTokensIterator struct {
	Event *LiquidStakingStakeTokens // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingStakeTokensIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingStakeTokens)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingStakeTokens)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingStakeTokensIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingStakeTokensIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingStakeTokens represents a StakeTokens event raised by the LiquidStaking contract.
type LiquidStakingStakeTokens struct {
	AccountAddress common.Address
	Tokens         *big.Int
	FinalTokens    *big.Int
	Timestamp      *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterStakeTokens is a free log retrieval operation binding the contract event 0x5fe79871cd2431c06447cbcf2557091da5d2ed5bc640f1028f42665913786e42.
//
// Solidity: event StakeTokens(address indexed accountAddress, uint256 indexed tokens, uint256 indexed finalTokens, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) FilterStakeTokens(opts *bind.FilterOpts, accountAddress []common.Address, tokens []*big.Int, finalTokens []*big.Int) (*LiquidStakingStakeTokensIterator, error) {

	var accountAddressRule []interface{}
	for _, accountAddressItem := range accountAddress {
		accountAddressRule = append(accountAddressRule, accountAddressItem)
	}
	var tokensRule []interface{}
	for _, tokensItem := range tokens {
		tokensRule = append(tokensRule, tokensItem)
	}
	var finalTokensRule []interface{}
	for _, finalTokensItem := range finalTokens {
		finalTokensRule = append(finalTokensRule, finalTokensItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "StakeTokens", accountAddressRule, tokensRule, finalTokensRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingStakeTokensIterator{contract: _LiquidStaking.contract, event: "StakeTokens", logs: logs, sub: sub}, nil
}

// WatchStakeTokens is a free log subscription operation binding the contract event 0x5fe79871cd2431c06447cbcf2557091da5d2ed5bc640f1028f42665913786e42.
//
// Solidity: event StakeTokens(address indexed accountAddress, uint256 indexed tokens, uint256 indexed finalTokens, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) WatchStakeTokens(opts *bind.WatchOpts, sink chan<- *LiquidStakingStakeTokens, accountAddress []common.Address, tokens []*big.Int, finalTokens []*big.Int) (event.Subscription, error) {

	var accountAddressRule []interface{}
	for _, accountAddressItem := range accountAddress {
		accountAddressRule = append(accountAddressRule, accountAddressItem)
	}
	var tokensRule []interface{}
	for _, tokensItem := range tokens {
		tokensRule = append(tokensRule, tokensItem)
	}
	var finalTokensRule []interface{}
	for _, finalTokensItem := range finalTokens {
		finalTokensRule = append(finalTokensRule, finalTokensItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "StakeTokens", accountAddressRule, tokensRule, finalTokensRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingStakeTokens)
				if err := _LiquidStaking.contract.UnpackLog(event, "StakeTokens", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStakeTokens is a log parse operation binding the contract event 0x5fe79871cd2431c06447cbcf2557091da5d2ed5bc640f1028f42665913786e42.
//
// Solidity: event StakeTokens(address indexed accountAddress, uint256 indexed tokens, uint256 indexed finalTokens, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) ParseStakeTokens(log types.Log) (*LiquidStakingStakeTokens, error) {
	event := new(LiquidStakingStakeTokens)
	if err := _LiquidStaking.contract.UnpackLog(event, "StakeTokens", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the LiquidStaking contract.
type LiquidStakingUnpausedIterator struct {
	Event *LiquidStakingUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingUnpaused represents a Unpaused event raised by the LiquidStaking contract.
type LiquidStakingUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_LiquidStaking *LiquidStakingFilterer) FilterUnpaused(opts *bind.FilterOpts) (*LiquidStakingUnpausedIterator, error) {

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &LiquidStakingUnpausedIterator{contract: _LiquidStaking.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_LiquidStaking *LiquidStakingFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *LiquidStakingUnpaused) (event.Subscription, error) {

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingUnpaused)
				if err := _LiquidStaking.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_LiquidStaking *LiquidStakingFilterer) ParseUnpaused(log types.Log) (*LiquidStakingUnpaused, error) {
	event := new(LiquidStakingUnpaused)
	if err := _LiquidStaking.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingUnstakeTokensIterator is returned from FilterUnstakeTokens and is used to iterate over the raw logs and unpacked data for UnstakeTokens events raised by the LiquidStaking contract.
type LiquidStakingUnstakeTokensIterator struct {
	Event *LiquidStakingUnstakeTokens // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingUnstakeTokensIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingUnstakeTokens)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingUnstakeTokens)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingUnstakeTokensIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingUnstakeTokensIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingUnstakeTokens represents a UnstakeTokens event raised by the LiquidStaking contract.
type LiquidStakingUnstakeTokens struct {
	AccountAddress common.Address
	Tokens         *big.Int
	FinalTokens    *big.Int
	Timestamp      *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUnstakeTokens is a free log retrieval operation binding the contract event 0x98850af53deeaf736beb3ea22091c900e8e5cfdf8153342daa366dc1d83a39fd.
//
// Solidity: event UnstakeTokens(address indexed accountAddress, uint256 indexed tokens, uint256 indexed finalTokens, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) FilterUnstakeTokens(opts *bind.FilterOpts, accountAddress []common.Address, tokens []*big.Int, finalTokens []*big.Int) (*LiquidStakingUnstakeTokensIterator, error) {

	var accountAddressRule []interface{}
	for _, accountAddressItem := range accountAddress {
		accountAddressRule = append(accountAddressRule, accountAddressItem)
	}
	var tokensRule []interface{}
	for _, tokensItem := range tokens {
		tokensRule = append(tokensRule, tokensItem)
	}
	var finalTokensRule []interface{}
	for _, finalTokensItem := range finalTokens {
		finalTokensRule = append(finalTokensRule, finalTokensItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "UnstakeTokens", accountAddressRule, tokensRule, finalTokensRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingUnstakeTokensIterator{contract: _LiquidStaking.contract, event: "UnstakeTokens", logs: logs, sub: sub}, nil
}

// WatchUnstakeTokens is a free log subscription operation binding the contract event 0x98850af53deeaf736beb3ea22091c900e8e5cfdf8153342daa366dc1d83a39fd.
//
// Solidity: event UnstakeTokens(address indexed accountAddress, uint256 indexed tokens, uint256 indexed finalTokens, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) WatchUnstakeTokens(opts *bind.WatchOpts, sink chan<- *LiquidStakingUnstakeTokens, accountAddress []common.Address, tokens []*big.Int, finalTokens []*big.Int) (event.Subscription, error) {

	var accountAddressRule []interface{}
	for _, accountAddressItem := range accountAddress {
		accountAddressRule = append(accountAddressRule, accountAddressItem)
	}
	var tokensRule []interface{}
	for _, tokensItem := range tokens {
		tokensRule = append(tokensRule, tokensItem)
	}
	var finalTokensRule []interface{}
	for _, finalTokensItem := range finalTokens {
		finalTokensRule = append(finalTokensRule, finalTokensItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "UnstakeTokens", accountAddressRule, tokensRule, finalTokensRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingUnstakeTokens)
				if err := _LiquidStaking.contract.UnpackLog(event, "UnstakeTokens", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnstakeTokens is a log parse operation binding the contract event 0x98850af53deeaf736beb3ea22091c900e8e5cfdf8153342daa366dc1d83a39fd.
//
// Solidity: event UnstakeTokens(address indexed accountAddress, uint256 indexed tokens, uint256 indexed finalTokens, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) ParseUnstakeTokens(log types.Log) (*LiquidStakingUnstakeTokens, error) {
	event := new(LiquidStakingUnstakeTokens)
	if err := _LiquidStaking.contract.UnpackLog(event, "UnstakeTokens", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStakingWithdrawUnstakeTokensIterator is returned from FilterWithdrawUnstakeTokens and is used to iterate over the raw logs and unpacked data for WithdrawUnstakeTokens events raised by the LiquidStaking contract.
type LiquidStakingWithdrawUnstakeTokensIterator struct {
	Event *LiquidStakingWithdrawUnstakeTokens // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LiquidStakingWithdrawUnstakeTokensIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStakingWithdrawUnstakeTokens)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LiquidStakingWithdrawUnstakeTokens)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LiquidStakingWithdrawUnstakeTokensIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStakingWithdrawUnstakeTokensIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStakingWithdrawUnstakeTokens represents a WithdrawUnstakeTokens event raised by the LiquidStaking contract.
type LiquidStakingWithdrawUnstakeTokens struct {
	AccountAddress common.Address
	Tokens         *big.Int
	Timestamp      *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterWithdrawUnstakeTokens is a free log retrieval operation binding the contract event 0x31453b59c0eb931edc765cc2400876a524254457b5534d660c9da1d74ab4666b.
//
// Solidity: event WithdrawUnstakeTokens(address indexed accountAddress, uint256 indexed tokens, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) FilterWithdrawUnstakeTokens(opts *bind.FilterOpts, accountAddress []common.Address, tokens []*big.Int) (*LiquidStakingWithdrawUnstakeTokensIterator, error) {

	var accountAddressRule []interface{}
	for _, accountAddressItem := range accountAddress {
		accountAddressRule = append(accountAddressRule, accountAddressItem)
	}
	var tokensRule []interface{}
	for _, tokensItem := range tokens {
		tokensRule = append(tokensRule, tokensItem)
	}

	logs, sub, err := _LiquidStaking.contract.FilterLogs(opts, "WithdrawUnstakeTokens", accountAddressRule, tokensRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStakingWithdrawUnstakeTokensIterator{contract: _LiquidStaking.contract, event: "WithdrawUnstakeTokens", logs: logs, sub: sub}, nil
}

// WatchWithdrawUnstakeTokens is a free log subscription operation binding the contract event 0x31453b59c0eb931edc765cc2400876a524254457b5534d660c9da1d74ab4666b.
//
// Solidity: event WithdrawUnstakeTokens(address indexed accountAddress, uint256 indexed tokens, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) WatchWithdrawUnstakeTokens(opts *bind.WatchOpts, sink chan<- *LiquidStakingWithdrawUnstakeTokens, accountAddress []common.Address, tokens []*big.Int) (event.Subscription, error) {

	var accountAddressRule []interface{}
	for _, accountAddressItem := range accountAddress {
		accountAddressRule = append(accountAddressRule, accountAddressItem)
	}
	var tokensRule []interface{}
	for _, tokensItem := range tokens {
		tokensRule = append(tokensRule, tokensItem)
	}

	logs, sub, err := _LiquidStaking.contract.WatchLogs(opts, "WithdrawUnstakeTokens", accountAddressRule, tokensRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStakingWithdrawUnstakeTokens)
				if err := _LiquidStaking.contract.UnpackLog(event, "WithdrawUnstakeTokens", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawUnstakeTokens is a log parse operation binding the contract event 0x31453b59c0eb931edc765cc2400876a524254457b5534d660c9da1d74ab4666b.
//
// Solidity: event WithdrawUnstakeTokens(address indexed accountAddress, uint256 indexed tokens, uint256 timestamp)
func (_LiquidStaking *LiquidStakingFilterer) ParseWithdrawUnstakeTokens(log types.Log) (*LiquidStakingWithdrawUnstakeTokens, error) {
	event := new(LiquidStakingWithdrawUnstakeTokens)
	if err := _LiquidStaking.contract.UnpackLog(event, "WithdrawUnstakeTokens", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
