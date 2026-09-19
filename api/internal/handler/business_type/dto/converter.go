package dto

import "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"

func ToResponse(businessType domain.BusinessType) *BusinessTypeResponse {
	return &BusinessTypeResponse{
		ID:          businessType.ID,
		InfraTypeID: businessType.InfraTypeID,
		InfraType: InfraTypeResponse{
			ID:        businessType.InfraType.ID,
			Slug:      businessType.InfraType.Slug,
			Name:      businessType.InfraType.Name,
			Weight:    businessType.InfraType.Weight,
			MaxRadius: businessType.InfraType.MaxRadius,
		},
	}
}

func ToResponseList(businessTypes []domain.BusinessType) []*BusinessTypeResponse {
	result := make([]*BusinessTypeResponse, 0, len(businessTypes))
	for _, businessType := range businessTypes {
		result = append(result, ToResponse(businessType))
	}
	return result
}
