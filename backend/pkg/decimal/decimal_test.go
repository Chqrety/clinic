package decimal_test

import (
	"testing"

	"clinic-api/pkg/decimal"
)

func TestFormatIDR(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
		noSen    string
	}{
		{1500000.00, "Rp 1.500.000,00", "Rp 1.500.000"},
		{25000.50, "Rp 25.000,50", "Rp 25.001"},
		{0.00, "Rp 0,00", "Rp 0"},
		{-75000.00, "-Rp 75.000,00", "-Rp 75.000"},
		{500000000.00, "Rp 500.000.000,00", "Rp 500.000.000"},
	}

	for _, tt := range tests {
		d := decimal.FromFloat(tt.input)
		res := decimal.FormatIDR(d)
		if res != tt.expected {
			t.Errorf("FormatIDR(%f) = %s, expected %s", tt.input, res, tt.expected)
		}

		resNoSen := decimal.FormatRupiahNoSen(d)
		if resNoSen != tt.noSen {
			t.Errorf("FormatRupiahNoSen(%f) = %s, expected %s", tt.input, resNoSen, tt.noSen)
		}
	}
}

