package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// ROUTER
	router := mux.NewRouter()

	// ROUTES
	router.HandleFunc("/", initServer).Methods("GET")

	// CORS
	handler := allowCORS(router)

	// SERVER
	fmt.Println("Server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}

func allowCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		handler.ServeHTTP(w, r)
	})
}

func initServer(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "<h1>Servidor Activo</h1>")
	if err != nil {
		return
	}
}
