package types

import (
	"fmt"

	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// DelegationState is the state indexed by staker_id + asset_id + operator_addr
// It tracks the undelegatable share and the amount of tokens that are unbonding.
// It is equivaent to `delegationtype.DelegationAmounts`
type DelegationState struct {
	StakerID               string
	AssetID                string
	OperatorAddr           string
	UndelegatableShare     string
	WaitUndelegationAmount string
}

// NewDelegationStateFromStr creates a new DelegationState instance using the given values in
// string format.
func NewDelegationStateFromStr(
	stakerID, assetID, operatorAddr, undelegatableShare, waitUndelegationAmount string,
) *DelegationState {
	return &DelegationState{
		StakerID:               stakerID,
		AssetID:                assetID,
		OperatorAddr:           operatorAddr,
		UndelegatableShare:     undelegatableShare,
		WaitUndelegationAmount: waitUndelegationAmount,
	}
}

// NewDelegationState creates a new DelegationState instance using the given values in
// `delegationtypes.DelegationAmounts` format.
func NewDelegationState(
	stakerID, assetID, operatorAddr string,
	delegationAmounts *delegationtypes.DelegationAmounts,
) *DelegationState {
	return NewDelegationStateFromStr(
		stakerID, assetID, operatorAddr,
		delegationAmounts.UndelegatableShare.String(),
		delegationAmounts.WaitUndelegationAmount.String(),
	)
}

// UndelegationRecord is the record indexed by record_id. It is the equivalent of
// `delegationtypes.UndelegationRecord` with additional `RecordID` and `HoldCount` fields.
type UndelegationRecord struct {
	RecordID                 string
	StakerID                 string
	AssetID                  string
	OperatorAddr             string
	TxHash                   string
	BlockNumber              string
	CompletedEpochIdentifier string
	CompletedEpochNumber     string
	UndelegationID           string
	Amount                   string
	ActualCompletedAmount    string
	HoldCount                string
}

// NewUndelegationRecordFromStr creates a new UndelegationRecord instance using the given values in
// string format.
func NewUndelegationRecordFromStr(
	recordID, stakerID, assetID,
	operatorAddr, txHash, blockNumber,
	completedEpochIdentifier, completedEpochNumber,
	undelegationID, amount, actualCompletedAmount,
	holdCount string,
) *UndelegationRecord {
	return &UndelegationRecord{
		RecordID:                 recordID,
		StakerID:                 stakerID,
		AssetID:                  assetID,
		OperatorAddr:             operatorAddr,
		TxHash:                   txHash,
		BlockNumber:              blockNumber,
		CompletedEpochIdentifier: completedEpochIdentifier,
		CompletedEpochNumber:     completedEpochNumber,
		UndelegationID:           undelegationID,
		Amount:                   amount,
	}
}

// NewUndelegationRecord creates a new UndelegationRecord instance using the given values in
// `delegationtypes.UndelegationRecord` format.
func NewUndelegationRecord(
	undelegationRecord *delegationtypes.UndelegationRecord,
	holdCount uint64,
) *UndelegationRecord {
	return NewUndelegationRecordFromStr(
		hexutil.Encode(undelegationRecord.GetKey()),
		undelegationRecord.StakerId,
		undelegationRecord.AssetId,
		undelegationRecord.OperatorAddr,
		undelegationRecord.TxHash,
		fmt.Sprintf("%d", undelegationRecord.BlockNumber),
		undelegationRecord.CompletedEpochIdentifier,
		fmt.Sprintf("%d", undelegationRecord.CompletedEpochNumber),
		fmt.Sprintf("%d", undelegationRecord.UndelegationId),
		undelegationRecord.Amount.String(),
		undelegationRecord.ActualCompletedAmount.String(),
		fmt.Sprintf("%d", holdCount),
	)
}

// ExoAssetDelegation is the delegation state indexed by staker_id. It is equivalent to
// the staker_assets table, but without the asset_id column.
// It tracks the delegated amount, the amount pending undelegation, and the amount slashed (
// via events, so not pictured here).
type ExoAssetDelegation struct {
	StakerID            string
	Delegated           string
	PendingUndelegation string
}

// NewExoAssetDelegationFromStr creates a new ExoAssetDelegation instance using the given values in
// string format.
func NewExoAssetDelegationFromStr(
	stakerID, delegated, pendingUndelegation string,
) *ExoAssetDelegation {
	return &ExoAssetDelegation{
		StakerID:            stakerID,
		Delegated:           delegated,
		PendingUndelegation: pendingUndelegation,
	}
}
