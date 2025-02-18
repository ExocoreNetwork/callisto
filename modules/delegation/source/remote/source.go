package remote

import (
	sdkmath "cosmossdk.io/math"

	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/forbole/juno/v5/node/remote"

	delegationsource "github.com/forbole/callisto/v4/modules/delegation/source"
)

// interface guard
var (
	_ delegationsource.Source = &Source{}
)

// Source implements delegationsource.Source using a remote node
type Source struct {
	*remote.Source
	querier delegationtypes.QueryClient
}

// NewSource implements a new Source instance
func NewSource(source *remote.Source, querier delegationtypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetDelegatedAmount implements delegationsource.Source
func (s Source) GetDelegatedAmount(
	height int64, stakerID string, assetID string, operatorAddr string,
) (sdkmath.Int, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	res, err := s.querier.QuerySingleDelegationInfo(
		ctx,
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
