package remote

import (
	exominttypes "github.com/ExocoreNetwork/exocore/x/exomint/types"
	"github.com/forbole/juno/v5/node/remote"

	exomintsource "github.com/forbole/callisto/v4/modules/exomint/source"
)

// interface guard
var (
	_ exomintsource.Source = &Source{}
)

// Source implements exomintsource.Source using a remote node
type Source struct {
	*remote.Source
	querier exominttypes.QueryClient
}

// NewSource implements a new Source instance
func NewSource(source *remote.Source, querier exominttypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetParams implements exomintsource.Source
func (s Source) GetParams(height int64) (exominttypes.Params, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	res, err := s.querier.Params(
		ctx, &exominttypes.QueryParamsRequest{},
	)
	if err != nil {
		return exominttypes.Params{}, err
	}

	return res.Params, nil
}
