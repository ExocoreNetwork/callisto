package operator

import (
	"fmt"

	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/forbole/callisto/v4/types"
)

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	// the below functions are triggered by transactions
	// register operator using CLI of x/operator
	if err := m.handleOperatorRegistrationEvents(tx.Events); err != nil {
		return fmt.Errorf("error while handling operator registration events: %s", err)
	}
	// (1) opt in using CLI of x/operator
	// (2) opt in using precompile of x/avs
	if err := m.handleOptInEvents(tx.Events); err != nil {
		return fmt.Errorf("error while handling operator opt-in events: %s", err)
	}
	// set consensus key using CLI of x/operator
	if err := m.handleSetConsKey(tx.Events); err != nil {
		return fmt.Errorf("error while handling set cons key events: %s", err)
	}
	// change consensus key using CLI of x/operator
	if err := m.handleSetPrevConsKey(tx.Events); err != nil {
		return fmt.Errorf("error while handling set prev cons key events: %s", err)
	}
	if err := m.handleTxAndBeginBlockEvents(tx.Events); err != nil {
		return fmt.Errorf("error while handling tx begin block events: %s", err)
	}
	// remove consensus key using CLI of x/operator by opting out
	// do this after opt-out events are handled
	if err := m.handleInitConsKeyRemoval(tx.Events); err != nil {
		return fmt.Errorf("error while handling cons key removal events: %s", err)
	}
	return nil
}

// handleOperatorRegistrationEvents handles the events emitted when an operator registers.
func (m *Module) handleOperatorRegistrationEvents(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeRegisterOperator)
	for _, event := range events {
		addr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while getting operator address: %s", err)
		}
		metaInfo, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyMetaInfo)
		if err != nil {
			return fmt.Errorf("error while getting operator meta info: %s", err)
		}
		rate, err := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyCommissionRate)
		if err != nil {
			return fmt.Errorf("error while getting commission rate: %s", err)
		}
		maxCommissionRate, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyMaxCommissionRate)
		if err != nil {
			return fmt.Errorf("error while getting max commission rate: %s", err)
		}
		maxChangeRate, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyMaxChangeRate)
		if err != nil {
			return fmt.Errorf("error while getting max change rate: %s", err)
		}
		lastUpdateTime, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyCommissionUpdateTime)
		if err != nil {
			return fmt.Errorf("error while getting commission update time: %s", err)
		}
		operator := types.NewOperatorFromStr(
			addr.Value, addr.Value, metaInfo.Value,
			rate.Value, maxCommissionRate.Value, maxChangeRate.Value,
			lastUpdateTime.Value,
		)
		err = m.db.SaveOperatorDetail(operator)
		if err != nil {
			return fmt.Errorf("error while saving operator details: %s", err)
		}
	}
	return nil
}

// handleOptInEvents handles the events emitted when an operator opts in.
func (m *Module) handleOptInEvents(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeOptIn)
	for _, event := range events {
		operatorAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while getting operator address: %s", err)
		}
		avsAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyAVSAddr)
		if err != nil {
			return fmt.Errorf("error while getting AVS address: %s", err)
		}
		slashContract, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeySlashContract)
		if err != nil {
			return fmt.Errorf("error while getting slash contract: %s", err)
		}
		optInHeight, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOptInHeight)
		if err != nil {
			return fmt.Errorf("error while getting opt-in height: %s", err)
		}
		// no opt-out attribute is included during opting in, and jailed is false.
		opted := types.NewOptedFromStr(
			operatorAddr.Value, avsAddr.Value, slashContract.Value, optInHeight.Value, "", false,
		)
		err = m.db.SaveOptedState(opted)
		if err != nil {
			return fmt.Errorf("error while saving operator opt state: %s", err)
		}
	}
	return nil
}

