package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"

	"github.com/ExocoreNetwork/exocore/x/operator/types"
)

// Operator is the string version of operatortypes.OperatorInfo
type Operator struct {
	EarningsAddress string
	ApproveAddress  string
	MetaInfo        string
	Rate            string
	MaxRate         string
	MaxChangeRate   string
	// prefixed with Commission to avoid conflict with other fields
	CommissionUpdateTime string
	// TODO: add client chain earnings address
}

// NewOperator converts an operatortypes.OperatorInfo into a Operator
func NewOperator(t *types.OperatorInfo) *Operator {
	return &Operator{
		EarningsAddress:      t.EarningsAddr,
		ApproveAddress:       t.ApproveAddr,
		MetaInfo:             t.OperatorMetaInfo,
		Rate:                 t.Commission.Rate.String(),
		MaxRate:              t.Commission.MaxRate.String(),
		MaxChangeRate:        t.Commission.MaxChangeRate.String(),
		CommissionUpdateTime: t.Commission.UpdateTime.String(),
	}
}

// NewOperatorFromStr creates a new Operator instance from the given strings
func NewOperatorFromStr(
	earningsAddress string, approveAddress string, metaInfo string,
	rate string, maxRate string, maxChangeRate string, commissionUpdateTime string,
) *Operator {
	return &Operator{
		EarningsAddress:      earningsAddress,
		ApproveAddress:       approveAddress,
		MetaInfo:             metaInfo,
		Rate:                 rate,
		MaxRate:              maxRate,
		MaxChangeRate:        maxChangeRate,
		CommissionUpdateTime: commissionUpdateTime,
	}
}

type Opted struct {
	OperatorAddress string
	AvsAddress      string
	SlashContract   string
	InHeight        string
	OutHeight       string
	Jailed          bool
}

// NewOpted converts an operatortypes.OptedInfo into a Opted
func NewOpted(
	operatorAddress string, avsAddress string, t *types.OptedInfo,
) *Opted {
	return &Opted{
		OperatorAddress: operatorAddress,
		AvsAddress:      avsAddress,
		SlashContract:   t.SlashContract,
		InHeight:        fmt.Sprintf("%d", t.OptedInHeight),
		OutHeight:       fmt.Sprintf("%d", t.OptedOutHeight),
		Jailed:          t.Jailed,
	}
}

// NewOptedFromStr creates a new Opted instance from the given string
// versions. Note that `jailed` is still a boolean.
func NewOptedFromStr(
	operatorAddress string, avsAddress string,
	slashContract string, inHeight string,
	outHeight string, jailed bool,
) *Opted {
	if outHeight == "" {
		outHeight = fmt.Sprintf("%d", types.DefaultOptedOutHeight)
	}
	return &Opted{
		OperatorAddress: operatorAddress,
		AvsAddress:      avsAddress,
		SlashContract:   slashContract,
		InHeight:        inHeight,
		OutHeight:       outHeight,
		Jailed:          jailed,
	}
}

// OperatorUSDValue is the string version of types.OperatorOptedUSDValue
type OperatorUSDValue struct {
	OperatorAddress string
	AvsAddress      string
	SelfUSDValue    string
	TotalUSDValue   string
	ActiveUSDValue  string
	// derived
	OtherUSDValue string
}

// NewOperatorUSDValueFromStr creates a new OperatorUSDValue instance from the given strings
func NewOperatorUSDValueFromStr(
	operatorAddr string, avsAddr string,
	selfUsdValue string, totalUsdValue string, activeUsdValue string,
) *OperatorUSDValue {
	totalUsdDec, err := sdkmath.LegacyNewDecFromStr(totalUsdValue)
	if err != nil {
		return nil
	}
	selfUsdDec, err := sdkmath.LegacyNewDecFromStr(selfUsdValue)
	if err != nil {
		return nil
	}
	otherUsdDec := totalUsdDec.Sub(selfUsdDec)
	return &OperatorUSDValue{
		OperatorAddress: operatorAddr,
		AvsAddress:      avsAddr,
		SelfUSDValue:    selfUsdValue,
		TotalUSDValue:   totalUsdValue,
		ActiveUSDValue:  activeUsdValue,
		OtherUSDValue:   otherUsdDec.String(),
	}
}

// NewOperatorUSDValue creates a new OperatorUSDValue instance from the given values
func NewOperatorUSDValue(
	operatorAddr string, avsAddr string,
	operatorUSDValue *types.OperatorOptedUSDValue,
) *OperatorUSDValue {
	return &OperatorUSDValue{
		OperatorAddress: operatorAddr,
		AvsAddress:      avsAddr,
		SelfUSDValue:    operatorUSDValue.SelfUSDValue.String(),
		TotalUSDValue:   operatorUSDValue.TotalUSDValue.String(),
		ActiveUSDValue:  operatorUSDValue.ActiveUSDValue.String(),
		OtherUSDValue:   operatorUSDValue.TotalUSDValue.Sub(operatorUSDValue.SelfUSDValue).String(),
	}
}

type AvsUSDValue struct {
	AvsAddress string
	USDValue   string
}

// NewAvsUSDValueFromStr creates a new AvsUSDValue instance from the given strings
func NewAvsUSDValueFromStr(
	avsAddr string, usdValue string,
) *AvsUSDValue {
	return &AvsUSDValue{
		AvsAddress: avsAddr,
		USDValue:   usdValue,
	}
}

// NewAvsUSDValue creates a new AvsUSDValue instance from the given values
func NewAvsUSDValue(
	avsUSDValue *types.AVSUSDValue,
) *AvsUSDValue {
	return &AvsUSDValue{
		AvsAddress: avsUSDValue.AVSAddr,
		USDValue:   avsUSDValue.Value.String(),
	}
}
