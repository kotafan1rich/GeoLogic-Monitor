package dto

import (
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_type/model"
)

func ToDomain(infraType model.InfraType) *domain.InfraType {
	return &domain.InfraType{
		ID:        infraType.ID,
		Slug:      infraType.Slug,
		Name:      infraType.Name,
		Weight:    infraType.Weight,
		MaxRadius: infraType.MaxRadius,
	}
}

func ToModel(infraType domain.InfraType) *model.InfraType {
	return &model.InfraType{
		ID:        infraType.ID,
		Slug:      infraType.Slug,
		Name:      infraType.Name,
		Weight:    infraType.Weight,
		MaxRadius: infraType.MaxRadius,
	}
}
