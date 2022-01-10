package integration_test

import (
	"bytes"
	"encoding/json"
	"github.com/remes2000/amu_financial_summary/category"
	"github.com/remes2000/amu_financial_summary/global"
	"github.com/remes2000/amu_financial_summary/regexp"
	"net/http"
	"strconv"
	"testing"
)

func getCategories() []category.Category {
	return []category.Category{
		{Id: 1, Name: "First category", Regexps: []regexp.Regexp{
			{
				Id:      1,
				Content: "First regexp",
			},
			{
				Id:      2,
				Content: "Second regexp",
			},
		}},
		{Id: 2, Name: "Second category", Regexps: []regexp.Regexp{}},
	}
}

func prepareDatabaseToCategoryTest() {
	clearDatabase()
	categoriesToCreate := getCategories()
	global.Database.Create(&categoriesToCreate)
}

// CREATE CATEGORY

type createCategoryTest struct {
	Content        string
	ExpectedStatus int
	TestCaseName   string
}

var createCategoryTestCases = []createCategoryTest{
	{`{ "name": "", "regexps": [{"content": "^abc$"}] }`, http.StatusBadRequest, "Bad request when name blank"},
	{`{ "name": null, "regexps": [{"content": "^abc$"}] }`, http.StatusBadRequest, "Bad request when name null"},
	{`{ "name": "CategoryName", "regexps": null }`, http.StatusBadRequest, "Bad request when regexps null"},
	{`{ "name": "CategoryName", "regexps": [null] }`, http.StatusBadRequest, "Bad request when regexps array with null"},
	{`{ "name": "CategoryName", "regexps": [{}] }`, http.StatusBadRequest, "Bad request when regexps array with empty regexp"},
	{`{ "name": "CategoryName", "regexps": [{"content": "^abc$"}] }`, http.StatusOK, "Create category with single regexp"},
	{`{ "name": "CategoryName", "regexps": [] }`, http.StatusOK, "Create category when regexps array is empty"},
	{`{ "name": "CategoryName", "regexps": [{"content": "^abc$"}, {"content": "^abc$"}] }`, http.StatusOK, "Create category with multiple regexps"},
}

func Test_givenCreateCategoryRequest_whenCreateCategory_thenReturnExpectedStatusAndCreateCategory(t *testing.T) {
	for _, testCase := range createCategoryTestCases {
		clearDatabase()
		var apiResponseBody category.Category
		req, _ := http.NewRequest("POST", "/category", bytes.NewBuffer([]byte(testCase.Content)))
		response := executeRequest(req, &apiResponseBody)

		if testCase.ExpectedStatus != response.Code {
			t.Errorf("[%s] Expected response code %d. Got %d", testCase.TestCaseName, testCase.ExpectedStatus, response.Code)
			continue
		}

		if testCase.ExpectedStatus == http.StatusOK {
			assertCategoryCreated(t, testCase, apiResponseBody)
		}
	}
}

func assertCategoryCreated(t *testing.T, testCase createCategoryTest, apiResponseBody category.Category) {
	var sentCreateRequestBody category.CreateCategory
	var createdCategory category.Category
	json.Unmarshal([]byte(testCase.Content), &sentCreateRequestBody)

	for _, regexp := range apiResponseBody.Regexps {
		if regexp.CategoryID != 0 {
			t.Errorf("[%s] CategoryID should not be present in api response", testCase.TestCaseName)
			return
		}
	}

	if err := global.Database.Preload("Regexps").Where("id = ?", apiResponseBody.Id).First(&createdCategory).Error; err != nil {
		t.Errorf("[%s] Category is not present in database", testCase.TestCaseName)
		return
	}
	compareCategoriesCreate(t, testCase.TestCaseName, createdCategory, sentCreateRequestBody)
}

func compareCategoriesCreate(t *testing.T, testCaseName string, actual category.Category, expected category.CreateCategory) bool {
	if actual.Name != expected.Name {
		t.Errorf("[%s] Expected category with name %s. Got %s", testCaseName, expected.Name, actual.Name)
		return false
	}
	if len(actual.Regexps) != len(expected.Regexps) {
		t.Errorf("[%s] Expected number of regexps %d. Got %d", testCaseName, len(expected.Regexps), len(actual.Regexps))
		return false
	}
	for i := range expected.Regexps {
		expectedRegexp := expected.Regexps[i]
		actualRegexp := actual.Regexps[i]
		if compareRegexpsCreate(t, testCaseName, actualRegexp, expectedRegexp) == false {
			return false
		}
	}
	return true
}

