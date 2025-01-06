package utils

import (
	"fmt"
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

func ConvertDecimal(amount string) string {
	if amount == "" {
		return ""
	}
	decAmount, err := decimal.NewFromString(amount)
	if err != nil {
		fmt.Println("Error converting string to decimal:", err)
		return ""
	}
	// 除以 10^18
	convertedAmount := decAmount.Div(decimal.NewFromInt(1e18))
	// 格式化为字符串，保留 18 位小数（或者可以调整小数位数，根据需求）
	return convertedAmount.String()
}
