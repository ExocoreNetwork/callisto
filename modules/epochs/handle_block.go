package epochs

import (
	"fmt"

	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	juno "github.com/forbole/juno/v5/types"
	"github.com/rs/zerolog/log"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	if err := m.saveEpochStates(block.Block.Height, res.BeginBlockEvents); err != nil {
		return fmt.Errorf("error while saving epoch states: %s", err)
	}
	return nil
}

// saveEpochStates saves the epoch states found in the given events
func (m *Module) saveEpochStates(height int64, events []abci.Event) error {
	log.Debug().Str("module", m.Name()).Int64("height", height).
		Msg("updating epoch states")
	// if an epoch end event is found, save the epoch state.
	events = juno.FindEventsByType(events, epochstypes.EventTypeEpochEnd)
	for _, event := range events {
		epochID, err := juno.FindAttributeByKey(event, epochstypes.AttributeEpochIdentifier)
		if err != nil {
			return fmt.Errorf("error while getting epoch ID: %s", err)
		}
		epoch, err := m.source.GetEpochInfo(height, epochID.Value)
		if err != nil {
			return fmt.Errorf("error while getting epoch info: %s", err)
		}
		if err = m.db.SaveEpochState(epoch); err != nil {
			return fmt.Errorf("error while saving epoch state: %s", err)
		}
	}
	return nil
}
