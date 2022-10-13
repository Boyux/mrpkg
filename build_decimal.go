//go:build decimal

package mrpkg

import "github.com/shopspring/decimal"

func init() {
	RegisterTypeConstructor[decimal.Decimal](func() any { return decimal.Zero })
}