func compareRegexpsCreate(t *testing.T, testCaseName string, actual regexp.Regexp, expected regexp.CreateRegexp) bool {
	if actual.Content != expected.Content {
		t.Errorf("[%s] Expected regexp with name %s. Got %s", testCaseName, expected.Content, actual.Content)
		return false
	}
	return true
}

// GET CATEGORY BY ID

type getByIdTest struct {
	Id             string
	ExpectedStatus int
	TestCaseName   string
}

var getByIdTestCases = []getByIdTest{
	{"0", http.StatusBadRequest, "Bad request when id 0"},
	{"random_string", http.StatusBadRequest, "Bad request when id string"},
	{"-5", http.StatusBadRequest, "Bad request when id lt 0"},
	{"1", http.StatusOK, "Return category with no regexps"},
	{"2", http.StatusOK, "Return category with multiple regexps"},
	{"3", http.StatusNotFound, "Not found when entity does not exist"},
}

func Test_givenCategoryId_whenGetCategoryById_thenReturnExpectedStatusAndExpectedBody(t *testing.T) {
	for _, testCase := range getByIdTestCases {
		prepareDatabaseToCategoryTest()
		var apiResponseBody category.Category
		req, _ := http.NewRequest("GET", "/category/"+testCase.Id, nil)
		response := executeRequest(req, &apiResponseBody)
		if testCase.ExpectedStatus != response.Code {
			t.Errorf("[%s] Expected response code %d. Got %d", testCase.TestCaseName, testCase.ExpectedStatus, response.Code)
			continue
		}
		if testCase.ExpectedStatus == http.StatusOK {
			categoryId, _ := strconv.Atoi(testCase.Id)
			expectedCategory := getCategories()[categoryId-1]
			compareCategories(t, testCase.TestCaseName, apiResponseBody, expectedCategory)
		}
	}
}

func compareCategories(t *testing.T, testCaseName string, actual category.Category, expected category.Category) bool {
	if actual.Id != expected.Id {
		t.Errorf("[%s] Expected category with id %d. Got %d", testCaseName, expected.Id, actual.Id)
		return false
	}
	if actual.Name != expected.Name {
		t.Errorf("[%s] Expected category with name %s. Got %s", testCaseName, expected.Name, actual.Name)
		return false
	}
	if len(actual.Regexps) != len(expected.Regexps) {
		t.Errorf("[%s] Expected number of regexps %d. Got %d", testCaseName, len(expected.Regexps), len(actual.Regexps))
		return false
	}
	for i := range expected.Regexps {
		expectedRegexp := expected.Regexps[i]
		actualRegexp := actual.Regexps[i]
		if compareRegexps(t, testCaseName, actualRegexp, expectedRegexp) == false {
			return false
		}
	}
	return true
}

func compareRegexps(t *testing.T, testCaseName string, actual regexp.Regexp, expected regexp.Regexp) bool {
	if actual.Id != expected.Id {
		t.Errorf("[%s] Expected regexp with id %d. Got %d", testCaseName, expected.Id, actual.Id)
		return false
	}
	if actual.Content != expected.Content {
		t.Errorf("[%s] Expected regexp with name %s. Got %s", testCaseName, expected.Content, actual.Content)
		return false
	}
	if actual.CategoryID != expected.CategoryID {
		t.Errorf("[%s] Expected regexp with category id %d. Got %d", testCaseName, expected.CategoryID, actual.CategoryID)
		return false
	}
	return true
}

// GET ALL CATEGORIES

type getAllTest struct {
	GetExpectedCategories func() []category.Category
	TestCaseName          string
}

var getAllTestCases = []getAllTest{
	{getCategories, "With full database"},
	{func() []category.Category { return []category.Category{} }, "With empty database"},
}

func Test_givenCategories_whenGetAllCategories_thenReturnAllCategories(t *testing.T) {
	for _, testCase := range getAllTestCases {
		clearDatabase()
		categoriesToCreate := testCase.GetExpectedCategories()
		global.Database.Create(&categoriesToCreate)

		var apiResponseBody []category.Category
		req, _ := http.NewRequest("GET", "/category", nil)
		response := executeRequest(req, &apiResponseBody)
		if http.StatusOK != response.Code {
			t.Errorf("[%s] Expected response code %d. Got %d", testCase.TestCaseName, http.StatusOK, response.Code)
			continue
		}

		if len(apiResponseBody) != len(testCase.GetExpectedCategories()) {
			t.Errorf("[%s] Expected number of categories %d. Got %d", testCase.TestCaseName, len(testCase.GetExpectedCategories()), len(apiResponseBody))
			return
		}

		for i := range apiResponseBody {
			actualRegexp := apiResponseBody[i]
			expectedRegexp := testCase.GetExpectedCategories()[i]
			if compareCategories(t, testCase.TestCaseName, actualRegexp, expectedRegexp) == false {
				return
			}
		}
	}
}

