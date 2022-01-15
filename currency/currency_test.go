package currency

import (
	"testing"
)

// IsValidCurrency

type currencyIsValidTest struct {
	Value          string
	ExpectedResult bool
}

var currencyIsValidTestCases = []currencyIsValidTest{
	{"", false},
	{".55", false},
	{",55", false},
	{"randomstring", false},
	{"255,111", false},
	{"255.111", false},
	{"-255.111", false},
	{"222222.11", true},
	{"-222222.11", true},
	{"222222,11", true},
	{"-222222,11", true},
}

func Test_givenInputString_whenCheckIsCurrencyValid_thenReturnTrue(t *testing.T) {
	for _, testCase := range currencyIsValidTestCases {
		isValid := IsValidCurrencyString(testCase.Value)
		if testCase.ExpectedResult != isValid {
			t.Errorf("Expected %t Got %t for value: %s", testCase.ExpectedResult, isValid, testCase.Value)
		}
	}
}

// FormatAsCurrency

type formatAsCurrencyTest struct {
	Input          int
	ExpectedResult string
}

var formatAsCurrencyTestCases = []formatAsCurrencyTest{
	{25, "0.25"},
	{0, "0.00"},
	{12007, "120.07"},
	{-12007, "-120.07"},
}

func Test_givenInputInteger_whenFormatAsCurrency_thenReturnIntegerInCurrencyFormat(t *testing.T) {
	for _, testCase := range formatAsCurrencyTestCases {
		result := FormatAsCurrency(testCase.Input)
		if testCase.ExpectedResult != result {
			t.Errorf("Expected %s Got %s for input %d", testCase.ExpectedResult, result, testCase.Input)
		}
	}
}

// CurrencyToInteger

type currencyToIntegerTest struct {
	Input          string
	ExpectedResult int
}

var currencyToIntegerTestCases = []currencyToIntegerTest{
	{"0.25", 25},
	{"0,25", 25},
	{"0.00", 0},
	{"0,00", 0},
	{"120.07", 12007},
	{"120,07", 12007},
	{"-120.07", -12007},
	{"-120,07", -12007},
}

func Test_givenInputCurrencyString_whenCurrencyToInteger_thenReturnCorrectInteger(t *testing.T) {
	for _, testCase := range currencyToIntegerTestCases {
		result := CurrencyToInteger(testCase.Input)
		if testCase.ExpectedResult != result {
			t.Errorf("Expected %d Got %d for input %s", testCase.ExpectedResult, result, testCase.Input)
		}
	}
}
