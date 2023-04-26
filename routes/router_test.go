package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"mta-hosting-optimizer/models"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSetupRoutes(t *testing.T) {

	conn, err := sql.Open("mysql", "root:Nikhil58@@tcp(127.0.0.1:3306)/ips")
	if err != nil {
		fmt.Println("Cannot connect to database", err)
		return
	}

	ipConfigData := []byte(`[{
		"ip": "127.0.0.1",
		"hostname": "mta-prod-1",
		"active": true
	},
	{
		"ip": "127.0.0.2",
		"hostname": "mta-prod-1",
		"active": false
	},
	{
		"ip": "127.0.0.3",
		"hostname": "mta-prod-2",
		"active": true
	},
	{
		"ip": "127.0.0.4",
		"hostname": "mta-prod-2",
		"active": true
	},
	{
		"ip": "127.0.0.5",
		"hostname": "mta-prod-2",
		"active": false
	},
	{
		"ip": "127.0.0.6",
		"hostname": "mta-prod-3",
		"active": false
	}]`)

	req, err := http.NewRequest("POST", "/inefficient_host", bytes.NewBuffer(ipConfigData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		models.CreateHost(w, r, conn)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	req, err = http.NewRequest("GET", "/inefficient_host", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		models.GetInefficientHost(w, r, conn)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var ipconfig []models.IPConfig
	err = json.Unmarshal(rr.Body.Bytes(), &ipconfig)
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("DELETE", fmt.Sprintf("/inefficient_host"), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		models.DeleteAllIPS(w, r, conn)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
