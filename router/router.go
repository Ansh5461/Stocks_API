package router

import (
	"github.com/gorilla/mux"

	"Stocks_API/middleware"
)

func Router() *mux.Router {

	//router is created
	router := mux.NewRouter()

	//we will have 5 commands, so 5 handle func

	//this is for get by ID, so api/stock/ id which will be provided
	router.HandleFunc("/api/stock/{id}", middleware.GetStock).Methods("GET", "OPTIONS")

	//this is for getting all stocks, so only api/stocks
	router.HandleFunc("/api/stock", middleware.GetAllStock).Methods("GET", "OPTIONS")

	//for creating new stock
	router.HandleFunc("/api/newstock", middleware.CreateStock).Methods("POST", "OPTIONS")

	//for deleting
	router.HandleFunc("/api/stock/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")

	//for updating
	router.HandleFunc("/api/stock/{id}", middleware.UpdateStock).Methods("PUT", "OPTIONS")

	return router

}