// UPDATE CATEGORY

type updateCategoryTest struct {
	Content        string
	ExpectedStatus int
	TestCaseName   string
}

var updateCategoryTestCases = []updateCategoryTest{
	{`{ "name": "CategoryName", "regexps": [{"content": "^abc$"}] }`, http.StatusBadRequest, "Bad request when id not provided"},
	{`{ "id": null, "name": "CategoryName", "regexps": [{"content": "^abc$"}] }`, http.StatusBadRequest, "Bad request when id null"},
	{`{ "id": 1, "name": "", "regexps": [{"content": "^abc$"}] }`, http.StatusBadRequest, "Bad request when name blank"},
	{`{ "id": 1, "name": null, "regexps": [{"content": "^abc$"}] }`, http.StatusBadRequest, "Bad request when name null"},
	{`{ "id": 1, "name": "CategoryName", "regexps": null }`, http.StatusBadRequest, "Bad request when regexps null"},
	{`{ "id": 1, "name": "CategoryName", "regexps": [null] }`, http.StatusBadRequest, "Bad request when regexps array with null"},
	{`{ "id": 1, "name": "CategoryName", "regexps": [{}] }`, http.StatusBadRequest, "Bad request when regexps array with empty regexp"},
	{`{ "id": 100, "name": "CategoryName", "regexps": [{"content": "^abc$"}] }`, http.StatusNotFound, "Not found when entity does not exist"},
	{`{ "id": 1, "name": "First category", "regexps": [{"content": "^abc$"}] }`, http.StatusOK, "Ok when change category name"},
	{`{ "id": 1, "name": "CategoryName", "regexps": [{"id": 1, "content": "First regexp"}, {"id": 2, "content": "Second regexp"}] }`, http.StatusOK, "Ok when replacing regexps"},
	{`{ "id": 1, "name": "CategoryName", "regexps": [{"id": 1, "content": "First regexp"}, {"id": 2, "content": "Second regexp"}, {"content": "Third regexp"}] }`, http.StatusOK, "Ok when inserting new regexp"},
	{`{ "id": 1, "name": "CategoryName", "regexps": [] }`, http.StatusOK, "Ok when deleting regexps"},
}

func Test_givenUpdateCategoryRequest_whenUpdateCategory_thenReturnExpectedStatusAndUpdateCategory(t *testing.T) {

	for _, testCase := range updateCategoryTestCases {
		prepareDatabaseToCategoryTest()

		var apiResponseBody category.Category
		req, _ := http.NewRequest("PUT", "/category", bytes.NewBuffer([]byte(testCase.Content)))
		response := executeRequest(req, &apiResponseBody)

		if testCase.ExpectedStatus != response.Code {
			t.Errorf("[%s] Expected response code %d. Got %d", testCase.TestCaseName, testCase.ExpectedStatus, response.Code)
			continue
		}
		if testCase.ExpectedStatus == http.StatusOK {
			if assertApiResponseAfterUpdate(t, testCase, apiResponseBody) == false {
				continue
			}
			if assertDbStateAfterUpdate(t, testCase, apiResponseBody) == false {
				continue
			}
		}
	}
}

func assertApiResponseAfterUpdate(t *testing.T, testCase updateCategoryTest, apiResponse category.Category) bool {
	var sentUpdateRequestBody category.Category
	json.Unmarshal([]byte(testCase.Content), &sentUpdateRequestBody)

	if apiResponse.Id != sentUpdateRequestBody.Id {
		t.Errorf("[%s] Expected category with id %d. Got %d", testCase.TestCaseName, sentUpdateRequestBody.Id, apiResponse.Id)
		return false
	}
	if apiResponse.Name != sentUpdateRequestBody.Name {
		t.Errorf("[%s] Expected category with name %s. Got %s", testCase.TestCaseName, sentUpdateRequestBody.Name, apiResponse.Name)
		return false
	}
	if len(apiResponse.Regexps) != len(sentUpdateRequestBody.Regexps) {
		t.Errorf("[%s] Expected number of regexps %d. Got %d", testCase.TestCaseName, len(sentUpdateRequestBody.Regexps), len(apiResponse.Regexps))
		return false
	}

	for i := range sentUpdateRequestBody.Regexps {
		expectedRegexp := sentUpdateRequestBody.Regexps[i]
		actualRegexp := apiResponse.Regexps[i]
		if assertApiResponseRegexpsAfterUpdate(t, testCase, actualRegexp, expectedRegexp) == false {
			return false
		}
	}
	return true
}

