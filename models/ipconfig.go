package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Hostname struct {
	Hostname string `json:"hostname"`
}

type IPConfig struct {
	IP       string `json:"ip"`
	HostName string `json:"hostname"`
	Active   bool   `json:"active"`
}

func GetInefficientHost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	xStr := os.Getenv("THRESHOLDX")
	if xStr == "" {
		xStr = "1"
	}
	x, err := strconv.Atoi(xStr)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT Hostname FROM ipconfig  GROUP BY Hostname HAVING SUM(Active) <= ?", x)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []Hostname
	for rows.Next() {
		var h Hostname
		err := rows.Scan(&h.Hostname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, h)
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func CreateHost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var ipconfigs []IPConfig
	err := json.NewDecoder(r.Body).Decode(&ipconfigs)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, ipconfig := range ipconfigs {

		_, err := db.ExecContext(r.Context(), `
			INSERT INTO ipconfig (ip, hostname, active)
			VALUES (?, ?, ?)
		`, ipconfig.IP, ipconfig.HostName, ipconfig.Active)
		if err != nil {
			http.Error(w, "Failed to create ipconfig", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(ipconfig)
		if err != nil {
			http.Error(w, "Failed to encode ipconfig JSON", http.StatusInternalServerError)
			return
		}
	}
}

func DeleteAllIPS(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	stmt, err := db.Prepare("DELETE FROM ipconfig")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Deleted %d ips record(s)", rowsAffected)
}
