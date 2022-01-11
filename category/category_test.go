package category

import (
	"github.com/remes2000/amu_financial_summary/regexp"
	"testing"
)

var alwaysTrueRegexp = regexp.Regexp{Content: ".*"}
var alwaysFalseRegexp = regexp.Regexp{Content: "(?!x)x"}

type categoryMatchesTest struct {
	Regexps        []regexp.Regexp
	ExpectedResult bool
	TestCaseName   string
}

var categoryMatchesTests = []categoryMatchesTest{
	{[]regexp.Regexp{}, false, "No regexps"},
	{[]regexp.Regexp{alwaysFalseRegexp, alwaysFalseRegexp}, false, "All false regexps"},
	{[]regexp.Regexp{alwaysFalseRegexp, alwaysFalseRegexp, alwaysTrueRegexp}, true, "One regexp true"},
}

func Test_givenRegexpAndText_whenRegexpMatch_thenReturnTrue(t *testing.T) {
	for _, testCase := range categoryMatchesTests {
		category := Category{Regexps: testCase.Regexps}
		result := category.Matches("randomtitle")
		if result != testCase.ExpectedResult {
			t.Errorf("[%s] Expected %t. Got %t", testCase.TestCaseName, testCase.ExpectedResult, result)
		}
	}
}
