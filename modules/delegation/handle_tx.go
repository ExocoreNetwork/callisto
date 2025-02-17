package delegation

import (
	"fmt"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/forbole/callisto/v4/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	// tx driven delegation or undelegation
	if err := m.handleDelegationStateUpdates(tx.Events); err != nil {
		return fmt.Errorf("error while handling delegation state updates: %s", err)
	}
	// tx driven staker operator association
	if err := m.handleStakerOperatorAssociations(tx.Events); err != nil {
		return fmt.Errorf("error while handling staker operator association updates: %s", err)
	}
	// tx driven staker operator disassociation
	if err := m.handleStakerOperatorDisassociations(tx.Events); err != nil {
		return fmt.Errorf("error while handling staker operator disassociation updates: %s", err)
	}
	// tx driven delegation of assetID by stakerID to operatorAddr
	if err := m.handleStakerAppendedToOperatorAsset(tx.Events); err != nil {
		return fmt.Errorf("error while handling staker appended to operator updates: %s", err)
	}
	// tx driven undelegation or slashing
	if err := m.handleStakerRemovedFromOperatorAsset(tx.Events); err != nil {
		return fmt.Errorf("error while handling staker removed from operator asset updates: %s", err)
	}
	// tx driven undelegation or slashing
	if err := m.handleAllStakersRemovedFromOperatorAsset(tx.Events); err != nil {
		return fmt.Errorf("error while handling all stakers removed from operator asset updates: %s", err)
	}
	// tx driven exo asset delegation
	if err := m.handleExoAssetDelegations(tx.Events); err != nil {
		return fmt.Errorf("error while handling exo asset delegations: %s", err)
	}
	// tx driven undelegation starts
	if err := m.handleUndelegationStarts(tx.Events); err != nil {
		return fmt.Errorf("error while handling undelegation starts: %s", err)
	}
	// tx driven undelegation hold count changes via hooks
	if err := m.handleUndelegationHoldCountChanges(tx.Events); err != nil {
		return fmt.Errorf("error while handling undelegation hold count changes: %s", err)
	}
	// slashing of undelegations via txs (not yet implemented)
	if err := m.handleUndelegationSlashings(tx.Events); err != nil {
		return fmt.Errorf("error while handling undelegation slashings: %s", err)
	}
	return nil
}

// handleDelegationStateUpdates handles the delegation state updates.
// the events are emitted in response to transactions: delegate, undelegate, slashing, nst change
// and in response to EndBlock for undelegation maturity
// and in response to BeginBlock for slashing.
func (m *Module) handleDelegationStateUpdates(events []abci.Event) error {
	events = juno.FindEventsByType(events, delegationtypes.EventTypeDelegationStateUpdated)
	for _, event := range events {
		stakerID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyStakerID)
		if err != nil {
			return fmt.Errorf("error while finding staker ID: %s", err)
		}
		operatorAddr, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyOperatorAddr)
		if err != nil {
			return fmt.Errorf("error while finding operator address: %s", err)
		}
		assetID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while finding asset ID: %s", err)
		}
		waitUndelegationAmount, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyWaitUndelegationAmount)
		if err != nil {
			return fmt.Errorf("error while finding wait undelegation amount: %s", err)
		}
		undelegatableShare, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyUndelegatableShare)
		if err != nil {
			return fmt.Errorf("error while finding undelegatable share: %s", err)
		}
		delegationState := types.NewDelegationStateFromStr(
			stakerID.Value, assetID.Value, operatorAddr.Value,
			undelegatableShare.Value, waitUndelegationAmount.Value,
		)
		if err := m.db.SaveDelegationState(delegationState); err != nil {
			return fmt.Errorf("error while saving delegation state: %s", err)
		}
	}
	return nil
}

