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
	client *Client
	defs   map[string]InfraTypeInput
	cache  *TypesCache
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
		client: c,
		defs:   bySlug,
		cache:  NewTypesCache(),
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
