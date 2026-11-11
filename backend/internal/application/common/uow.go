package common

import (
	"context"
)

// UnitOfWork defines atomic transaction execution across multiple repositories.
// Essential for ACID compliance in Double-Entry Accounting operations.
type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

