package types

import (
	exominttypes "github.com/ExocoreNetwork/exocore/x/exomint/types"
)

// ExomintParams represents the x/exomint parameters
type ExomintParams struct {
	exominttypes.Params
	Height int64
}

// NewExomintParams allows to build a new ExomintParams instance
func NewExomintParams(params exominttypes.Params, height int64) *ExomintParams {
	return &ExomintParams{
		Params: params,
		Height: height,
	}
}

// MintHistory represents the mint history
type MintHistory struct {
	// this must be a string and not sdkmath.Int because a string can be directly
	// inserted into the database, while an sdkmath.Int cannot.
	Amount      string
	Height      int64
	EpochID     string
	EpochNumber int64
	Denom       string
}

// NewMintHistory allows to build a new MintHistory instance
func NewMintHistory(
	height int64, amount string, epochID string, epochNumber int64, denom string,
) *MintHistory {
	return &MintHistory{
		Height:      height,
		Amount:      amount,
		EpochID:     epochID,
		EpochNumber: epochNumber,
		Denom:       denom,
	}
}
