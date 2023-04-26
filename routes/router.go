package routes

import (
	"database/sql"
	"mta-hosting-optimizer/models"
	"net/http"
)

func SetupRoutes(db *sql.DB) {

	http.HandleFunc("/inefficient_host", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			models.GetInefficientHost(w, r, db)
		case http.MethodPost:
			models.CreateHost(w, r, db)
		case http.MethodDelete:
			models.DeleteAllIPS(w, r, db)
		default:
			http.NotFound(w, r)
		}
	})

}
