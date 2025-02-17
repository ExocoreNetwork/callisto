package local

import (
	"fmt"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v5/node/local"

	assetssource "github.com/forbole/callisto/v4/modules/assets/source"
)

// interface guard
var (
	_ assetssource.Source = &Source{}
)

// Source implements assetssource.Source using a local node
type Source struct {
	*local.Source
	querier assetstypes.QueryServer
}

// NewSource implements a new Source instance
func NewSource(source *local.Source, querier assetstypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetParams implements assetssource.Source
func (s Source) GetParams(height int64) (assetstypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return assetstypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(
		sdk.WrapSDKContext(ctx),
		&assetstypes.QueryParamsRequest{},
	)
	if err != nil {
		return assetstypes.Params{}, err
	}

	return *res.Params, nil
}
