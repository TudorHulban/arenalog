package arenalog

const _DeltaEstimation = 128
const _PrecisionFloat = 12

// Baseline for wrappers, commas, structures, and newlines.
const (
	_EstimatedBaselineCapacityJSON = 32
	_EstimatedBaselineCapacityRaw  = 16
	_EstimatedBaselineRoot         = 32
	_EstimatedBaselineFields       = 48
	_EstimatedBaselineFieldCount   = 48
)

const (
	MessageSmallSize  uint32 = 256
	MessageMediumSize uint32 = 512
	MessageJSONSize   uint32 = 768
	MessageLargeSize  uint32 = 2048
	MessageExtraLarge uint32 = 4096
)

const delim = ": "
