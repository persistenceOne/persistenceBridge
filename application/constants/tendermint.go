package constants

import "time"

const (
	// TendermintAvgBlockTime Putting this as constant because different network will have different block times,
	// and we will release different version of bridge for each network as SCs will be different.
	TendermintAvgBlockTime = 6000 * time.Millisecond
)