// handleStakerOperatorAssociations handles the staker operator associations.
// only triggered by transactions
func (m *Module) handleStakerOperatorAssociations(events []abci.Event) error {
	events = juno.FindEventsByType(events, delegationtypes.EventTypeOperatorAssociated)
	for _, event := range events {
		stakerID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyStakerID)
		if err != nil {
			return fmt.Errorf("error while finding staker ID: %s", err)
		}
		operatorAddr, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyOperatorAddr)
		if err != nil {
			return fmt.Errorf("error while finding operator address: %s", err)
		}
		if err := m.db.SaveStakerOperatorAssociation(stakerID.Value, operatorAddr.Value); err != nil {
			return fmt.Errorf("error while saving staker operator association: %s", err)
		}
	}
	return nil
}

// handleStakerOperatorDisassociations handles the staker operator disassociations.
// only triggered by transactions
func (m *Module) handleStakerOperatorDisassociations(events []abci.Event) error {
	events = juno.FindEventsByType(events, delegationtypes.EventTypeOperatorDisassociated)
	for _, event := range events {
		stakerID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyStakerID)
		if err != nil {
			return fmt.Errorf("error while finding staker ID: %s", err)
		}
		if err := m.db.DeleteStakerOperatorAssociation(stakerID.Value); err != nil {
			return fmt.Errorf("error while deleting staker operator association: %s", err)
		}
	}
	return nil
}

// handleStakerAppendedToOperatorAsset handles the staker appended to operator asset updates.
// only triggered by transactions
func (m *Module) handleStakerAppendedToOperatorAsset(events []abci.Event) error {
	events = juno.FindEventsByType(events, delegationtypes.EventTypeStakerAppended)
	for _, event := range events {
		stakerID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyStakerID)
		if err != nil {
			return fmt.Errorf("error while finding staker ID: %s", err)
		}
		operatorAddr, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyOperatorAddr)
		if err != nil {
			return fmt.Errorf("error while finding operator address: %s", err)
		}
		assetID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while finding asset ID: %s", err)
		}
		if err := m.db.AppendStakerToOperatorAsset(stakerID.Value, operatorAddr.Value, assetID.Value); err != nil {
			return fmt.Errorf("error while appending staker to operator asset: %s", err)
		}
	}
	return nil
}

// handleStakerRemovedFromOperatorAsset handles the staker removed from operator asset updates.
// only triggered by transactions
func (m *Module) handleStakerRemovedFromOperatorAsset(events []abci.Event) error {
	events = juno.FindEventsByType(events, delegationtypes.EventTypeStakerRemoved)
	for _, event := range events {
		stakerID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyStakerID)
		if err != nil {
			return fmt.Errorf("error while finding staker ID: %s", err)
		}
		operatorAddr, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyOperatorAddr)
		if err != nil {
			return fmt.Errorf("error while finding operator address: %s", err)
		}
		assetID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while finding asset ID: %s", err)
		}
		if err := m.db.RemoveStakerFromOperatorAsset(stakerID.Value, operatorAddr.Value, assetID.Value); err != nil {
			return fmt.Errorf("error while removing staker from operator asset: %s", err)
		}
	}
	return nil
}

// handleAllStakersRemovedFromOperatorAsset handles the all stakers removed from operator asset updates.
// triggered in response to slashing, which may be via txs or in BeginBlock
func (m *Module) handleAllStakersRemovedFromOperatorAsset(events []abci.Event) error {
	events = juno.FindEventsByType(events, delegationtypes.EventTypeAllStakersRemoved)
	for _, event := range events {
		operatorAddr, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyOperatorAddr)
		if err != nil {
			return fmt.Errorf("error while finding operator address: %s", err)
		}
		assetID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while finding asset ID: %s", err)
		}
		if err := m.db.DeleteAllStakersFromOperatorAsset(operatorAddr.Value, assetID.Value); err != nil {
			return fmt.Errorf("error while deleting all stakers from operator asset: %s", err)
		}
	}
	return nil
}

