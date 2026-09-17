package dto

import (
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/business_type/model"
	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

func ToDomain(businessType model.BusinessType) *domain.BusinessType {
	return &domain.BusinessType{
		ID:          businessType.ID,
		InfraTypeID: businessType.InfraTypeID,
		InfraType: domain.InfraType{
			ID:        businessType.InfraType.ID,
			Slug:      businessType.InfraType.Slug,
			Name:      businessType.InfraType.Name,
			Weight:    businessType.InfraType.Weight,
			MaxRadius: businessType.InfraType.MaxRadius,
		},
	}
}

func ToModel(businessType domain.BusinessType) *model.BusinessType {
	return &model.BusinessType{
		Base: basemodel.Base{
			ID: businessType.ID,
		},
		InfraTypeID: businessType.InfraTypeID,
	}
}
