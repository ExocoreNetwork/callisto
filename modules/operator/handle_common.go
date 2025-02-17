package operator

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
)

// handleTxAndBeginBlockEvents handles the events emitted during transactions and
// begin blockers. it is not designed to distinguish between the two, and thus,
// should not be used with end blocker events. there are no significant downsides,
// however, beyond 2x work for no reason.
func (m *Module) handleTxAndBeginBlockEvents(events []abci.Event) error {
	// (1) opt out using CLI of x/operator
	// (2) opt out using precompile of x/avs
	// (3) jailed by x/slashing based on transaction
	// (4) jailed by x/slashing based on begin blocker's signing calc
	if err := m.handleOptInfoUpdated(events); err != nil {
		return fmt.Errorf("error while handling opt info updated events: %s", err)
	}
	// (1) tx triggered slashing may cause operator usd value updates
	// although has so far not been implemented
	// (2) jailing by x/slashing based on begin blocker's signing calc
	// (3) epoch change during begin blocker
	if err := m.handleOperatorUSDValues(events); err != nil {
		return fmt.Errorf("error while handling operator usd values: %s", err)
	}
	// deletion of operator usd values can happen
	// (1) upon opt-out
	// (2) during begin blocker if x/avs errors
	if err := m.handleOperatorUSDValueDeletion(events); err != nil {
		return fmt.Errorf("error while handling operator usd values deletion: %s", err)
	}
	// avs usd values can be updated during
	// (1) tx triggered slashing
	// (2) epoch change in begin blocker
	if err := m.handleAvsUSDValues(events); err != nil {
		return fmt.Errorf("error while handling avs usd values: %s", err)
	}
	// deletion of avs usd values can happen if x/avs errors during
	// (1) tx triggered slashing
	// (2) epoch change in begin blocker
	// (3) slashing, whether via tx or begin blocker
	if err := m.handleAvsUSDValueDeletion(events); err != nil {
		return fmt.Errorf("error while handling avs usd values deletion: %s", err)
	}
	return nil
}
