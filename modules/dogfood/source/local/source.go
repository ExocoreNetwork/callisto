package local

import (
	"fmt"

	dogfoodtypes "github.com/ExocoreNetwork/exocore/x/dogfood/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v5/node/local"

	dogfoodsource "github.com/forbole/callisto/v4/modules/dogfood/source"
)

// interface guard
var (
	_ dogfoodsource.Source = &Source{}
)

// Source implements dogfoodsource.Source using a local node
type Source struct {
	*local.Source
	querier dogfoodtypes.QueryServer
}

// NewSource implements a new Source instance
func NewSource(source *local.Source, querier dogfoodtypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetParams implements dogfoodsource.Source
func (s Source) GetParams(height int64) (dogfoodtypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return dogfoodtypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(
		sdk.WrapSDKContext(ctx),
		&dogfoodtypes.QueryParamsRequest{},
	)
	if err != nil {
		return dogfoodtypes.Params{}, err
	}

	return res.Params, nil
}

// GetValidators implements dogfoodsource.Source
func (s Source) GetValidators(height int64) ([]dogfoodtypes.ExocoreValidator, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var validators []dogfoodtypes.ExocoreValidator
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.querier.Validators(
			sdk.WrapSDKContext(ctx),
			&dogfoodtypes.QueryAllValidatorsRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100,
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
