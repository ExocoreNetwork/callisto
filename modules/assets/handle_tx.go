package assets

import (
	"fmt"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/forbole/callisto/v4/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	// SetClientChainInfo is only called during a transaction and not in End/BeginBlock, so
	// it only exists here.
	if err := m.handleClientChainEvents(tx.Events); err != nil {
		return fmt.Errorf("error while handling client chain events: %s", err)
	}
	// same logic as above; only triggered by tx
	if err := m.handleNewTokenEvents(tx.Events); err != nil {
		return fmt.Errorf("error while handling new token events: %s", err)
	}
	if err := m.handleUpdateTokenEvents(tx.Events); err != nil {
		return fmt.Errorf("error while handling updated token events: %s", err)
	}
	if err := m.handleUpdateStakingTotalAmount(tx.Events); err != nil {
		return fmt.Errorf("error while handling updated staking total amount: %s", err)
	}
	// the same event (as HandleBlock) is emitted in response to transactions
	// for example, if a staker delegates to an operator or deposits an asset.
	if err := m.handleStakerEvents(tx.Height, tx.Events); err != nil {
		return fmt.Errorf("error while handling staker events: %s", err)
	}
	// handle operator-level information changing such as operator shares
	if err := m.handleOperatorEvents(tx.Height, tx.Events); err != nil {
		return fmt.Errorf("error while handling operator events: %s", err)
	}
	return nil
}

// handleClientChainEvents filters, parses and indexes the client chain events.
func (m *Module) handleClientChainEvents(events []abci.Event) error {
	if err := m.handleClientChainEventsByType(events, assetstypes.EventTypeNewClientChain); err != nil {
		return fmt.Errorf("error while handling new client chain events: %s", err)
	}
	if err := m.handleClientChainEventsByType(events, assetstypes.EventTypeUpdatedClientChain); err != nil {
		return fmt.Errorf("error while handling updated client chain events: %s", err)
	}
	return nil
}

// handleClientChainEventsByType filters, parses and indexes the client chain events of the provided type.
func (m *Module) handleClientChainEventsByType(events []abci.Event, ty string) error {
	events = juno.FindEventsByType(events, ty)
	for _, event := range events {
		name, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyName)
		if err != nil {
			return fmt.Errorf("error while getting client chain name: %s", err)
		}
		metaInfo, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyMetaInfo)
		if err != nil {
			return fmt.Errorf("error while getting client chain meta info: %s", err)
		}
		chainId, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyChainID)
		if err != nil {
			return fmt.Errorf("error while getting client chain ID: %s", err)
		}
		exocoreChainIdx, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyExocoreChainIdx)
		if err != nil {
			return fmt.Errorf("error while getting exocore chain index: %s", err)
		}
		finalizationBlocks, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyFinalizationBlocks)
		if err != nil {
			return fmt.Errorf("error while getting finalization blocks: %s", err)
		}
		lzID, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyLZID)
		if err != nil {
			return fmt.Errorf("error while getting lzID: %s", err)
		}
		sigType, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeySigType)
		if err != nil {
			return fmt.Errorf("error while getting signature type: %s", err)
		}
		addrLength, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyAddrLength)
		if err != nil {
			return fmt.Errorf("error while getting address length: %s", err)
		}
		chain := types.NewClientChainFromStr(
			name.Value, metaInfo.Value, chainId.Value,
			exocoreChainIdx.Value, finalizationBlocks.Value,
			lzID.Value, sigType.Value, addrLength.Value,
		)
		if err := m.db.SaveOrUpdateClientChain(chain); err != nil {
			return fmt.Errorf("error while saving client chain: %s", err)
		}
	}
	return nil
}

// handleNewTokenEvents filters, parses and indexes the new token events.
func (m *Module) handleNewTokenEvents(events []abci.Event) error {
	events = juno.FindEventsByType(events, assetstypes.EventTypeNewToken)
	for _, event := range events {
		assetID, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while getting asset ID: %s", err)
		}
		name, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyName)
		if err != nil {
			return fmt.Errorf("error while getting token name: %s", err)
		}
		symbol, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeySymbol)
		if err != nil {
			return fmt.Errorf("error while getting token symbol: %s", err)
		}
		address, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyAddress)
		if err != nil {
			return fmt.Errorf("error while getting token address: %s", err)
		}
		decimals, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyDecimals)
		if err != nil {
			return fmt.Errorf("error while getting token decimals: %s", err)
		}
		lzID, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyLZID)
		if err != nil {
			return fmt.Errorf("error while getting lzID: %s", err)
		}
		metaInfo, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyMetaInfo)
		if err != nil {
			return fmt.Errorf("error while getting token meta info: %s", err)
		}
		exocoreChainIdx, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyExocoreChainIdx)
		if err != nil {
			return fmt.Errorf("error while getting exocore chain index: %s", err)
		}
		stakingTotalAmount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyTotalAmount)
		if err != nil {
			return fmt.Errorf("error while getting staking total amount: %s", err)
		}
		token := types.NewAssetsTokenFromStr(
			assetID.Value, name.Value, symbol.Value, address.Value, decimals.Value,
			lzID.Value, metaInfo.Value, exocoreChainIdx.Value, stakingTotalAmount.Value,
		)
		if err := m.db.SaveAssetsToken(token); err != nil {
			return fmt.Errorf("error while saving token: %s", err)
		}
	}
	return nil
}

