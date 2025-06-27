package app

import (
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	evmencoding "github.com/cosmos/evm/encoding"
)

// MakeEncodingConfig creates the EncodingConfig for realio network
func MakeEncodingConfig(evmChainID uint64) testutil.TestEncodingConfig {
	encCfg := evmencoding.MakeConfig(evmChainID)
	// ethcryptocodec.RegisterInterfaces(interfaceRegistry)
	return encCfg
}
