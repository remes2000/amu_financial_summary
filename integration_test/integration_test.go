package integration_test

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/remes2000/amu_financial_summary/app"
	"github.com/remes2000/amu_financial_summary/global"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env_test"); err != nil {
		panic("Error loading .env_test file")
	}
	app.Initialize()
	code := m.Run()
	os.Exit(code)
}

func clearDatabase() {
	global.Database.Exec("DELETE FROM regexps")
	global.Database.Exec("DELETE FROM categories")
}

func executeRequest(req *http.Request, body interface{}) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	global.Rest.ServeHTTP(responseRecorder, req)
	if body != nil {
		json.Unmarshal(responseRecorder.Body.Bytes(), body)
	}
	return responseRecorder
}
