package router

import (
	"backend/controllers"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/api/campaign", controllers.GetAllCampaigns).Methods("GET", "OPTIONS")


	return router
}