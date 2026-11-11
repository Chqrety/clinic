package decimal

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// Currency standard constants for Indonesian Rupiah (IDR).
const (
	CurrencyCode   = "IDR"
	CurrencySymbol = "Rp"
	CurrencyName   = "Indonesian Rupiah"
)

// Decimal aliases shopspring/decimal.Decimal for consistent application-wide usage.
type Decimal = decimal.Decimal

// Common zero and unit values
var (
	Zero = decimal.Zero
	One  = decimal.NewFromInt(1)
)

// FromFloat creates a Decimal from a float64 (safely converted via string).
func FromFloat(f float64) Decimal {
	return decimal.NewFromFloat(f)
}

// NewFromFloat creates a Decimal from a float64 (alias for FromFloat).
func NewFromFloat(f float64) Decimal {
	return decimal.NewFromFloat(f)
}

// FromInt creates a Decimal from an int64.
func FromInt(i int64) Decimal {
	return decimal.NewFromInt(i)
}

// FromString parses a string into Decimal.
func FromString(s string) (Decimal, error) {
	return decimal.NewFromString(s)
}

// NewFromString parses a string into Decimal (alias for FromString).
func NewFromString(s string) (Decimal, error) {
	return decimal.NewFromString(s)
}

// RoundCurrency rounds a Decimal to standard 2 decimal places for financial calculations.
func RoundCurrency(d Decimal) Decimal {
	return d.Round(2)
}

// IsBalanced checks if total debits equal total credits within accounting tolerance.
func IsBalanced(debit, credit Decimal) bool {
	return debit.Equal(credit)
}

// NullDecimal represents a nullable Decimal for SQL integration.
type NullDecimal struct {
	Decimal Decimal
	Valid   bool
}

// Scan implements database/sql.Scanner interface.
func (nd *NullDecimal) Scan(value interface{}) error {
	if value == nil {
		nd.Decimal, nd.Valid = Zero, false
		return nil
	}
	nd.Valid = true
	return nd.Decimal.Scan(value)
}

// Value implements driver.Valuer interface.
func (nd NullDecimal) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}
	return nd.Decimal.Value()
}

// FormatIDR formats a Decimal to Indonesian Rupiah standard format (e.g. Rp 1.500.000,00).
func FormatIDR(d Decimal) string {
	fixed := d.StringFixed(2)
	isNegative := false
	if len(fixed) > 0 && fixed[0] == '-' {
		isNegative = true
		fixed = fixed[1:]
	}

	parts := strings.Split(fixed, ".")
	intPart := parts[0]
	decPart := "00"
	if len(parts) > 1 {
		decPart = parts[1]
	}

	var formattedInt strings.Builder
	n := len(intPart)
	for i, ch := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			formattedInt.WriteRune('.')
		}
		formattedInt.WriteRune(ch)
	}

	if isNegative {
		return fmt.Sprintf("-Rp %s,%s", formattedInt.String(), decPart)
	}
	return fmt.Sprintf("Rp %s,%s", formattedInt.String(), decPart)
}

// FormatRupiahNoSen formats a Decimal to Rupiah without fractional sen (e.g. Rp 1.500.000).
func FormatRupiahNoSen(d Decimal) string {
	rounded := d.Round(0).StringFixed(0)
	isNegative := false
	if len(rounded) > 0 && rounded[0] == '-' {
		isNegative = true
		rounded = rounded[1:]
	}

	var formattedInt strings.Builder
	n := len(rounded)
	for i, ch := range rounded {
		if i > 0 && (n-i)%3 == 0 {
			formattedInt.WriteRune('.')
		}
		formattedInt.WriteRune(ch)
	}

	if isNegative {
		return fmt.Sprintf("-Rp %s", formattedInt.String())
	}
	return fmt.Sprintf("Rp %s", formattedInt.String())
}


