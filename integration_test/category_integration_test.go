package integration_test

import (
	"net/http"
	"testing"
)

// CREATE CATEGORY

type createCategoryTest struct {
	Name           string
	Regexps        []string
	ExpectedStatus int
}

var createCategoryTestCases = []createCategoryTest{
	{"", []string{"some-content"}, http.StatusBadRequest},
	{"Category", []string{}, http.StatusBadRequest},
	{"Category", []string{""}, http.StatusBadRequest},
}

//func Test_givenCreateCategoryRequest_whenCreateCategory_thenReturnExpectedStatus(t *testing.T) {
//	for _, testCase := range createCategoryTestCases {
//
//	}
//}

// GET CATEGORY BY ID

type getByIdTest struct {
	Id             string
	ExpectedStatus int
}

var getByIdTestCases = []getByIdTest{
	{"0", http.StatusBadRequest},
	{"random_string", http.StatusBadRequest},
	{"-5", http.StatusBadRequest},
	{"1", http.StatusNotFound},
}

func Test_givenCategoryId_whenGetCategoryById_thenReturnExpectedStatus(t *testing.T) {
	for _, testCase := range getByIdTestCases {
		clearDatabase()
		req, _ := http.NewRequest("GET", "/category/"+testCase.Id, nil)
		response := executeRequest(req, nil)
		assertResponseCode(t, testCase.ExpectedStatus, response.Code)
	}
}
