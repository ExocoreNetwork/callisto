package local

import (
	"fmt"

	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v5/node/local"

	"github.com/cosmos/cosmos-sdk/types/query"

	epochssource "github.com/forbole/callisto/v4/modules/epochs/source"
)

// interface guard
var (
	_ epochssource.Source = &Source{}
)

// Source implements epochssource.Source using a local node
type Source struct {
	*local.Source
	querier epochstypes.QueryServer
}

// NewSource implements a new Source instance
func NewSource(source *local.Source, querier epochstypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetEpochInfos implements epochssource.Source
func (s Source) GetEpochInfos(height int64) ([]epochstypes.EpochInfo, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var ret []epochstypes.EpochInfo
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.querier.EpochInfos(
			sdk.WrapSDKContext(ctx),
			&epochstypes.QueryEpochsInfoRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 1000,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		ret = append(ret, res.Epochs...)
	}

	return ret, nil
}

// GetEpochInfo implements epochssource.Source
func (s Source) GetEpochInfo(height int64, epochID string) (epochstypes.EpochInfo, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return epochstypes.EpochInfo{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.EpochInfo(
		sdk.WrapSDKContext(ctx),
		&epochstypes.QueryEpochInfoRequest{Identifier: epochID},
	)
	if err != nil {
		return epochstypes.EpochInfo{}, err
	}

	return res.Epoch, nil
}
