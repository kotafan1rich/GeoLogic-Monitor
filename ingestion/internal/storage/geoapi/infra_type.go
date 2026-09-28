package geoapi

import (
	"context"
	"fmt"
)

const infraTypeEndpoint = "internal/v1/infra-types"

func (c *Client) PutInfraType(ctx context.Context, obj InfraTypeInput) (InfraType, error) {
	const op = "geoapi.Client.PutInfraType"

	var raw InfraType

	if err := c.put(ctx, infraTypeEndpoint, obj, &raw); err != nil {
		return InfraType{}, fmt.Errorf("%s: %w", op, err)
	}

	return raw, nil
}

type TypeRegistry struct {
	client   *Client
	defs     map[string]InfraTypeInput
	cache    *TypesCache
	business *BusinessTypesCache
}

func NewTypeRegistry(c *Client, defs []InfraTypeInput) (*TypeRegistry, error) {
	const op = "geoapi.NewTypeRegistry"

	if c == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidClient)
	}

	if len(defs) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrEmptyInfraTypes)
	}

	bySlug := make(map[string]InfraTypeInput, len(defs))

	for _, def := range defs {
		if def.Slug == "" {
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidInfraType)
		}

		if _, ok := bySlug[def.Slug]; ok {
			return nil, fmt.Errorf("%s: %w [%s]", op, ErrDuplicateInfraType, def.Slug)
		}

		bySlug[def.Slug] = def
	}

	return &TypeRegistry{
		client:   c,
		defs:     bySlug,
		cache:    NewTypesCache(),
		business: NewBusinessTypesCache(),
	}, nil
}

func MustNewTypeRegistry(c *Client, defs []InfraTypeInput) *TypeRegistry {
	const op = "geoapi.MustNewTypeRegistry"

	r, err := NewTypeRegistry(c, defs)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init infra type registry: %v", op, err))
	}
	return r
}

func (r *TypeRegistry) TypeID(ctx context.Context, slug string) (string, error) {
	const op = "geoapi.TypeRegistry.TypeID"

	if id, ok := r.cache.Get(slug); ok {
		return string(id), nil
	}

	def, ok := r.defs[slug]
	if !ok {
		return "", fmt.Errorf("%s: %w [%s]", op, ErrUnknownInfraType, slug)
	}

	t, err := r.client.PutInfraType(ctx, def)
	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, ErrPutInfraType, err)
	}

	if t.ID == "" {
		return "", fmt.Errorf("%s: %w", op, ErrEmptyTypeID)
	}

	id := TypeID(t.ID)
	r.cache.Set(slug, id)

	return string(id), nil
}

func (r *TypeRegistry) SetBusinessType(id, slug string) {
	if id == "" {
		return
	}

	r.business.Set(id, slug)
}

func (r *TypeRegistry) BusinessType(ctx context.Context, id string) (InfraTypeInput, bool, error) {
	const op = "geoapi.TypeRegistry.BusinessType"

	slug, ok := r.business.Get(id)
	if !ok {
		if err := r.loadBusinessTypes(ctx); err != nil {
			return InfraTypeInput{}, false, fmt.Errorf("%s: %w", op, err)
		}

		if slug, ok = r.business.Get(id); !ok {
			return InfraTypeInput{}, false, nil
		}
	}

	def, ok := r.defs[slug]

	return def, ok, nil
}

func (r *TypeRegistry) BusinessTypeID(ctx context.Context, slug string) (string, bool, error) {
	const op = "geoapi.TypeRegistry.BusinessTypeID"

	if id, ok := r.business.ID(slug); ok {
		return id, true, nil
	}

	if err := r.loadBusinessTypes(ctx); err != nil {
		return "", false, fmt.Errorf("%s: %w", op, err)
	}

	id, ok := r.business.ID(slug)

	return id, ok, nil
}

func (r *TypeRegistry) loadBusinessTypes(ctx context.Context) error {
	types, err := r.client.BusinessTypes(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrGetBusinessTypes, err)
	}

	for _, bt := range types {
		if bt.ID == "" || bt.InfraType.Slug == "" {
			continue
		}

		r.business.Set(bt.ID, bt.InfraType.Slug)
	}

	return nil
}
