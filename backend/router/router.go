package router

import (
	"backend/controllers"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/api/campaign", controllers.GetAllCampaigns).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newcampaign", controllers.CreateCampaign).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deletecampaign/{id}", controllers.DeleteCampaign).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/updatecampaign/{id}", controllers.UpdateCampaign).Methods("PUT", "OPTIONS")
	
	return router
}