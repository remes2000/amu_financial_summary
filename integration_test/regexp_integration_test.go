package integration_test

import (
	"github.com/remes2000/amu_financial_summary/global"
	"github.com/remes2000/amu_financial_summary/regexp"
	"log"
	"net/http"
	"strconv"
	"testing"
)

func Test_givenNonExistingRegexpId_whenGetSingleRegexp_then404(t *testing.T) {
	clearDatabase()
	//Given
	regexpId := strconv.Itoa(1)
	//When
	req, _ := http.NewRequest("GET", "/regexp/"+regexpId, nil)
	response := executeRequest(req, nil)
	//Then
	assertResponseCode(t, http.StatusNotFound, response.Code)
}

func Test_givenExistingRegexpId_whenGetSingleRegexp_thenReturnCorrectModel(t *testing.T) {
	clearDatabase()
	//Given
	regexpContent := "some-content"
	createdRegexp := createRegexp(regexpContent)
	regexpId := strconv.Itoa(int(createdRegexp.Id))
	//When
	var responseRegexp regexp.Regexp
	req, _ := http.NewRequest("GET", "/regexp/"+regexpId, nil)
	response := executeRequest(req, &responseRegexp)
	//Then
	assertResponseCode(t, http.StatusOK, response.Code)
	if responseRegexp.Content != createdRegexp.Content {
		t.Errorf("Expected regexp content to be '%s'. Got '%s'", regexpContent, responseRegexp.Content)
	}
	if responseRegexp.Id != createdRegexp.Id {
		t.Errorf("Expected regexp id to be '%d'. Bot '%d'", createdRegexp.Id, createdRegexp.Id)
	}
}

type createRegexpTest struct {
}

func Test_givenCreateRegexpModel_whenCreateRegexp_thenCreateOrAbort(t *testing.T) {
	clearDatabase()
	//Given

}

func createRegexp(c string) regexp.Regexp {
	newRegexp := regexp.Regexp{Content: c}
	if err := global.Database.Create(&newRegexp).Error; err != nil {
		log.Print(err)
		panic("Cannot create regexp")
	}
	return newRegexp
}
