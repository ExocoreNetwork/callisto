package assets

import (
	"fmt"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/forbole/callisto/v4/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	if err := m.handleStakerEvents(tx.Height, tx.Events); err != nil {
		return fmt.Errorf("error while handling staker events: %s", err)
	}
	if err := m.handleOperatorEvents(tx.Height, tx.Events); err != nil {
		return fmt.Errorf("error while handling operator events: %s", err)
	}
	return nil
}

// handleStakerEvents filters, parses and indexes the staker events.
func (m *Module) handleStakerEvents(height int64, events []abci.Event) error {
	events = juno.FindEventsByType(events, assetstypes.EventTypeUpdatedStakerAsset)
	for _, event := range events {
		stakerID, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyStakerID)
		if err != nil {
			return fmt.Errorf("error while getting staker ID: %s", err)
		}
		assetID, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while getting asset ID: %s", err)
		}
		depositAmount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyDepositAmount)
		if err != nil {
			return fmt.Errorf("error while getting deposit amount: %s", err)
		}
		withdrawableAmount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyWithdrawableAmount)
		if err != nil {
			return fmt.Errorf("error while getting withdrawable amount: %s", err)
		}
		pendingUndelegationAmount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyPendingUndelegationAmount)
		if err != nil {
			return fmt.Errorf("error while getting pending undelegation amount: %s", err)
		}
		asset := types.NewStakerAssetFromStr(
			stakerID.Value, assetID.Value,
			depositAmount.Value, withdrawableAmount.Value, pendingUndelegationAmount.Value,
			height,
		)
		if err := m.db.SaveStakerAsset(asset); err != nil {
			return fmt.Errorf("error while saving staker asset: %s", err)
		}
	}
	return nil
}

// handleOperatorEvents filters, parses and indexes the operator events.
func (m *Module) handleOperatorEvents(height int64, events []abci.Event) error {
	events = juno.FindEventsByType(events, assetstypes.EventTypeUpdatedOperatorAsset)
	for _, event := range events {
		operatorAddr, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyOperatorAddress)
		if err != nil {
			return fmt.Errorf("error while getting operator ID: %s", err)
		}
		assetID, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while getting asset ID: %s", err)
		}
		totalAmount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyTotalAmount)
		if err != nil {
			return fmt.Errorf("error while getting deposit amount: %s", err)
		}
		pendingUndelegationAmount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyPendingUndelegationAmount)
		if err != nil {
			return fmt.Errorf("error while getting pending undelegation amount: %s", err)
		}
		totalShare, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyTotalShare)
		if err != nil {
			return fmt.Errorf("error while getting total share: %s", err)
		}
		operatorShare, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyOperatorShare)
		if err != nil {
			return fmt.Errorf("error while getting operator share: %s", err)
		}
		asset := types.NewOperatorAssetFromStr(
			operatorAddr.Value, assetID.Value,
			totalAmount.Value, pendingUndelegationAmount.Value,
			totalShare.Value, operatorShare.Value,
			height,
		)
		if err := m.db.SaveOperatorAsset(asset); err != nil {
			return fmt.Errorf("error while saving operator asset: %s", err)
		}
	}
	return nil
}
