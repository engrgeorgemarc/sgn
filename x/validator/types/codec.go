package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgInitializeCandidate{}, "validator/MsgInitializeCandidate", nil)
	cdc.RegisterConcrete(MsgClaimValidator{}, "validator/MsgClaimValidator", nil)
	cdc.RegisterConcrete(MsgSyncValidator{}, "validator/MsgSyncValidator", nil)
}
