package delegation

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	juno "github.com/forbole/juno/v5/types"

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
	// tracked in x/asset under `operator_assets`
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
	}
	return nil
}

// handleUndelegationSlashings handles the slashing of undelegation records. Technically, this
// happens in x/operator, however, it modifies the state of x/delegation schema items. So, I
// put it here.
func (m *Module) handleUndelegationSlashings(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeUndelegationSlashed)
	for _, event := range events {
		recordID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyRecordID)
		if err != nil {
			return fmt.Errorf("error while finding record ID: %s", err)
		}
		amount, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyAmount)
		if err != nil {
			return fmt.Errorf("error while finding amount: %s", err)
		}
		if err := m.db.SlashUndelegationRecord(recordID.Value, amount.Value); err != nil {
			return fmt.Errorf("error while slashing undelegation record: %s", err)
		}
	}
	return nil
}