// handleSetConsKey handles the events emitted when an operator sets a consensus key.
func (m *Module) handleSetConsKey(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeSetConsKey)
	for _, event := range events {
		addr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while getting operator address: %s", err)
		}
		chainID, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyChainID)
		if err != nil {
			return fmt.Errorf("error while getting chain ID: %s", err)
		}
		consAddress, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyConsensusAddress)
		if err != nil {
			return fmt.Errorf("error while getting consensus address: %s", err)
		}
		consKeyHex, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyConsKeyHex)
		if err != nil {
			return fmt.Errorf("error while getting consensus key hex: %s", err)
		}
		err = m.db.SaveOperatorConsKey(addr.Value, chainID.Value, consKeyHex.Value, consAddress.Value)
		if err != nil {
			return fmt.Errorf("error while saving operator cons key: %s", err)
		}
	}
	return nil
}

// handleSetPrevConsKey handles the events emitted when an operator sets a previous consensus key.
func (m *Module) handleSetPrevConsKey(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeSetPrevConsKey)
	for _, event := range events {
		addr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while getting operator address: %s", err)
		}
		chainID, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyChainID)
		if err != nil {
			return fmt.Errorf("error while getting chain ID: %s", err)
		}
		consAddress, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyConsensusAddress)
		if err != nil {
			return fmt.Errorf("error while getting consensus address: %s", err)
		}
		consKeyHex, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyConsKeyHex)
		if err != nil {
			return fmt.Errorf("error while getting consensus key hex: %s", err)
		}
		err = m.db.SaveOperatorConsKey(addr.Value, chainID.Value, consKeyHex.Value, consAddress.Value)
		if err != nil {
			return fmt.Errorf("error while saving operator cons key: %s", err)
		}
	}
	return nil
}

// handleOptInfoUpdated handles the events emitted when an operator's opted-info is updated.
// Such an update may occur when an operator opts out, or when an operator is jailed. Typically,
// the former happens during a transaction, while the latter happens during a begin/end blocker.
func (m *Module) handleOptInfoUpdated(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeOptInfoUpdated)
	for _, event := range events {
		operatorAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while getting operator address: %s", err)
		}
		avsAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyAVSAddr)
		if err != nil {
			return fmt.Errorf("error while getting AVS address: %s", err)
		}
		slashContract, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeySlashContract)
		if err != nil {
			return fmt.Errorf("error while getting slash contract: %s", err)
		}
		optInHeight, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOptInHeight)
		if err != nil {
			return fmt.Errorf("error while getting opt-in height: %s", err)
		}
		optOutHeight, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOptOutHeight)
		if err != nil {
			return fmt.Errorf("error while getting opt-out height: %s", err)
		}
		jailed, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyJailed)
		if err != nil {
			return fmt.Errorf("error while getting jailed status: %s", err)
		}
		opted := types.NewOptedFromStr(
			operatorAddr.Value, avsAddr.Value, slashContract.Value,
			optInHeight.Value, optOutHeight.Value, jailed.Value == "true",
		)
		err = m.db.SaveOptedState(opted)
		if err != nil {
			return fmt.Errorf("error while saving operator opt state: %s", err)
		}
	}
	return nil
}

// handleConsKeyRemoval handles the events emitted when an operator removes a consensus key.
func (m *Module) handleInitConsKeyRemoval(events []abci.Event) error {
	events = juno.FindEventsByType(events, operatortypes.EventTypeInitRemoveConsKey)
	for _, event := range events {
		operatorAddr, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyOperator)
		if err != nil {
			return fmt.Errorf("error while getting operator address: %s", err)
		}
		chainID, err := juno.FindAttributeByKey(event, operatortypes.AttributeKeyChainID)
		if err != nil {
			return fmt.Errorf("error while getting chain ID: %s", err)
		}
		err = m.db.MarkOperatorKeyRemoval(chainID.Value, operatorAddr.Value)
		if err != nil {
			return fmt.Errorf("error while marking operator key removal: %s", err)
		}
	}
	return nil
}
