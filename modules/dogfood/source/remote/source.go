package remote

import (
	dogfoodtypes "github.com/ExocoreNetwork/exocore/x/dogfood/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v5/node/remote"

	dogfoodsource "github.com/forbole/callisto/v4/modules/dogfood/source"
)

// interface guard
var (
	_ dogfoodsource.Source = &Source{}
)

// Source implements dogfoodsource.Source using a remote node
type Source struct {
	*remote.Source
	querier dogfoodtypes.QueryClient
}

// NewSource implements a new Source instance
func NewSource(source *remote.Source, querier dogfoodtypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetParams implements dogfoodsource.Source
func (s Source) GetParams(height int64) (dogfoodtypes.Params, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	res, err := s.querier.Params(
		ctx, &dogfoodtypes.QueryParamsRequest{},
	)
	if err != nil {
		return dogfoodtypes.Params{}, err
	}

	return res.Params, nil
}

// GetValidators implements dogfoodsource.Source
func (s Source) GetValidators(height int64) ([]dogfoodtypes.ExocoreValidator, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var validators []dogfoodtypes.ExocoreValidator
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.querier.Validators(
			ctx,
			&dogfoodtypes.QueryAllValidatorsRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 validators at time
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		validators = append(validators, res.Validators...)
	}

	return validators, nil
}
