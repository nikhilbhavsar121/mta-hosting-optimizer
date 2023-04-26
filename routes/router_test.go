package routes

import (
	"encoding/json"
	"mta-hosting-optimizer/models"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSetupRoutes(t *testing.T) {

	req, err := http.NewRequest("GET", "/inefficient_host", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		models.GetInefficientHost1(w, r)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expectedResult := []string{"mta-prod-1", "mta-prod-3"}
	var actualResult []string
	if err := json.Unmarshal(rr.Body.Bytes(), &actualResult); err != nil {
		t.Errorf("error unmarshalling response body: %v", err)
	}

	sort.Strings(expectedResult)
	sort.Strings(actualResult)

	if !reflect.DeepEqual(actualResult, expectedResult) {
		t.Errorf("handler returned unexpected result: got %v want %v",
			actualResult, expectedResult)
	}

}
