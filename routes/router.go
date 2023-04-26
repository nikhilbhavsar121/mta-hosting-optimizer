package routes

import (
	"mta-hosting-optimizer/models"
	"net/http"
)

func SetupRoutes() {

	http.HandleFunc("/inefficient_host", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			models.GetInefficientHost1(w, r)
		default:
			http.NotFound(w, r)
		}
	})

}
