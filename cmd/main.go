package main

import (
	"encoding/json"
	"github.com/greblin/smarthome/ya_sdk"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func DialogWrappedHandler(rw http.ResponseWriter, rq *http.Request) {
	var req ya_sdk.Request
	if err := json.NewDecoder(rq.Body).Decode(&req); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rsp, err := DialogHandler(rq.Context(), req)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(rsp); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	http.HandleFunc("/dialog", DialogWrappedHandler)
	http.HandleFunc("/kuzya", KuzyaHandler)
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
