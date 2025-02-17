package local

import (
	"fmt"

	exominttypes "github.com/ExocoreNetwork/exocore/x/exomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v5/node/local"

	exomintsource "github.com/forbole/callisto/v4/modules/exomint/source"
)

// interface guard
var (
	_ exomintsource.Source = &Source{}
)

// Source implements exomintsource.Source using a local node
type Source struct {
	*local.Source
	querier exominttypes.QueryServer
}

// NewSource implements a new Source instance
func NewSource(source *local.Source, querier exominttypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetParams implements exomintsource.Source
func (s Source) GetParams(height int64) (exominttypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return exominttypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(
		sdk.WrapSDKContext(ctx),
		&exominttypes.QueryParamsRequest{},
	)
	if err != nil {
		return exominttypes.Params{}, err
	}

	return res.Params, nil
}
