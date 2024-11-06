package sdk

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewGuilinBankSDK, NewSPDBankSDK, NewPinganBankSDK, NewIcbcBankSDK, NewMinShengBankSDK)

func NewGuilinBankSDK() GuilinBankSDK {
	return &guilinBankSDK{}
}

func NewSPDBankSDK() SPDBankSDK { return &spdBankSDK{} }

func NewPinganBankSDK() PinganBankSDK {
	return &pinganBankSDK{}
}

func NewIcbcBankSDK() IcbcBankSDK {
	return &icbcBankSDK{}
}
func NewMinShengBankSDK() MinShengSDK {
	return &minShengSDK{}
}
