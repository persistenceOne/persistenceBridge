/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package application

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

// EncodingConfiguration specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfiguration struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Marshaler
	TransactionConfig client.TxConfig

	// nolint fixme: remove deprecated
	// nolint:staticcheck
	Amino *codec.LegacyAmino
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() EncodingConfiguration {
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	encodingConfig := EncodingConfiguration{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TransactionConfig: tx.NewTxConfig(marshaler, tx.DefaultSignModes),
		Amino:             codec.NewLegacyAmino(),
	}
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}
