package dto

import (
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_object/model"
	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

func ToDomain(infraObject model.InfraObject) *domain.InfraObject {
	return &domain.InfraObject{
		ID:       infraObject.ID,
		TypeID:   infraObject.TypeID,
		GeoPoint: domain.GeoPoint(infraObject.Location),
		Address:  infraObject.Address,
		Name:     infraObject.Name,
		Type: domain.InfraType{
			ID:        infraObject.Type.ID,
			Slug:      infraObject.Type.Slug,
			Name:      infraObject.Type.Name,
			Weight:    infraObject.Type.Weight,
			MaxRadius: infraObject.Type.MaxRadius,
		},
	}
}

func ToModel(infraObject domain.InfraObject) *model.InfraObject {
	return &model.InfraObject{
		ID:       infraObject.ID,
		TypeID:   infraObject.TypeID,
		Location: basemodel.GeoPoint(infraObject.GeoPoint),
		Address:  infraObject.Address,
		Name:     infraObject.Name,
	}
}
