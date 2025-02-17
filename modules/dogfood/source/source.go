package source

import (
	dogfoodtypes "github.com/ExocoreNetwork/exocore/x/dogfood/types"
)

type Source interface {
	GetParams(height int64) (dogfoodtypes.Params, error)
	GetValidators(height int64) ([]dogfoodtypes.ExocoreValidator, error)
}
