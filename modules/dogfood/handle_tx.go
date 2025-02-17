package dogfood

import (
	"fmt"

	dogfoodtypes "github.com/ExocoreNetwork/exocore/x/dogfood/types"
	abci "github.com/cometbft/cometbft/abci/types"
	juno "github.com/forbole/juno/v5/types"

	callistotypes "github.com/forbole/callisto/v4/types"
)

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	if err := m.handleOptOutBeganEvents(tx.Events); err != nil {
		return err
	}
	if err := m.handleConsAddrPruningScheduledEvents(tx.Events); err != nil {
		return err
	}
	if err := m.handleUndelegationMaturityScheduledEvents(tx.Events); err != nil {
		return err
	}
	return nil
}

// handleOptOutBeganEvents handles the events emitted, in response to a tx,
// when an operator opts out of the dogfood AVS.
func (m *Module) handleOptOutBeganEvents(events []abci.Event) error {
	events = juno.FindEventsByType(events, dogfoodtypes.EventTypeOptOutBegan)
	for _, event := range events {
		operatorAddress, err := juno.FindAttributeByKey(event, dogfoodtypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while getting operator address: %s", err)
		}
		epoch, err := juno.FindAttributeByKey(event, dogfoodtypes.AttributeKeyEpoch)
		if err != nil {
			return fmt.Errorf("error while getting epoch: %s", err)
		}
		if err := m.db.SaveOptOutExpiry(
			callistotypes.NewOptOutExpiryFromStr(
				epoch.Value, operatorAddress.Value,
			),
		); err != nil {
			return fmt.Errorf("error while saving opt out expiry: %s", err)
		}
	}
	return nil
}

// handleConsAddrPruningScheduledEvents handles the events emitted, in response to a tx,
// when a consensus address is scheduled to be pruned.
func (m *Module) handleConsAddrPruningScheduledEvents(events []abci.Event) error {
	events = juno.FindEventsByType(events, dogfoodtypes.EventTypeConsAddrPruningScheduled)
	for _, event := range events {
		consAddr, err := juno.FindAttributeByKey(event, dogfoodtypes.AttributeKeyConsAddr)
		if err != nil {
			return fmt.Errorf("error while getting consensus address: %s", err)
		}
		epoch, err := juno.FindAttributeByKey(event, dogfoodtypes.AttributeKeyEpoch)
		if err != nil {
			return fmt.Errorf("error while getting epoch: %s", err)
		}
		if err := m.db.SaveConsensusAddrToPrune(
			callistotypes.NewConsensusAddrToPruneFromStr(
				epoch.Value, consAddr.Value,
			),
		); err != nil {
			return fmt.Errorf("error while saving consensus addr to prune: %s", err)
		}
	}
	return nil
}

// handleUndelegationMaturityScheduledEvents handles the events emitted, in response to a tx,
// when an undelegation maturity is scheduled.
func (m *Module) handleUndelegationMaturityScheduledEvents(events []abci.Event) error {
	events = juno.FindEventsByType(events, dogfoodtypes.EventTypeUndelegationMaturityScheduled)
	for _, event := range events {
		recordKey, err := juno.FindAttributeByKey(event, dogfoodtypes.AttributeKeyRecordID)
		if err != nil {
			return fmt.Errorf("error while getting record key: %s", err)
		}
		epoch, err := juno.FindAttributeByKey(event, dogfoodtypes.AttributeKeyEpoch)
		if err != nil {
			return fmt.Errorf("error while getting epoch: %s", err)
		}
		if err := m.db.SaveUndelegationMaturity(
			callistotypes.NewUndelegationMaturityFromStr(
				epoch.Value, recordKey.Value,
			),
		); err != nil {
			return fmt.Errorf("error while saving undelegation maturity: %s", err)
		}
	}
	return nil
}
