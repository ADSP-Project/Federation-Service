package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Private key needs to be saved somewhere else for Production for example in .env or in a file
var globalPrivKey rsa.PrivateKey

type Shop struct {
	Name       string `json:"name"`
	WebhookURL string `json:"webhookURL"`
	PublicKey  string `json:"publicKey"`
}

type JWT struct {
	SigningString string
	Secret        string
}

type Partner struct {
	JWTToken          JWT
	Name              string `json:"name"`
	CanEarnCommission bool   `json:"comm"`
	CanShareInventory bool   `json:"shareInv"`
	CanShareData      bool   `json:"shareData"`
	CanCoPromote      bool   `json:"promote"`
	CanSell           bool   `json:"sell"`
	PublicKey         string `json:"publicKey"`
}

var federationServer = "http://localhost:8000"

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run shop.go [port] [name]")
	}

	port := os.Args[1]
	shopName := os.Args[2]

	router := mux.NewRouter()
	router.HandleFunc("/webhook", handleWebhook).Methods("POST")
	router.HandleFunc("/partner", handlePartner).Methods("POST")

	httpServer := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: router}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	joinFederation(shopName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}
	log.Printf("Shop has been created!")
}

func dbConn() (db *sql.DB) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbDriver := "postgres"
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName)

	db, err = sql.Open(dbDriver, dbInfo)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func GetJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "shop"
	claims["aud"] = "test"
	claims["iss"] = "shop"
	claims["exp"] = time.Now().Add(time.Minute * 120).Unix()

	tokenString, err := token.SignedString(globalPrivKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func sendRequest(webhookURL string, port string) {
	validToken, err := GetJWT()
	fmt.Println(validToken)
	if err != nil {
		fmt.Println("Failed to generate token")
	}
	jsonData, _ := json.Marshal(validToken)

	// resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	// http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Printf("Failed to send webhook to %s: %v\n", webhookURL, err)
		return
	}
	// defer resp.Body.Close()

}

func handlePartner(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	var partnerRequest Partner
	json.NewDecoder(r.Body).Decode(&partnerRequest)
	//TODO Check JWT validity

	//Inserting into database partners
	insForm, err := db.Prepare("INSERT INTO partners(shopName, canEarnCommission, canShareInventory, canShareData, canCoPromote, canSell, publicKey) VALUES($1,$2,$3,$4,$5,$6,$7)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(partnerRequest.Name, partnerRequest.CanEarnCommission, partnerRequest.CanShareInventory, partnerRequest.CanShareData, partnerRequest.CanCoPromote, partnerRequest.CanSell, partnerRequest.PublicKey)

}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	var newShop Shop
	json.NewDecoder(r.Body).Decode(&newShop)

	fmt.Printf("New shop joined the federation: %s\n", newShop.Name)

	fmt.Printf("Public Key: %s", newShop.PublicKey)
}

func exportPublicKeyAsPemStr(pubkey *rsa.PublicKey) string {
	PublicKey := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(pubkey)}))
	return PublicKey
}

func exportPrivateKeyAsPemStr(privatekey *rsa.PrivateKey) string {
	privatekey_pem := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}))
	return privatekey_pem
}

func joinFederation(shopName string) {

	privKey, err := rsa.GenerateKey(rand.Reader, 128)

	//for authentication in GetJWT func
	globalPrivKey = *privKey

	privatekey_pem := exportPrivateKeyAsPemStr(privKey)
	PublicKey := exportPublicKeyAsPemStr(&privKey.PublicKey)

	newShop := Shop{Name: shopName, WebhookURL: fmt.Sprintf("http://localhost:%s/webhook", os.Args[1]), PublicKey: PublicKey}

	log.Printf("New Shop Private Key is \n %s", privatekey_pem)
	log.Printf("New Shop Public key is \n %s", newShop.PublicKey)

	resp, err := http.PostForm("http://localhost:8081/login", url.Values{"name": {shopName}, "webhookURL": {newShop.WebhookURL}, "publicKey": {newShop.PublicKey}})
	if err != nil {
		log.Fatal("Failed to authenticate with auth server")
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	accessToken := result["access_token"]

	jsonData, _ := json.Marshal(newShop)
	req, err := http.NewRequest("POST", federationServer+"/shops", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", accessToken)

	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Failed to join federation: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Shop joined the federation")
	os.Exit(0)
}

func pollFederationServer() {
	db := dbConn()
	defer db.Close()

	for {
		time.Sleep(10 * time.Second)

		resp, err := http.Get(federationServer + "/shops")
		if err != nil {
			log.Printf("Failed to poll federation server: %v\n", err)
			continue
		}

		var shops []Shop
		json.NewDecoder(resp.Body).Decode(&shops)

		for _, shop := range shops {
			insForm, err := db.Prepare("INSERT INTO shops(name, webhookURL, publicKey) VALUES($1,$2,$3)")
			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(shop.Name, shop.WebhookURL, shop.PublicKey)
		}

		fmt.Printf("Current shops in the federation: %v\n", shops)
	}
}
