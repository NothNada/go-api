package main

import (
	"go-api/db"
	"log"
	"net/http"

	"go-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB() // Inicializa o banco de dados

	go handlers.CheckForExpiredData()

	router := mux.NewRouter()

	// Rotas da API
	router.HandleFunc("/data", handlers.CreateData).Methods("POST")
	router.HandleFunc("/data", handlers.GetAllData).Methods("GET")
	router.HandleFunc("/data/{id}", handlers.GetData).Methods("GET")
	router.HandleFunc("/data/{id}", handlers.UpdateData).Methods("PUT")
	router.HandleFunc("/data/{id}", handlers.DeleteData).Methods("DELETE")

	log.Println("Servidor iniciado na porta :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
