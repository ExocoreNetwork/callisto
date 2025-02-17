package source

import (
	sdkmath "cosmossdk.io/math"
)

type Source interface {
	GetDelegatedAmount(
		height int64, stakerID string, assetID string, operatorAddr string,
	) (sdkmath.Int, error)
}
