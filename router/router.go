package router

import (
	"todo-list/middleware"
	"net/http"

	"github.com/gorilla/mux"
)
// gorilla/mux implements a request router and dispatcher for matching incoming requests to their respective handler
// Router is exported and used in main.go
func Router() *mux.Router {
	router := mux.NewRouter()

	// Serving static files
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	// API Endpoint
	// router.HandleFunc takes the api endpoint and match those with CRUD  VERBS
	router.HandleFunc("/", middleware.HomePage).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/todo", middleware.GetAllTodoList).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/todo", middleware.CreateTodoList).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/todo", middleware.UpdateTodoList).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/todo/{id}", middleware.DeleteTodoList).Methods("DELETE", "OPTIONS")

	return router
}
