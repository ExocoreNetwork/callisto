package types

import (
	"fmt"

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

// StakerAsset is a helper struct with the state of a staker's asset. In
// addition to the key (StakerID, AssetID) and the base structure used by the
// x/assets module, it also contains the lifetime slashed amount and the
// delegated amount. The lifetime slashed amount is derived from events,
// and the deposited amount is calculated from the other parts.
type StakerAsset struct {
	StakerID            string
	AssetID             string
	Deposited           string
	Withdrawable        string
	PendingUndelegation string
	// AdditionalSlashed is non-cumulative in nature and may be "" to indicate no slashing
	AdditionalSlashed string
}

// NewStakerAssetFromInfo creates a new StakerAsset instance from the given
// StakerID, AssetID, StakerAssetInfo and additionalSlashed amount
func NewStakerAssetFromInfo(
	stakerID string, assetID string,
	info assetstypes.StakerAssetInfo,
	additionalSlashed string, // may be "" to indicate no slashing
) *StakerAsset {
	return &StakerAsset{
		StakerID:            stakerID,
		AssetID:             assetID,
		Deposited:           info.TotalDepositAmount.String(),
		Withdrawable:        info.WithdrawableAmount.String(),
		PendingUndelegation: info.PendingUndelegationAmount.String(),
		AdditionalSlashed:   additionalSlashed,
	}
}

// NewStakerAssetFromStr creates a new StakerAsset instance from the given
// StakerID, AssetID, and string versions of the amounts.
func NewStakerAssetFromStr(
	stakerID string, assetID string,
	deposited string, withdrawable string, pendingUndelegation string,
	additionalSlashed string, // may be "" to indicate no slashing
) *StakerAsset {
	return &StakerAsset{
		StakerID:            stakerID,
		AssetID:             assetID,
		Deposited:           deposited,
		Withdrawable:        withdrawable,
		PendingUndelegation: pendingUndelegation,
		AdditionalSlashed:   additionalSlashed,
	}
}

// OperatorAsset is a helper struct containing string versions of OperatorAssetInfo
// with indexing by OperatorAddress and AssetID.
type OperatorAsset struct {
	OperatorAddress           string
	AssetID                   string
	TotalAmount               string
	PendingUndelegationAmount string
	TotalShare                string
	SelfShare                 string
}

// NewOperatorAssetFromInfo creates a new OperatorAsset instance from the given
// OperatorAddress, AssetID and OperatorAssetInfo.
func NewOperatorAssetFromInfo(
	operatorAddress string, assetID string,
	info assetstypes.OperatorAssetInfo,
) *OperatorAsset {
	return &OperatorAsset{
		OperatorAddress:           operatorAddress,
		AssetID:                   assetID,
		TotalAmount:               info.TotalAmount.String(),
		PendingUndelegationAmount: info.PendingUndelegationAmount.String(),
		TotalShare:                info.TotalShare.String(),
		SelfShare:                 info.OperatorShare.String(),
	}
}

// NewOperatorAssetFromStr creates a new OperatorAsset instance from the given
// OperatorAddress, AssetID, and string versions of the amounts + shares.
func NewOperatorAssetFromStr(
	operatorAddress string, assetID string,
	totalAmount string, pendingUndelegationAmount string,
	totalShare string, selfShare string,
) *OperatorAsset {
	return &OperatorAsset{
		OperatorAddress:           operatorAddress,
		AssetID:                   assetID,
		TotalAmount:               totalAmount,
		PendingUndelegationAmount: pendingUndelegationAmount,
		TotalShare:                totalShare,
		SelfShare:                 selfShare,
	}
}

// ClientChain is the string version of assetstypes.ClientChainInfo
type ClientChain struct {
	Name               string
	MetaInfo           string
	ChainId            string
	ExocoreChainIndex  string
	FinalizationBlocks string
	LayerZeroChainID   string
	SignatureType      string
	AddressLength      string
}

// NewClientChain converts assetstypes.ClientChainInfo to ClientChain
func NewClientChain(info assetstypes.ClientChainInfo) *ClientChain {
	return &ClientChain{
		Name:               info.Name,
		MetaInfo:           info.MetaInfo,
		ChainId:            fmt.Sprintf("%d", info.ChainId),
		ExocoreChainIndex:  fmt.Sprintf("%d", info.ExocoreChainIndex),
		FinalizationBlocks: fmt.Sprintf("%d", info.FinalizationBlocks),
		LayerZeroChainID:   fmt.Sprintf("%d", info.LayerZeroChainID),
		SignatureType:      info.SignatureType,
		AddressLength:      fmt.Sprintf("%d", info.AddressLength),
	}
}

// NewClientChainFromStr creates a new ClientChain instance from the given
// string versions of the fields.
func NewClientChainFromStr(
	name string, metaInfo string, chainId string, exocoreChainIndex string,
	finalizationBlocks string, layerZeroChainID string, signatureType string,
	addressLength string,
) *ClientChain {
	return &ClientChain{
		Name:               name,
		MetaInfo:           metaInfo,
		ChainId:            chainId,
		ExocoreChainIndex:  exocoreChainIndex,
		FinalizationBlocks: finalizationBlocks,
		LayerZeroChainID:   layerZeroChainID,
		SignatureType:      signatureType,
		AddressLength:      addressLength,
	}
}

// AssetsToken is a helper struct to represent an assetstypes.StakingAssetInfo with an assetID
// string.
type AssetsToken struct {
	AssetID           string
	Name              string
	Symbol            string
	Address           string
	Decimals          string
	LayerZeroChainID  string
	ExocoreChainIndex string
	MetaInfo          string
	Amount            string
}

// NewAssetsToken creates a new AssetsToken instance from the given assetstypes.StakingAssetInfo
func NewAssetsToken(info *assetstypes.StakingAssetInfo) *AssetsToken {
	basic := info.AssetBasicInfo
	return &AssetsToken{
		AssetID:           basic.AssetID(),
		Name:              basic.Name,
		Symbol:            basic.Symbol,
		Address:           basic.Address,
		Decimals:          fmt.Sprintf("%d", basic.Decimals),
		LayerZeroChainID:  fmt.Sprintf("%d", basic.LayerZeroChainID),
		ExocoreChainIndex: fmt.Sprintf("%d", basic.ExocoreChainIndex),
		MetaInfo:          basic.MetaInfo,
		Amount:            info.StakingTotalAmount.String(),
	}
}

// NewAssetsTokenFromStr creates a new AssetsToken instance from the given string versions
// of the fields.
func NewAssetsTokenFromStr(
	assetID string, name string, symbol string, address string, decimals string,
	layerZeroChainID string, exocoreChainIndex string, metaInfo string, amount string,
) *AssetsToken {
	return &AssetsToken{
		AssetID:           assetID,
		Name:              name,
		Symbol:            symbol,
		Address:           address,
		Decimals:          decimals,
		LayerZeroChainID:  layerZeroChainID,
		ExocoreChainIndex: exocoreChainIndex,
		MetaInfo:          metaInfo,
		Amount:            amount,
	}
}
