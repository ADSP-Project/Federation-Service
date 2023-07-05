package main

import (
	"github.com/ADSP-Project/Federation-Service/handlers"
	"github.com/ADSP-Project/Federation-Service/federation"
	"fmt"
	"log"
	"net/http"
	"crypto/rsa"
	"os"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

var privKey *rsa.PrivateKey

var federationServer = "http://localhost:8000"

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run shop.go [port] [name]")
	}

	port := os.Args[1]
	shopName := os.Args[2]
	shopDescription := os.Args[3]

	router := mux.NewRouter()
	router.HandleFunc("/webhook", handlers.HandleWebhook).Methods("POST")
	router.HandleFunc("/api/v1/partnerships/request", func(w http.ResponseWriter, r *http.Request) {
		handlers.RequestPartnership(w, r, privKey)
	  }).Methods("POST")
	router.HandleFunc("/api/v1/partnerships/process", handlers.ProcessPartnership).Methods("POST")
	router.HandleFunc("/api/v1/partners", handlers.GetPartners).Methods("GET")


	privKey = federation.JoinFederation(shopName, shopDescription)
	go federation.PollFederationServer()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}