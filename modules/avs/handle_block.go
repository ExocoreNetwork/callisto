package avs

import (
	"fmt"

	dogfoodtypes "github.com/ExocoreNetwork/exocore/x/dogfood/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	if err := m.handleDogfoodAvsCreationEvents(res.BeginBlockEvents); err != nil {
		return fmt.Errorf("error while handling dogfood avs creation events: %s", err)
	}
	return nil
}

// handleDogfoodAvsCreationEvents handles the events emitted when a dogfood AVS is created.
func (m *Module) handleDogfoodAvsCreationEvents(events []abci.Event) error {
	events = juno.FindEventsByType(events, dogfoodtypes.EventTypeDogfoodAvsCreated)
	for _, event := range events {
		avsAddress, err := juno.FindAttributeByKey(event, dogfoodtypes.AttributeKeyAvsAddress)
		if err != nil {
			return fmt.Errorf("error while getting avs address: %s", err)
		}
		chainID, err := juno.FindAttributeByKey(event, dogfoodtypes.AttributeKeyChainIDWithoutRev)
		if err != nil {
			return fmt.Errorf("error while getting chain id: %s", err)
		}
		// handle avs creation
		if err := m.db.SaveAvsAddr(avsAddress.Value); err != nil {
			return fmt.Errorf("error while saving avs address: %s", err)
		}
		// and that it was a chain avs
		if err := m.db.SaveChainIdToAvsAddr(chainID.Value, avsAddress.Value); err != nil {
			return fmt.Errorf("error while saving chain id to avs address: %s", err)
		}
	}
	return nil
}
