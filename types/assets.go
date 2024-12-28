package types

import (
	sdkmath "cosmossdk.io/math"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
)

// AssetsParams represents the x/assets parameters
type AssetsParams struct {
	assetstypes.Params
	Height int64
}

// NewAssetsParams allows to build a new AssetsParams instance
func NewAssetsParams(params assetstypes.Params, height int64) *AssetsParams {
	return &AssetsParams{
		Params: params,
		Height: height,
	}
}

// StakerAsset is a helper struct containing string versions of StakerAssetInfo
// with indexing by AssetID, StakerID, and Height (optional).
type StakerAsset struct {
	StakerID            string
	AssetID             string
	Deposited           string
	Free                string
	Delegated           string
	PendingUndelegation string
	Height              int64
}

// NewStakerAssetFromInfo creates a new StakerAsset instance from the given
// StakerID, AssetID, StakerAssetInfo, and height. The height may be 0.
func NewStakerAssetFromInfo(
	stakerID string, assetID string,
	info assetstypes.StakerAssetInfo, height int64,
) *StakerAsset {
	delegated := info.TotalDepositAmount.
		Sub(info.WithdrawableAmount).
		Sub(info.PendingUndelegationAmount)
	return &StakerAsset{
		StakerID:            stakerID,
		AssetID:             assetID,
		Deposited:           info.TotalDepositAmount.String(),
		Free:                info.WithdrawableAmount.String(),
		Delegated:           delegated.String(),
		PendingUndelegation: info.PendingUndelegationAmount.String(),
		Height:              height,
	}
}

// NewStakerAssetFromStr creates a new StakerAsset instance from the given
// StakerID, AssetID, and string versions of the deposited, free, and pendingUndelegation
// amounts. The height may be 0.
func NewStakerAssetFromStr(
	stakerID string, assetID string,
	deposited string, free string, pendingUndelegation string,
	height int64,
) *StakerAsset {
	// to calculate delegated, we need the sdkmath.Int versions :/
	// we ignore the errors for now but could be improved
	depositedInt, _ := sdkmath.NewIntFromString(deposited)
	freeInt, _ := sdkmath.NewIntFromString(free)
	pendingUndelegationInt, _ := sdkmath.NewIntFromString(pendingUndelegation)
	delegated := depositedInt.Sub(freeInt).Sub(pendingUndelegationInt)
	return &StakerAsset{
		StakerID:            stakerID,
		AssetID:             assetID,
		Deposited:           deposited,
		Free:                free,
		Delegated:           delegated.String(),
		PendingUndelegation: pendingUndelegation,
		Height:              height,
	}
}

// OperatorAsset is a helper struct containing string versions of OperatorAssetInfo
// with indexing by OperatorAddress, AssetID, and Height (optional).
type OperatorAsset struct {
	OperatorAddress     string
	AssetID             string
	Delegated           string
	PendingUndelegation string
	Share               string
	SelfShare           string
	DelegatedShare      string
	Height              int64
}

// NewOperatorAssetFromInfo creates a new OperatorAsset instance from the given
// OperatorAddress, AssetID, OperatorAssetInfo, and height. The height may be 0.
func NewOperatorAssetFromInfo(
	operatorAddress string, assetID string,
	info assetstypes.OperatorAssetInfo, height int64,
) *OperatorAsset {
	delegatedShare := info.TotalShare.Sub(info.OperatorShare)
	return &OperatorAsset{
		OperatorAddress:     operatorAddress,
		AssetID:             assetID,
		Delegated:           info.TotalAmount.String(),
		PendingUndelegation: info.PendingUndelegationAmount.String(),
		Share:               info.TotalShare.String(),
		SelfShare:           info.OperatorShare.String(),
		DelegatedShare:      delegatedShare.String(),
		Height:              height,
	}
}

// NewOperatorAssetFromStr creates a new OperatorAsset instance from the given
// OperatorAddress, AssetID, and string versions of the delegated, pendingUndelegation,
// share, selfShare, and delegatedShare amounts. The height may be 0.
func NewOperatorAssetFromStr(
	operatorAddress string, assetID string,
	delegated string, pendingUndelegation string,
	share string, selfShare string,
	height int64,
) *OperatorAsset {
	// to calculate delegatedShare, we need the sdkmath.Dec versions :/
	shareDec := sdkmath.LegacyMustNewDecFromStr(delegated)
	selfShareDec := sdkmath.LegacyMustNewDecFromStr(selfShare)
	delegatedShare := shareDec.Sub(selfShareDec)
	return &OperatorAsset{
		OperatorAddress:     operatorAddress,
		AssetID:             assetID,
		Delegated:           delegated,
		PendingUndelegation: pendingUndelegation,
		Share:               share,
		SelfShare:           selfShare,
		DelegatedShare:      delegatedShare.String(),
		Height:              height,
	}
}
