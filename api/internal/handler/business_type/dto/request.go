package dto

import "uuid"

type UpsertRequest struct {
	InfraTypeID *uuid.UUID `json:"infra_type_id"`
}
