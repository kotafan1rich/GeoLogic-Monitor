package dto

import "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"

func InfraTypeToResponse(infraType domain.InfraType) *InfraTypeResponse {
	return &InfraTypeResponse{
		ID:        infraType.ID,
		Slug:      infraType.Slug,
		Name:      infraType.Name,
		Weight:    infraType.Weight,
		MaxRadius: infraType.MaxRadius,
	}
}

func InfraObjectToResponse(infraObject domain.InfraObject) *InfraObjectResponse {
	return &InfraObjectResponse{
		ID:         infraObject.ID,
		ExternalID: infraObject.ExternalID,
		TypeID:     infraObject.TypeID,
		Lat:        infraObject.GeoPoint.Lat,
		Lon:        infraObject.GeoPoint.Lng,
		Address:    infraObject.Address,
		Name:       infraObject.Name,
	}
}
