package utils

import (
	"github.com/shopspring/decimal"
)

func ParseString2DBDecimal(str string) decimal.Decimal {
	scale := decimal.NewFromInt(1_000_000_000_000_000_000)
	value, _ := decimal.NewFromString(str)
	result := value.Mul(scale)
	return result
}

func ParseDBDecimal2String(originValue decimal.Decimal) string {
	scale := decimal.NewFromInt(1_000_000_000_000_000_000)
	value := originValue.Div(scale)
	return value.String()
}
