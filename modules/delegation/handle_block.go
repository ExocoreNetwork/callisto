package delegation

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	juno "github.com/forbole/juno/v5/types"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	// slashing
	if err := m.handleDelegationStateUpdates(res.BeginBlockEvents); err != nil {
		return fmt.Errorf("error while handling slashing events: %s", err)
	}
	// undelegation maturity
	if err := m.handleDelegationStateUpdates(res.EndBlockEvents); err != nil {
		return fmt.Errorf("error while handling delegation state updates: %s", err)
	}
	// slashing in BeginBlock
	if err := m.handleAllStakersRemovedFromOperatorAsset(res.BeginBlockEvents); err != nil {
		return fmt.Errorf("error while handling all stakers removed from operator asset updates: %s", err)
	}
	// undelegation hold count reduction in EndBlock
	if err := m.handleUndelegationHoldCountChanges(res.EndBlockEvents); err != nil {
		return fmt.Errorf("error while handling undelegation hold count changes: %s", err)
	}
	// undelegation completions in EndBlock
	if err := m.handleUndelegationCompletions(block.Block.Height, res.EndBlockEvents); err != nil {
		return fmt.Errorf("error while handling undelegation completions: %s", err)
	}
	// slashing of undelegations via BeginBlock
	if err := m.handleUndelegationSlashings(res.BeginBlockEvents); err != nil {
		return fmt.Errorf("error while handling undelegation slashings: %s", err)
	}
	// remember that the other thing being slashed is the delegation itself.
	// such a slashing is applied to the operator state (by assetID) and it is
	// tracked in x/asset under `operator_assets`. however, we need to track it
	// on a staker-level.
	// it is applied in the BeginBlock, for example, in response to less signing.
	if err := m.handleDelegationSlashings(block.Block.Height, res.BeginBlockEvents); err != nil {
		return fmt.Errorf("error while handling delegation slashings: %s", err)
	}
	return nil
}

// handleUndelegationCompletions handles the undelegation completions.
// triggered in response to the unbonding epoch duration end in x/delegation's EndBlocker
func (m *Module) handleUndelegationCompletions(height int64, events []abci.Event) error {
	events = juno.FindEventsByType(events, delegationtypes.EventTypeUndelegationMatured)
	for _, event := range events {
		recordID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyRecordID)
		if err != nil {
			return fmt.Errorf("error while finding record ID: %s", err)
		}
		amount, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyAmount)
		if err != nil {
			return fmt.Errorf("error while finding amount: %s", err)
		}
		if err := m.db.MatureUndelegationRecord(recordID.Value, amount.Value, height); err != nil {
			return fmt.Errorf("error while maturing undelegation record: %s", err)
		}
		stakerID, assetID, err := m.db.GetStakerIDAssetIDFromUndelegationRecord(recordID.Value)
		if err != nil {
			return fmt.Errorf("error while getting staker ID and asset ID from undelegation record: %s", err)
		}
		if assetID == assetstypes.ExocoreAssetID {
			// there is no staker asset for exocore asset, so we operate on
			// exo_asset_delegation's pending_undelegation amount.
			operatorAddr, err := m.db.GetOperatorAddrFromUndelegationRecord(recordID.Value)
			if err != nil {
				return fmt.Errorf("error while getting operator address from undelegation record: %s", err)
			}
			if err := m.db.MatureExoAssetUndelegation(stakerID, operatorAddr, amount.Value); err != nil {
				return fmt.Errorf("error while maturing exo asset undelegation: %s", err)
			}
		}
	}
	return nil
}

// handleUndelegationSlashings handles the slashing of undelegation records. Technically, this
// happens in x/operator, however, it modifies the state of x/delegation schema items. So, I
// put it here.
func (m *Module) handleUndelegationSlashings(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeUndelegationSlashed)
	for _, event := range events {
		recordID, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyRecordID)
		if err != nil {
			return fmt.Errorf("error while finding record ID: %s", err)
		}
		amount, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyAmount)
		if err != nil {
			return fmt.Errorf("error while finding amount: %s", err)
		}
		slashedAmount, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeySlashAmount)
		if err != nil {
			return fmt.Errorf("error while finding slashed amount: %s", err)
		}
		if err := m.db.SlashUndelegationRecord(recordID.Value, amount.Value, slashedAmount.Value); err != nil {
			return fmt.Errorf("error while slashing undelegation record: %s", err)
		}
	}
	return nil
}

// handleDelegationSlashings handles the slashing of delegation records. Technically, this
// happens in x/operator, however, it modifies the state of x/delegation-related schema items.
// So, I put it here but it could go in x/assets as well.
func (m *Module) handleDelegationSlashings(height int64, events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeOperatorAssetSlashed)
	for _, event := range events {
		assetID, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while finding asset ID: %s", err)
		}
		operatorAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while finding operator address: %s", err)
		}
		stakerIDs, err := m.db.GetStakersByOperatorAsset(operatorAddr.Value, assetID.Value)
		if err != nil {
			return fmt.Errorf("error while getting stakers by operator asset: %s", err)
		}
		for _, stakerID := range stakerIDs {
			delegatedAmount, err := m.source.GetDelegatedAmount(
				height, stakerID, assetID.Value, operatorAddr.Value,
			)
			if err != nil {
				return fmt.Errorf("error while getting delegated amount: %s", err)
			}
			// previously delegated amount
			prevAmount, err := m.db.GetDelegatedAmount(stakerID, assetID.Value)
			if err != nil {
				return fmt.Errorf("error while getting delegated amount: %s", err)
			}
			slashedAmount := prevAmount.Sub(delegatedAmount)
			if assetID.Value != assetstypes.ExocoreAssetID {
				if err := m.db.SlashStakerDelegation(
					stakerID, assetID.Value, slashedAmount.String(),
				); err != nil {
					return fmt.Errorf("error while accumulating staker lifetime slashing: %s", err)
				}
			} else {
				if err := m.db.SlashExoAssetDelegation(
					stakerID, operatorAddr.Value, slashedAmount.String(),
				); err != nil {
					return fmt.Errorf("error while slashing exo asset delegation: %s", err)
				}
			}
		}
	}
	return nil
}