// handleExoAssetDelegations handles the exo asset delegations.
// only triggered by transactions
func (m *Module) handleExoAssetDelegations(events []abci.Event) error {
	events = juno.FindEventsByType(events, delegationtypes.EventTypeExoAssetDelegation)
	for _, event := range events {
		stakerID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyStakerID)
		if err != nil {
			return fmt.Errorf("error while finding staker ID: %s", err)
		}
		operatorAddr, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyOperatorAddr)
		if err != nil {
			return fmt.Errorf("error while finding operator address: %s", err)
		}
		amount, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyAmount)
		if err != nil {
			return fmt.Errorf("error while finding amount: %s", err)
		}
		delegation := types.NewExoAssetDelegationFromStr(
			stakerID.Value, operatorAddr.Value, amount.Value,
			"0", "0",
		)
		if err := m.db.AccumulateExoAssetDelegation(delegation); err != nil {
			return fmt.Errorf("error while accumulating exo asset delegation: %s", err)
		}
	}
	return nil
}

// handleUndelegationStarts handles the undelegation starts.
// only triggered by transactions
func (m *Module) handleUndelegationStarts(events []abci.Event) error {
	events = juno.FindEventsByType(events, delegationtypes.EventTypeUndelegationStarted)
	for _, event := range events {
		stakerID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyStakerID)
		if err != nil {
			return fmt.Errorf("error while finding staker ID: %s", err)
		}
		assetID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while finding asset ID: %s", err)
		}
		operatorAddr, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyOperatorAddr)
		if err != nil {
			return fmt.Errorf("error while finding operator address: %s", err)
		}
		recordID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyRecordID)
		if err != nil {
			return fmt.Errorf("error while finding record ID: %s", err)
		}
		amount, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyAmount)
		if err != nil {
			return fmt.Errorf("error while finding amount: %s", err)
		}
		completedEpochIdentifier, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyCompletedEpochID)
		if err != nil {
			return fmt.Errorf("error while finding completed epoch ID: %s", err)
		}
		completedEpochNumber, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyCompletedEpochNumber)
		if err != nil {
			return fmt.Errorf("error while finding completed epoch number: %s", err)
		}
		undelegationID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyUndelegationID)
		if err != nil {
			return fmt.Errorf("error while finding undelegation ID: %s", err)
		}
		txHash, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyTxHash)
		if err != nil {
			return fmt.Errorf("error while finding tx hash: %s", err)
		}
		blockNumber, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyBlockNumber)
		if err != nil {
			return fmt.Errorf("error while finding block number: %s", err)
		}
		undelegation := types.NewUndelegationRecordFromStr(
			recordID.Value, stakerID.Value, assetID.Value, operatorAddr.Value,
			txHash.Value, blockNumber.Value, completedEpochIdentifier.Value, completedEpochNumber.Value,
			undelegationID.Value, amount.Value, amount.Value, "0", /* holdCount */
		)
		if err := m.db.SaveUndelegationRecord(undelegation); err != nil {
			return fmt.Errorf("error while saving undelegation record: %s", err)
		}
		// if there is an exo-asset undelegation, figure out what to do
		// (1) pending_undelegation += amount
		// (2) delegated -= amount
		if assetID.Value == assetstypes.ExocoreAssetID {
			if err := m.db.UndelegateExoAsset(stakerID.Value, operatorAddr.Value, amount.Value); err != nil {
				return fmt.Errorf("error while undelegating exo asset: %s", err)
			}
		}
	}
	return nil
}

// handleUndelegationHoldCountChanges handles the undelegation hold count changes.
// triggered by transactions (undelegation => held by AVS) and
// during EndBlock (released upon unbonding epoch duration end)
func (m *Module) handleUndelegationHoldCountChanges(events []abci.Event) error {
	events = juno.FindEventsByType(events, delegationtypes.EventTypeUndelegationHoldCountChanged)
	for _, event := range events {
		recordID, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyRecordID)
		if err != nil {
			return fmt.Errorf("error while finding record ID: %s", err)
		}
		holdCount, err := juno.FindAttributeByKey(event, delegationtypes.AttributeKeyHoldCount)
		if err != nil {
			return fmt.Errorf("error while finding hold count: %s", err)
		}
		if err := m.db.UpdateUndelegationRecordHoldCount(recordID.Value, holdCount.Value); err != nil {
			return fmt.Errorf("error while updating undelegation hold count: %s", err)
		}
	}
	return nil
}
