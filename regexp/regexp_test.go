package regexp

import "testing"

type regexpMatchesTest struct {
	Content        string
	TestValue      string
	ExpectedResult bool
}

var regexpMatchesTests = []regexpMatchesTest{
	{"^abc$", "abc", true},
	{"^abc$", "abcc", false},
	{".*KAUFLAND.*", "ASD-A12315ASDAKAUFLAND123123", true},
}

func Test_givenRegexpAndText_whenRegexpMatch_thenReturnTrue(t *testing.T) {
	for _, testCase := range regexpMatchesTests {
		regexp := Regexp{Content: testCase.Content}
		result := regexp.Matches(testCase.TestValue)
		if result != testCase.ExpectedResult {
			t.Errorf("Regexp: %s, Value: %s, Expected %t, Got %t", testCase.Content, testCase.TestValue, testCase.ExpectedResult, result)
		}
	}
}
