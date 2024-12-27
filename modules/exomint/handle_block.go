package exomint

import (
	"fmt"
	"strconv"

	exominttypes "github.com/ExocoreNetwork/exocore/x/exomint/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/callisto/v4/types"
	juno "github.com/forbole/juno/v5/types"
	"github.com/rs/zerolog/log"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	if err := m.saveMintHistory(block.Block.Height, res.BeginBlockEvents); err != nil {
		return fmt.Errorf("error while saving mint history: %s", err)
	}
	return nil
}

// saveEpochStates saves the mint history found in the given events
func (m *Module) saveMintHistory(height int64, events []abci.Event) error {
	log.Debug().Str("module", m.Name()).Int64("height", height).
		Msg("updating mint history")
	events = juno.FindEventsByType(events, exominttypes.EventTypeMint)
	for _, event := range events {
		amountAttr, err := juno.FindAttributeByKey(event, sdk.AttributeKeyAmount)
		if err != nil {
			return fmt.Errorf("error while getting amount: %s", err)
		}
		epochIDAttr, err := juno.FindAttributeByKey(event, exominttypes.AttributeEpochIdentifier)
		if err != nil {
			return fmt.Errorf("error while getting epoch ID: %s", err)
		}
		epochNumberAttr, err := juno.FindAttributeByKey(event, exominttypes.AttributeEpochNumber)
		if err != nil {
			return fmt.Errorf("error while getting epoch number: %s", err)
		}
		epochNumber, err := strconv.ParseInt(epochNumberAttr.Value, 10, 64)
		if err != nil {
			return fmt.Errorf("error while converting epoch number to int: %s", epochNumberAttr.Value)
		}
		denomAttr, err := juno.FindAttributeByKey(event, exominttypes.AttributeDenom)
		if err != nil {
			return fmt.Errorf("error while getting denom: %s", err)
		}
		mintHistory := types.NewMintHistory(
			height, amountAttr.Value, epochIDAttr.Value, epochNumber, denomAttr.Value,
		)
		if err := m.db.AppendMintHistory(mintHistory); err != nil {
			return fmt.Errorf("error while appending mint history: %s", err)
		}
	}
	return nil
}
