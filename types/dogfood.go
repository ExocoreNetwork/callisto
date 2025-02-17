package types

// OptOutExpiry represents the opt out expiry for a given operator.
type OptOutExpiry struct {
	EpochNumber  string
	OperatorAddr string
}

// NewOptOutExpiryFromStr creates a new OptOutExpiry instance from the given values in
// string format.
func NewOptOutExpiryFromStr(epochNumber, operatorAddr string) *OptOutExpiry {
	return &OptOutExpiry{
		EpochNumber:  epochNumber,
		OperatorAddr: operatorAddr,
	}
}

// ConsensusAddrToPrune represents the consensus address to prune for a given epoch.
type ConsensusAddrToPrune struct {
	EpochNumber   string
	ConsensusAddr string
}

// NewConsensusAddrToPruneFromStr creates a new ConsensusAddrToPrune instance from the given values in
// string format.
func NewConsensusAddrToPruneFromStr(epochNumber, consensusAddr string) *ConsensusAddrToPrune {
	return &ConsensusAddrToPrune{
		EpochNumber:   epochNumber,
		ConsensusAddr: consensusAddr,
	}
}

// UndelegationMaturity represents the undelegation maturity for a given operator.
type UndelegationMaturity struct {
	EpochNumber string
	RecordKey   string
}

// NewUndelegationMaturityFromStr creates a new UndelegationMaturity instance from the given values in
// string format.
func NewUndelegationMaturityFromStr(epochNumber, recordKey string) *UndelegationMaturity {
	return &UndelegationMaturity{
		EpochNumber: epochNumber,
		RecordKey:   recordKey,
	}
}
