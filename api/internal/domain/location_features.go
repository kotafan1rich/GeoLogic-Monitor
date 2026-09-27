package domain

import "uuid"

// LocationFeatures contains walking distances supplied by the caller, never coordinates.
type LocationFeatures struct {
	BusinessTypeID uuid.UUID
	BusinessSlug   string
	Objects        []RatingObjectDistance
	Profile        map[string]float64
	StreetFactor   *float64
}

type RatingObjectDistance struct {
	Type           string
	DistanceMeters float64
}

// RatingAssessment is transient; only Rating and an external timestamp are persisted.
type RatingAssessment struct {
	Rating             float64
	Confidence         float64
	Breakdown          []RatingContribution
	CompetitionPenalty float64
}

type RatingContribution struct {
	Type         string
	Weight       float64
	Sat          float64
	Contribution float64
}
