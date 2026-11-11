package dto

import "github.com/google/uuid"

type YearEndClosingRequest struct {
	FiscalYear int        `json:"fiscal_year"`
	ClosedBy   *uuid.UUID `json:"closed_by,omitempty"`
}

