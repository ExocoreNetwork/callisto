package operator

import (
	"fmt"

	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/forbole/callisto/v4/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	if err := m.handleTxAndBeginBlockEvents(res.BeginBlockEvents); err != nil {
		return fmt.Errorf("error while handling tx and begin block events: %s", err)
	}
	// 2 known end blocker events for x/operator: operator key removal completed
	// it is triggered via x/dogfood's EndBlocker
	if err := m.handleCompleteConsKeyRemoval(res.EndBlockEvents); err != nil {
		return fmt.Errorf("error while handling complete cons key removal: %s", err)
	}
	// previous consensus key removal, again, via x/dogfood's EndBlocker
	if err := m.handleClearOperatorPrevConsKey(res.EndBlockEvents); err != nil {
		return fmt.Errorf("error while handling clear operator prev cons key: %s", err)
	}
	return nil
}

// handleOperatorUSDValues handles the events emitted when an operator's usd values are updated.
func (m *Module) handleOperatorUSDValues(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeUpdateOperatorUSDValue)
	for _, event := range events {
		operatorAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while getting operator address: %s", err)
		}
		avsAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyAVSAddr)
		if err != nil {
			return fmt.Errorf("error while getting avs address: %s", err)
		}
		selfUsdValue, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeySelfUSDValue)
		if err != nil {
			return fmt.Errorf("error while getting self usd value: %s", err)
		}
		totalUsdValue, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyTotalUSDValue)
		if err != nil {
			return fmt.Errorf("error while getting total usd value: %s", err)
		}
		activeUsdValue, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyActiveUSDValue)
		if err != nil {
			return fmt.Errorf("error while getting active usd value: %s", err)
		}
		operator := types.NewOperatorUSDValueFromStr(
			operatorAddr.Value, avsAddr.Value,
			selfUsdValue.Value, totalUsdValue.Value, activeUsdValue.Value,
		)
		if err := m.db.SaveOperatorUSDValue(operator); err != nil {
			return fmt.Errorf("error while saving operator usd value: %s", err)
		}
	}
	return nil
}

// handleAvsUSDValues handles the events emitted when an avs's usd values are updated.
func (m *Module) handleAvsUSDValues(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeUpdateAVSUSDValue)
	for _, event := range events {
		avsAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyAVSAddr)
		if err != nil {
			return fmt.Errorf("error while getting avs address: %s", err)
		}
		usdValue, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyTotalUSDValue)
		if err != nil {
			return fmt.Errorf("error while getting usd value: %s", err)
		}
		avs := types.NewAvsUSDValueFromStr(avsAddr.Value, usdValue.Value)
		if err := m.db.SaveAvsUSDValue(avs); err != nil {
			return fmt.Errorf("error while saving avs usd value: %s", err)
		}
	}
	return nil
}

// handleOperatorUSDValueDeletion handles the events emitted when an operator's usd values are deleted.
func (m *Module) handleOperatorUSDValueDeletion(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeDeleteOperatorUSDValue)
	for _, event := range events {
		operatorAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while getting operator address: %s", err)
		}
		avsAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyAVSAddr)
		if err != nil {
			return fmt.Errorf("error while getting avs address: %s", err)
		}
		// no need for a dedicated type for this situation when 2 strings can suffice
		if err := m.db.DeleteOperatorUSDValue(operatorAddr.Value, avsAddr.Value); err != nil {
			return fmt.Errorf("error while deleting operator usd value: %s", err)
		}
	}
	return nil
}

// handleAvsUSDValueDeletion handles the events emitted when an avs's usd values are deleted.
func (m *Module) handleAvsUSDValueDeletion(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeDeleteAVSUSDValue)
	for _, event := range events {
		avsAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyAVSAddr)
		if err != nil {
			return fmt.Errorf("error while getting avs address: %s", err)
		}
		if err := m.db.DeleteAvsUSDValue(avsAddr.Value); err != nil {
			return fmt.Errorf("error while deleting avs usd value: %s", err)
		}
	}
	return nil
}

// handleCompleteConsKeyRemoval handles the events emitted when an operator completes the removal of a consensus key.
func (m *Module) handleCompleteConsKeyRemoval(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeEndRemoveConsKey)
	for _, event := range events {
		operatorAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while getting operator address: %s", err)
		}
		chainID, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyChainID)
		if err != nil {
			return fmt.Errorf("error while getting chain ID: %s", err)
		}
		err = m.db.RemoveOperatorConsKey(chainID.Value, operatorAddr.Value)
		if err != nil {
			return fmt.Errorf("error while removing operator cons key: %s", err)
		}
	}
	return nil
}

// handleClearOperatorPrevConsKey handles the events emitted when an operator's previous consensus key is cleared.
func (m *Module) handleClearOperatorPrevConsKey(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeRemovePrevConsKey)
	for _, event := range events {
		operatorAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while getting operator address: %s", err)
		}
		chainID, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyChainID)
		if err != nil {
			return fmt.Errorf("error while getting chain ID: %s", err)
		}
		if err := m.db.ClearOperatorPrevConsKey(operatorAddr.Value, chainID.Value); err != nil {
			return fmt.Errorf("error while clearing operator prev cons key: %s", err)
		}
	}
	return nil
}
