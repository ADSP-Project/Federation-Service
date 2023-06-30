package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"strings"
)

var privKey *rsa.PrivateKey

type Shop struct {
	Name       string `json:"name"`
	WebhookURL string `json:"webhookURL"`
	PublicKey string `json:"publicKey"`
}

type ShopDisplay struct {
	Name       string `json:"name"`
	WebhookURL string `json:"webhookURL"`
}

type PartnershipRequest struct {
	ShopId     string   `json:"shopId"`
	PartnerId  string   `json:"partnerId"`
	Rights     []string `json:"rights"`
}

type PartnershipProcessRequest struct {
	ShopId string   `json:"shopId"`
	Jwt    string   `json:"jwt"`
	Rights []string `json:"rights"`
}

type tokenClaims struct {
	ShopId     string   `json:"shopId"`
	PartnerId  string   `json:"partnerId"`
	Rights     []string `json:"rights"`
	jwt.RegisteredClaims
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
	router.HandleFunc("/api/v1/partnerships/request", requestPartnership).Methods("POST")
	router.HandleFunc("/api/v1/partnerships/process", processPartnership).Methods("POST")


	privKey = joinFederation(shopName)
	go pollFederationServer()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
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


func handleWebhook(w http.ResponseWriter, r *http.Request) {
	var newShop Shop
	json.NewDecoder(r.Body).Decode(&newShop)

	fmt.Printf("New shop joined the federation: %s\n", newShop.Name)

	// fmt.Printf("Public Key: %s", newShop.PublicKey)
}

func requestPartnership(w http.ResponseWriter, r *http.Request) {
	var request PartnershipRequest
	json.NewDecoder(r.Body).Decode(&request)

	// Here, generate the JWT token with the rights embedded in it and send it to the partner.
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"shopId":    request.ShopId,
		"partnerId": request.PartnerId,
		"rights":    request.Rights,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(privKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while signing the token: %v", err), http.StatusInternalServerError)
		return
	}

	// For simplicity, I'm just printing the JWT
	fmt.Printf("Generated JWT for partnership request: %s\n", tokenString)
}

func processPartnership(w http.ResponseWriter, r *http.Request) {
    var request PartnershipProcessRequest
    json.NewDecoder(r.Body).Decode(&request)
    
    tokenString := request.Jwt
    shopId := request.ShopId
    
    publicKeyStr, err := getPublicKeyFromDB(shopId)
    if err != nil {
		detailedError := fmt.Errorf("Failed to retrieve public key: %w", err)
		fmt.Printf("%+v\n", detailedError)
		http.Error(w, detailedError.Error(), http.StatusInternalServerError)
		return
	}	

	fmt.Printf("PublicKeyStr: %s\n", publicKeyStr)

	publicKeyBlock, rest := pem.Decode([]byte(publicKeyStr))
	if publicKeyBlock == nil {
		detailedError := fmt.Errorf("Failed to decode public key. Remaining data: %s", string(rest))
		fmt.Printf("%+v\n", detailedError)
		http.Error(w, "Failed to decode public key", http.StatusInternalServerError)
		return
	}
	
    
    publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
    if err != nil {
        log.Printf("Failed to parse public key: %v\n", err)
        http.Error(w, "Failed to parse public key", http.StatusInternalServerError)
        return
    }
    
    // Now we can use publicKey in jwt.Parse
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return publicKey, nil
    })
    
    if err != nil {
        log.Printf("Error while validating the token: %v\n", err)
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }
    
    if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        // Process the partnership here
    } else {
        log.Println("Invalid token")
        http.Error(w, "Invalid token", http.StatusUnauthorized)
    }
}

func getPublicKeyFromDB(shopId string) (string, error) {
    db := dbConn()
    defer db.Close()
    
    var publicKey string
    row := db.QueryRow("SELECT publicKey FROM shops WHERE id = $1", shopId)
    err := row.Scan(&publicKey)
    
    if err != nil {
        return "", fmt.Errorf("error getting public key from DB: %w", err)
    }
    
    return publicKey, nil
}


func exportPublicKeyAsPemStr(pubkey *rsa.PublicKey) string {
    publicKeyBytes := x509.MarshalPKCS1PublicKey(pubkey)
    publicKeyPem := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: publicKeyBytes,
    })

    // Convert to string and remove the headers and footers.
    publicKeyPemStr := string(publicKeyPem)
    publicKeyPemStr = strings.Replace(publicKeyPemStr, "-----BEGIN RSA PUBLIC KEY-----\n", "", 1)
    publicKeyPemStr = strings.Replace(publicKeyPemStr, "\n-----END RSA PUBLIC KEY-----\n", "", 1)
    
    return publicKeyPemStr
}


func exportPrivateKeyAsPemStr(privatekey *rsa.PrivateKey) string {
	privatekey_pem := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}))
	return privatekey_pem
}

func joinFederation(shopName string) *rsa.PrivateKey {

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	// privatekey_pem := exportPrivateKeyAsPemStr(privKey)
	PublicKey := exportPublicKeyAsPemStr(&privKey.PublicKey)

	newShop := Shop{Name: shopName, WebhookURL: fmt.Sprintf("http://localhost:%s/webhook", os.Args[1]), PublicKey: PublicKey}

	// log.Printf("New Shop Private Key is \n %s", privatekey_pem)
	// log.Printf("New Shop Public key is \n %s", newShop.PublicKey)

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
		return privKey
	}
	defer resp.Body.Close()

	fmt.Println("Shop joined the federation")

	return privKey
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

		var shopsDisplay []ShopDisplay
		for _, shop := range shops {
			insForm, err := db.Prepare("INSERT INTO shops(name, webhookURL, publicKey) VALUES($1,$2,$3)")
			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(shop.Name, shop.WebhookURL, shop.PublicKey)
			shopsDisplay = append(shopsDisplay, ShopDisplay{Name: shop.Name, WebhookURL: shop.WebhookURL})
		}

		shopsDisplayJSON, _ := json.MarshalIndent(shopsDisplay, "", "    ")
		fmt.Printf("Current shops in the federation: \n%s\n", string(shopsDisplayJSON))
	}
}
