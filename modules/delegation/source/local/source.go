package local

import (
	"fmt"

	sdkmath "cosmossdk.io/math"

	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v5/node/local"

	delegationsource "github.com/forbole/callisto/v4/modules/delegation/source"
)

// interface guard
var (
	_ delegationsource.Source = &Source{}
)

// Source implements delegationsource.Source using a local node
type Source struct {
	*local.Source
	querier delegationtypes.QueryServer
}

// NewSource implements a new Source instance
func NewSource(source *local.Source, querier delegationtypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetDelegatedAmount implements delegationsource.Source
func (s Source) GetDelegatedAmount(
	height int64, stakerID string, assetID string, operatorAddr string,
) (sdkmath.Int, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return sdkmath.ZeroInt(), fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.QuerySingleDelegationInfo(
		sdk.WrapSDKContext(ctx),
		&delegationtypes.SingleDelegationInfoReq{
			StakerID:     stakerID,
			AssetID:      assetID,
			OperatorAddr: operatorAddr,
		},
	)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	return res.MaxUndelegatableAmount, nil
}
