package integration_test

import (
	"bytes"
	"net/http"
	"testing"
)

// IMPORT TRANSACTIONS
type importTransactionsTest struct {
	Content        string
	ExpectedStatus int
	TestCaseName   string
}

var importTransactionsTestCases = []importTransactionsTest{
	{"", http.StatusBadRequest, "Bad request when empty input"},
	{"[]", http.StatusBadRequest, "Bad request when emtpy array in input"},
	{`[{"amount": 21.37, "date": "28-12-2021"}]`, http.StatusBadRequest, "Bad request when title not specified"},
	{`[{"title": "abc", "date": "28-12-2021"}]`, http.StatusBadRequest, "Bad request when amount not specified"},
	{`[{"title": "abc", "amount": 21.37}]`, http.StatusBadRequest, "Bad request when date not specified"},
	{`[{"title": "abc", "amount": 21.37, "date": "randomstring"}]`, http.StatusBadRequest, "Bad request when date not valid format"},
}

func Test_givenImportTransactionsJson_whenImportTransactions_thenReturnExpectedStatusAndPerformImport(t *testing.T) {
	for _, testCase := range importTransactionsTestCases {
		clearDatabase()

		req, _ := http.NewRequest("POST", "/account-transaction", bytes.NewBuffer([]byte(testCase.Content)))
		response := executeRequest(req, nil)

		if testCase.ExpectedStatus != response.Code {
			t.Errorf("[%s] Expected response code %d. Got %d", testCase.TestCaseName, testCase.ExpectedStatus, response.Code)
			continue
		}
	}
}
