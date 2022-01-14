package currency

import (
	"fmt"
	goRegexp "regexp"
	"strconv"
)

func FormatAsCurrency(amount int) string {
	beforeDecimalSeparator := amount / 100
	afterDecimalSeparator := amount % 100
	if amount < 0 {
		afterDecimalSeparator *= -1
	}
	return fmt.Sprintf("%d.%.2d", beforeDecimalSeparator, afterDecimalSeparator)
}

func CurrencyToInteger(currency string) int {
	re, _ := goRegexp.Compile("\\.|,")
	currency = re.ReplaceAllString(currency, "")
	currencyAsInt, _ := strconv.Atoi(currency)
	return currencyAsInt
}

func IsValidCurrencyString(s string) bool {
	isValid, _ := goRegexp.MatchString("^((-?[0-9]+)|(-?[0-9]+(\\.|,)[0-9]{2}))$", s)
	return isValid
}
