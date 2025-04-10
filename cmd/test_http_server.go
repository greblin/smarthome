package main

import (
	"encoding/json"
	"github.com/greblin/smarthome/ya_sdk"
	"log"
	"net/http"
)

func TestHandler(rw http.ResponseWriter, rq *http.Request) {
	var req ya_sdk.Request
	if err := json.NewDecoder(rq.Body).Decode(&req); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rsp, err := Handler(rq.Context(), req)
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
	http.HandleFunc("/", TestHandler)
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
