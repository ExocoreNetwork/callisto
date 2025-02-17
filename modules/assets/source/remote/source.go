package remote

import (
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/forbole/juno/v5/node/remote"

	assetssource "github.com/forbole/callisto/v4/modules/assets/source"
)

// interface guard
var (
	_ assetssource.Source = &Source{}
)

// Source implements assetssource.Source using a remote node
type Source struct {
	*remote.Source
	querier assetstypes.QueryClient
}

// NewSource implements a new Source instance
func NewSource(source *remote.Source, querier assetstypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetParams implements assetssource.Source
func (s Source) GetParams(height int64) (assetstypes.Params, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	res, err := s.querier.Params(
		ctx, &assetstypes.QueryParamsRequest{},
	)
	if err != nil {
		return assetstypes.Params{}, err
	}

	return *res.Params, nil
}