func assertApiResponseRegexpsAfterUpdate(t *testing.T, testCase updateCategoryTest, actual regexp.Regexp, expected regexp.Regexp) bool {
	if expected.Id == 0 && actual.Id == 0 {
		t.Errorf("[%s] Expected new regexp to got id", testCase.TestCaseName)
		return false
	}
	if expected.Content != actual.Content {
		t.Errorf("[%s] Expected regexp with name %s. Got %s", testCase.TestCaseName, expected.Content, actual.Content)
		return false
	}
	if expected.CategoryID != actual.CategoryID {
		t.Errorf("[%s] Expected regexp with category id %d. Got %d", testCase.TestCaseName, expected.CategoryID, actual.CategoryID)
		return false
	}
	return true
}

func assertDbStateAfterUpdate(t *testing.T, testCase updateCategoryTest, apiResponse category.Category) bool {
	var sentUpdateRequestBody category.Category
	json.Unmarshal([]byte(testCase.Content), &sentUpdateRequestBody)

	var categoriesInDb []category.Category
	global.Database.Preload("Regexps").Find(&categoriesInDb)

	var regexpsInDb []regexp.Regexp
	global.Database.Find(&regexpsInDb)

	if len(categoriesInDb) != len(getCategories()) {
		t.Errorf("[%s] Expected %d categories in database but found %d", testCase.TestCaseName, len(getCategories()), len(categoriesInDb))
		return false
	}

	numberOfRegexpsAfterPut := getNumberOfRegexpsAfterPut(getCategories(), sentUpdateRequestBody)
	if len(regexpsInDb) != numberOfRegexpsAfterPut {
		t.Errorf("[%s] Expected %d regexps in database but found %d", testCase.TestCaseName, numberOfRegexpsAfterPut, len(regexpsInDb))
		return false
	}

	var updatedCategoryFromDb category.Category
	for _, category := range categoriesInDb {
		if category.Id == sentUpdateRequestBody.Id {
			updatedCategoryFromDb = category
			break
		}
	}
	if updatedCategoryFromDb.Id == 0 {
		t.Errorf("[%s] Category with id %d does not exist in database", testCase.TestCaseName, sentUpdateRequestBody.Id)
		return false
	}

	if sentUpdateRequestBody.Name != updatedCategoryFromDb.Name {
		t.Errorf("[%s] Expected category with name %s. Got %s", testCase.TestCaseName, sentUpdateRequestBody.Name, updatedCategoryFromDb.Name)
		return false
	}

	if len(updatedCategoryFromDb.Regexps) != len(sentUpdateRequestBody.Regexps) {
		t.Errorf("[%s] Category with id %d expected %d regexps got %d", testCase.TestCaseName, sentUpdateRequestBody.Id, len(sentUpdateRequestBody.Regexps), len(updatedCategoryFromDb.Regexps))
		return false
	}

	for _, sentRegexp := range sentUpdateRequestBody.Regexps {
		if sentRegexp.Id != 0 {

			var updatedRegexp regexp.Regexp
			for _, regexpInDb := range updatedCategoryFromDb.Regexps {
				if regexpInDb.Id == sentRegexp.Id {
					updatedRegexp = regexpInDb
				}
			}
			if &updatedRegexp == nil {
				t.Errorf("[%s] Regexp with id %d is not attached to category", testCase.TestCaseName, sentRegexp.Id)
				return false
			}
			if updatedRegexp.Id != sentRegexp.Id {
				t.Errorf("[%s] Expected regexp with id %d. Got %d", testCase.TestCaseName, sentRegexp.Id, updatedRegexp.Id)
				return false
			}
			if updatedRegexp.Content != sentRegexp.Content {
				t.Errorf("[%s] Expected regexp with name %s. Got %s", testCase.TestCaseName, sentRegexp.Content, updatedRegexp.Content)
				return false
			}

		} else if isRegexpPresentInArray(regexpsInDb, sentRegexp.Content) == false {
			t.Errorf("[%s] Regexp with content %s is not attached to category", testCase.TestCaseName, sentRegexp.Content)
			return false
		}
	}

	return true
}

func isRegexpPresentInArray(regexps []regexp.Regexp, regexpContent string) bool {
	for _, regexp := range regexps {
		if regexp.Content == regexpContent {
			return true
		}
	}
	return false
}

func getNumberOfRegexpsAfterPut(dbState []category.Category, updateRequest category.Category) int {
	numberOfRegexpsOutsideUpdateRequest := 0
	for _, category := range dbState {
		if category.Id != updateRequest.Id {
			numberOfRegexpsOutsideUpdateRequest += len(category.Regexps)
		}
	}
	return numberOfRegexpsOutsideUpdateRequest + len(updateRequest.Regexps)
}
