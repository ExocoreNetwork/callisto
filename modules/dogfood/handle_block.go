package dogfood

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dogfoodtypes "github.com/ExocoreNetwork/exocore/x/dogfood/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	juno "github.com/forbole/juno/v5/types"

	callistotypes "github.com/forbole/callisto/v4/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	// validator set can only change at the end of a block, so can the voting power.
	if err := m.handleLastTotalPowerUpdated(res.EndBlockEvents); err != nil {
		return fmt.Errorf("error while handling dogfood last total power updated: %s", err)
	}
	// handle validator set change
	if err := m.handleValidatorSetChange(block.Block.Height, res.EndBlockEvents); err != nil {
		return fmt.Errorf("error while handling dogfood validator set change: %s", err)
	}
	// the following items are handled in BeginBlock in response to epoch hooks
	// technically, they are "upgraded" during the BeginBlock, to be handled within the EndBlock
	// however, that is only a 3.5s difference, so we can just handle them here. the height,
	// is, of course, the same in either begin or end block.
	if err := m.handleOptOutsFinished(block.Block.Height, res.BeginBlockEvents); err != nil {
		return fmt.Errorf("error while handling dogfood opt outs finished: %s", err)
	}
	if err := m.handleConsAddrsPruned(block.Block.Height, res.BeginBlockEvents); err != nil {
		return fmt.Errorf("error while handling dogfood consensus addr pruned: %s", err)
	}
	if err := m.handleUndelegationsMatured(block.Block.Height, res.BeginBlockEvents); err != nil {
		return fmt.Errorf("error while handling dogfood undelegation matured: %s", err)
	}
	return nil
}

// handleLastTotalPowerUpdated handles the event emitted when the last total power is updated.
func (m *Module) handleLastTotalPowerUpdated(events []abci.Event) error {
	events = juno.FindEventsByType(events, dogfoodtypes.EventTypeLastTotalPowerUpdated)
	for _, event := range events {
		// there is only one attribute
		m.db.SaveLastTotalPower(event.Attributes[0].Value)
	}
	return nil
}

// handleValidatorSetChange handles the event emitted when the validator set changes.
func (m *Module) handleValidatorSetChange(height int64, events []abci.Event) error {
	events = juno.FindEventsByType(events, dogfoodtypes.EventTypeLastTotalPowerUpdated)
	if len(events) == 0 {
		return nil
	}
	// now we need to get the validator set from the source
	validators, err := m.source.GetValidators(height)
	if err != nil {
		return fmt.Errorf("error while getting validators: %s", err)
	}
	// remember that, this function is called after worker.SaveValidators, so we do not
	// need to save the validators - we can instead focus on the vote power.
	votingPowers := make([]callistotypes.ValidatorVotingPower, len(validators))
	for i, validator := range validators {
		votingPowers[i] = callistotypes.NewValidatorVotingPower(
			sdk.ConsAddress(validator.Address).String(),
			validator.Power, height,
		)
	}
	if err := m.db.SaveValidatorsVotingPowers(votingPowers); err != nil {
		return fmt.Errorf("error while saving validator voting powers: %s", err)
	}
	return nil
}

// handleOptOutsFinished handles the event emitted when an opt out is finished.
func (m *Module) handleOptOutsFinished(height int64, events []abci.Event) error {
	events = juno.FindEventsByType(events, dogfoodtypes.EventTypeOptOutsFinished)
	for _, event := range events {
		epoch, err := juno.FindAttributeByKey(event, dogfoodtypes.AttributeKeyEpoch)
		if err != nil {
			return fmt.Errorf("error while getting epoch: %s", err)
		}
		if err := m.db.CompleteOptOuts(epoch.Value, height); err != nil {
			return fmt.Errorf("error while completing opt outs: %s", err)
		}
	}
	return nil
}

// handleConsAddrsPruned handles the event emitted when a consensus address is pruned.
func (m *Module) handleConsAddrsPruned(height int64, events []abci.Event) error {
	events = juno.FindEventsByType(events, dogfoodtypes.EventTypeConsAddrsPruned)
	for _, event := range events {
		epoch, err := juno.FindAttributeByKey(event, dogfoodtypes.AttributeKeyEpoch)
		if err != nil {
			return fmt.Errorf("error while getting epoch: %s", err)
		}
		if err := m.db.CompleteConsensusAddrsPruning(epoch.Value, height); err != nil {
			return fmt.Errorf("error while completing consensus addrs pruning: %s", err)
		}
	}
	return nil
}

// handleUndelegationsMatured handles the event emitted when an undelegation is matured.
func (m *Module) handleUndelegationsMatured(height int64, events []abci.Event) error {
	events = juno.FindEventsByType(events, dogfoodtypes.EventTypeUndelegationsMatured)
	for _, event := range events {
		epoch, err := juno.FindAttributeByKey(event, dogfoodtypes.AttributeKeyEpoch)
		if err != nil {
			return fmt.Errorf("error while getting epoch: %s", err)
		}
		if err := m.db.MatureUndelegations(
			epoch.Value, height,
		); err != nil {
			return fmt.Errorf("error while maturing undelegations: %s", err)
		}
	}
	return nil
}
