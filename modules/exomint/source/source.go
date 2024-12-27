package source

import (
	exominttypes "github.com/ExocoreNetwork/exocore/x/exomint/types"
)

type Source interface {
	GetParams(height int64) (exominttypes.Params, error)
}
