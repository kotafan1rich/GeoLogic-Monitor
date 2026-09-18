package domain

import "uuid"

type InfraTypeFeatures struct {
	TypeID          uuid.UUID
	Slug            string
	Weight          float64
	MaxRadiusMeters uint16
	ObjectCount     uint64
	NearestMeters   *float64
	MeanMeters      *float64
	DistanceSum     float64
	Influence       float64
	Count100m       uint64
	Count300m       uint64
	Count500m       uint64
	IsCompetitor    bool
}

type LocationFeatures struct {
	BusinessTypeID uuid.UUID
	InfraTypes     []InfraTypeFeatures
}
