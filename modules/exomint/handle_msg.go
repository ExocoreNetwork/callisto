package exomint

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/authz"

	exominttypes "github.com/ExocoreNetwork/exocore/x/exomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/callisto/v4/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsgExec implements AuthzMessageModule. It handles the case wherein
// a grantee is permitted to execute a message on behalf of the granter.
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	switch cosmosMsg := msg.(type) {
	case *exominttypes.MsgUpdateParams:
		return m.handleMsgUpdateParams(tx.Height, cosmosMsg)
	}
	return nil
}

// handleMsgUpdateParams handles the MsgUpdateParams message type by overwriting
// the existing parameters with the new ones in the database.
func (m *Module) handleMsgUpdateParams(
	height int64, _ *exominttypes.MsgUpdateParams,
) error {
	// we can parse the params from here, or we can just load them from the module source.
	// it is easier to do the latter.
	params, err := m.source.GetParams(height)
	if err != nil {
		return fmt.Errorf("error while getting params: %s", err)
	}
	err = m.db.SaveExomintParams(types.NewExomintParams(params, height))
	if err != nil {
		return fmt.Errorf("error while saving params: %s", err)
	}
	return nil
}
