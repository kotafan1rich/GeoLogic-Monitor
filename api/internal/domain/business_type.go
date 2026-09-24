package domain

import "uuid"

type BusinessType struct {
	ID          uuid.UUID
	InfraTypeID uuid.UUID
	InfraType   InfraType
}

func NewBusinessType(infraTypeID uuid.UUID) *BusinessType {
	return &BusinessType{
		InfraTypeID: infraTypeID,
	}
}

