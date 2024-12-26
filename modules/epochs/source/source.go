package source

import (
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
)

type Source interface {
	GetEpochInfos(height int64) ([]epochstypes.EpochInfo, error)
	GetEpochInfo(height int64, epochID string) (epochstypes.EpochInfo, error)
}
