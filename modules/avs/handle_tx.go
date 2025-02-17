package avs

import (
	"fmt"

	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	abci "github.com/cometbft/cometbft/abci/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	if err := m.handleAvsCreatedEvents(tx.Events); err != nil {
		return err
	}
	if err := m.handleChainAvsCreatedEvents(tx.Events); err != nil {
		return err
	}
	return nil
}

// handleAvsCreatedEvents handles the events emitted when an AVS is created
func (m *Module) handleAvsCreatedEvents(events []abci.Event) error {
	events = juno.FindEventsByType(events, avstypes.EventTypeAvsCreated)
	for _, event := range events {
		avsAddress, err := juno.FindAttributeByKey(event, avstypes.AttributeKeyAvsAddress)
		if err != nil {
			return fmt.Errorf("error while getting avs address: %s", err)
		}
		if err := m.db.SaveAvsAddr(avsAddress.Value); err != nil {
			return fmt.Errorf("error while saving avs address: %s", err)
		}
	}
	return nil
}

// handleChainAvsCreatedEvents handles the events emitted when a chain AVS is created
func (m *Module) handleChainAvsCreatedEvents(events []abci.Event) error {
	events = juno.FindEventsByType(events, avstypes.EventTypeChainAvsCreated)
	for _, event := range events {
		chainID, err := juno.FindAttributeByKey(event, avstypes.AttributeKeyChainID)
		if err != nil {
			return fmt.Errorf("error while getting chain id: %s", err)
		}
		avsAddress, err := juno.FindAttributeByKey(event, avstypes.AttributeKeyAvsAddress)
		if err != nil {
			return fmt.Errorf("error while getting avs address: %s", err)
		}
		if err := m.db.SaveChainIdToAvsAddr(chainID.Value, avsAddress.Value); err != nil {
			return fmt.Errorf("error while saving chain id to avs address: %s", err)
		}
	}
	return nil
}
