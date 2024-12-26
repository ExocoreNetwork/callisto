package remote

import (
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	"github.com/forbole/juno/v5/node/remote"

	"github.com/cosmos/cosmos-sdk/types/query"
	epochssource "github.com/forbole/callisto/v4/modules/epochs/source"
)

// interface guard
var (
	_ epochssource.Source = &Source{}
)

// Source implements epochssource.Source using a remote node
type Source struct {
	*remote.Source
	querier epochstypes.QueryClient
}

// NewSource implements a new Source instance
func NewSource(source *remote.Source, querier epochstypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetEpochInfos implements epochssource.Source
func (s Source) GetEpochInfos(height int64) ([]epochstypes.EpochInfo, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var ret []epochstypes.EpochInfo
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.querier.EpochInfos(
			ctx,
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
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	res, err := s.querier.EpochInfo(
		ctx,
		&epochstypes.QueryEpochInfoRequest{
			Identifier: epochID,
		},
	)
	if err != nil {
		return epochstypes.EpochInfo{}, err
	}

	return res.Epoch, nil
}
