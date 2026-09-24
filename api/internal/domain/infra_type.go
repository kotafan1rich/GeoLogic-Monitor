package domain

import (
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
)

type InfraType struct {
	ID        uuid.UUID
	Slug      string
	Name      string
	Weight    float64
	MaxRadius uint16
}

func NewInfraType(slug, name string, weight float64, maxRadius uint16) (*InfraType, error) {
	if slug == "" {
		return nil, errs.ErrInvalidSlug
	}
	if name == "" {
		return nil, errs.ErrInvalidName
	}
	if weight <= 0 {
		return nil, errs.ErrInvalidWeight
	}
	if maxRadius <= 0 {
		return nil, errs.ErrInvalidRadius
	}
	return &InfraType{
		Slug:      slug,
		Name:      name,
		Weight:    weight,
		MaxRadius: maxRadius,
	}, nil
}

func (t *InfraType) UpdateSlug(slug string) error {
	if slug == "" {
		return errs.ErrInvalidSlug
	}
	t.Slug = slug
	return nil
}

func (t *InfraType) UpdateName(name string) error {
	if name == "" {
		return errs.ErrInvalidName
	}
	t.Name = name
	return nil
}

func (t *InfraType) UpdateWeight(weight float64) error {
	if weight <= 0 {
		return errs.ErrInvalidWeight
	}
	t.Weight = weight
	return nil
}

func (t *InfraType) UpdateMaxRadius(maxRadius uint16) error {
	if maxRadius <= 0 {
		return errs.ErrInvalidRadius
	}
	t.MaxRadius = maxRadius
	return nil
}

func (t *InfraType) Update(slug, name *string, weight *float64, maxRadius *uint16) error {
	if slug != nil {
		err := t.UpdateSlug(*slug)
		if err != nil {
			return err
		}
	}

	if name != nil {
		err := t.UpdateName(*name)
		if err != nil {
			return err
		}
	}

	if weight != nil {
		err := t.UpdateWeight(*weight)
		if err != nil {
			return err
		}
	}

	if maxRadius != nil {
		err := t.UpdateMaxRadius(*maxRadius)
		if err != nil {
			return err
		}
	}
	return nil
}
