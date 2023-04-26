package models

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type IPConfig struct {
	IP       string
	Hostname string
	Active   bool
}

type MockIPConfigService struct{}

func (s MockIPConfigService) GetIPConfig() []IPConfig {
	return []IPConfig{
		{"127.0.0.1", "mta-prod-1", true},
		{"127.0.0.2", "mta-prod-1", false},
		{"127.0.0.3", "mta-prod-2", true},
		{"127.0.0.4", "mta-prod-2", true},
		{"127.0.0.5", "mta-prod-2", false},
		{"127.0.0.6", "mta-prod-3", false},
	}
}

func GetInefficientHost1(w http.ResponseWriter, r *http.Request) {
	service := MockIPConfigService{}
	xStr := os.Getenv("THRESHOLDX")
	if xStr == "" {
		xStr = "1"
	}
	x, err := strconv.Atoi(xStr)
	if err != nil {
		log.Fatal(err)
	}

	// Get IP configuration data from mock service
	ipConfig := service.GetIPConfig()

	// Collect hostnames with less than or equal to X active IP addresses
	hostnames := make(map[string]int)
	for _, row := range ipConfig {

		if _, ok := hostnames[row.Hostname]; !ok && !row.Active {
			hostnames[row.Hostname] = 0
		} else if row.Active {
			hostnames[row.Hostname]++
		}

	}
	var result []string
	for hostname, activeIPs := range hostnames {
		if activeIPs <= x {
			result = append(result, hostname)
		}
	}

	// Marshal and write result as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
