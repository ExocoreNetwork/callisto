package source

import (
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
)

type Source interface {
	GetParams(height int64) (assetstypes.Params, error)
}