// handleUpdateTokenEvents filters, parses and indexes the updated token events.
func (m *Module) handleUpdateTokenEvents(events []abci.Event) error {
	events = juno.FindEventsByType(events, assetstypes.EventTypeUpdatedToken)
	for _, event := range events {
		assetID, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while getting asset ID: %s", err)
		}
		// only metaInfo can be changed after token registration
		metaInfo, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyMetaInfo)
		if err != nil {
			return fmt.Errorf("error while getting token meta info: %s", err)
		}
		if err := m.db.UpdateAssetMetadata(assetID.Value, metaInfo.Value); err != nil {
			return fmt.Errorf("error while updating token metadata: %s", err)
		}
	}
	return nil
}

// handleUpdateStakingTotalAmount filters, parses and indexes the updated staking total amount
// events.
func (m *Module) handleUpdateStakingTotalAmount(events []abci.Event) error {
	events = juno.FindEventsByType(events, assetstypes.EventTypeUpdatedStakingTotalAmount)
	for _, event := range events {
		amount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyTotalAmount)
		if err != nil {
			return fmt.Errorf("error while getting staking total amount: %s", err)
		}
		assetID, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while getting asset ID: %s", err)
		}
		if err := m.db.UpdateStakingTotalAmount(assetID.Value, amount.Value); err != nil {
			return fmt.Errorf("error while updating staking total amount: %s", err)
		}
	}
	return nil
}

// handleStakerEvents filters, parses and indexes the staker events.
func (m *Module) handleStakerEvents(height int64, events []abci.Event) error {
	events = juno.FindEventsByType(events, assetstypes.EventTypeUpdatedStakerAsset)
	for _, event := range events {
		stakerID, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyStakerID)
		if err != nil {
			return fmt.Errorf("error while getting staker ID: %s", err)
		}
		assetID, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while getting asset ID: %s", err)
		}
		depositAmount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyDepositAmount)
		if err != nil {
			return fmt.Errorf("error while getting deposit amount: %s", err)
		}
		withdrawableAmount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyWithdrawableAmount)
		if err != nil {
			return fmt.Errorf("error while getting withdrawable amount: %s", err)
		}
		pendingUndelegationAmount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyPendingUndelegationAmount)
		if err != nil {
			return fmt.Errorf("error while getting pending undelegation amount: %s", err)
		}
		asset := types.NewStakerAssetFromStr(
			stakerID.Value, assetID.Value,
			depositAmount.Value, withdrawableAmount.Value, pendingUndelegationAmount.Value,
			"", // no update to slashed amount
		)
		if err := m.db.SaveStakerAsset(asset); err != nil {
			return fmt.Errorf("error while saving staker asset: %s", err)
		}
	}
	return nil
}

// handleOperatorEvents filters, parses and indexes the operator events.
func (m *Module) handleOperatorEvents(height int64, events []abci.Event) error {
	events = juno.FindEventsByType(events, assetstypes.EventTypeUpdatedOperatorAsset)
	for _, event := range events {
		operatorAddr, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyOperatorAddress)
		if err != nil {
			return fmt.Errorf("error while getting operator ID: %s", err)
		}
		assetID, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyAssetID)
		if err != nil {
			return fmt.Errorf("error while getting asset ID: %s", err)
		}
		totalAmount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyTotalAmount)
		if err != nil {
			return fmt.Errorf("error while getting deposit amount: %s", err)
		}
		pendingUndelegationAmount, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyPendingUndelegationAmount)
		if err != nil {
			return fmt.Errorf("error while getting pending undelegation amount: %s", err)
		}
		totalShare, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyTotalShare)
		if err != nil {
			return fmt.Errorf("error while getting total share: %s", err)
		}
		operatorShare, err := juno.FindAttributeByKey(event, assetstypes.AttributeKeyOperatorShare)
		if err != nil {
			return fmt.Errorf("error while getting operator share: %s", err)
		}
		asset := types.NewOperatorAssetFromStr(
			operatorAddr.Value, assetID.Value,
			totalAmount.Value, pendingUndelegationAmount.Value,
			totalShare.Value, operatorShare.Value,
		)
		if err := m.db.SaveOperatorAsset(asset); err != nil {
			return fmt.Errorf("error while saving operator asset: %s", err)
		}
	}
	return nil
}
