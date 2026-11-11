package shared

import (
	"time"

	"github.com/google/uuid"
)

// BaseEntity contains standard audit fields across all domain models.
type BaseEntity struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
}

// BranchScoped ensures entities comply with multi-branch isolation.
type BranchScoped struct {
	BranchID uuid.UUID `json:"branch_id"`
}

