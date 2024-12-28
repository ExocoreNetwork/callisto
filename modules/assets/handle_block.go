package assets

import (
	"fmt"

	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	// staker events are emitted during EndBlock for undelegation maturity
	// and in response to transactions (handle_tx.go)
	if err := m.handleStakerEvents(block.Block.Height, res.EndBlockEvents); err != nil {
		return fmt.Errorf("error while handling staker events: %s", err)
	}
	// similarly, operator events are emitted during EndBlock for undelegation maturity
	// and in response to transactions (handle_tx.go)
	if err := m.handleOperatorEvents(block.Block.Height, res.EndBlockEvents); err != nil {
		return fmt.Errorf("error while handling operator events: %s", err)
	}
	return nil
}
