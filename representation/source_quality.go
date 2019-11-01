package representation

const (
	// SourceQualityPerfect indicates that the representation is perfect
	// quality with no degradation.
	SourceQualityPerfect float32 = 1.0

	// SourceQualityNearlyPerfect indicates the threshold of noticeable loss
	// of quality for the representation.
	SourceQualityNearlyPerfect float32 = 0.9

	// SourceQualityAcceptable indicates that the representation has
	// noticeable but acceptable quality reduction.
	SourceQualityAcceptable float32 = 0.8

	// SourceQualityBarelyAcceptable indicates that the representation
	// has barely acceptable quality.
	SourceQualityBarelyAcceptable float32 = 0.5

	// SourceQualitySeverelyDegraded indicates that the representation
	// has severely degraded quality.
	SourceQualitySeverelyDegraded float32 = 0.3

	// SourceQualityCompletelyDegraded indicates that the representation
	// has completed degraded quality.
	SourceQualityCompletelyDegraded float32 = 0.0
)
