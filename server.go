package main

import (
	"employee-service/handlers"
	"employee-service/repo"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	employeeHandler := handlers.NewHandler(repo.NewSqlRepository("postgres://username:password@localhost:5432/employee-service"))
	router := mux.NewRouter()
	router.HandleFunc("/add", employeeHandler.AddHandler)
	router.HandleFunc("/get/{sort}", employeeHandler.GetHandler)

	http.Handle("/", router)

	fmt.Println("Слушаю шорохи...")
	_ = http.ListenAndServe(":80", nil)
}
