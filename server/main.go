package main

import (
	"fmt"
	"log"
	"net/http"

	"history-map/server/db"
	"history-map/server/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db.InitDB()
	router := mux.NewRouter()

	// Register Handlers
	router.HandleFunc("/maps", handlers.GetAllMaps).Methods("GET")
	router.HandleFunc("/maps/{id}", handlers.GetMapByID).Methods("GET")
	router.HandleFunc("/maps", handlers.CreateMap).Methods("POST")

	handler := cors.Default().Handler(router)

	fmt.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", handler))

}
