package handlers

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ADSP-Project/Federation-Service/database"
	"github.com/ADSP-Project/Federation-Service/types"
	"github.com/golang-jwt/jwt/v5"
)

func ProcessPartnership(w http.ResponseWriter, r *http.Request) {
    var request types.PartnershipRequest

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)  // prints the actual request body

	// you need to put back the body content into the request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}
    
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		http.Error(w, "No Authorization header provided", http.StatusBadRequest)
		return
	}
	
	splitToken := strings.Split(authorizationHeader, "Bearer ")
	if len(splitToken) != 2 {
		http.Error(w, "Malformed Authorization header", http.StatusBadRequest)
		return
	}
	
	tokenString := splitToken[1] // Here is your token
	shopId := request.ShopId
    
    publicKeyStr, err := getPublicKeyFromDB(shopId)
    if err != nil {
		detailedError := fmt.Errorf("Failed to retrieve public key: %w", err)
		fmt.Printf("%+v\n", detailedError)
		http.Error(w, detailedError.Error(), http.StatusInternalServerError)
		return
	}

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
    
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		db := database.DbConn()
		shopId := claims["shopId"].(string)
		partnerId := claims["partnerId"].(string)

		partner, err := database.GetShopById(partnerId)
		if err != nil {
			log.Printf("Failed to get shop with id %s: %v\n", shopId, err)
			http.Error(w, "Failed to process partnership", http.StatusInternalServerError)
			return
		}
	
		// Insert new partnership entry
		sqlStatement := `
			INSERT INTO partners (shopid, shopname, canearncommission, canshareinventory, cansharedata, cancopromote, cansell, requeststatus)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending')
		`
		fmt.Print(request.Rights)
		_, err = db.Exec(sqlStatement, partnerId, partner.Name, request.Rights.CanEarnCommission, request.Rights.CanShareInventory, request.Rights.CanShareData, request.Rights.CanCoPromote, request.Rights.CanSell)
		if err != nil {
			log.Printf("Failed to insert new partnership: %v\n", err)
			http.Error(w, "Failed to process partnership", http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("Invalid token")
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	}
}

func RequestPartnership(w http.ResponseWriter, r *http.Request, privKey *rsa.PrivateKey) {
	var request types.PartnershipRequest
	json.NewDecoder(r.Body).Decode(&request)

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

	fmt.Printf("Generated JWT for partnership request: %s\n", tokenString)

	db := database.DbConn()
	defer db.Close()

	var partnerWebhookURL string
	err = db.QueryRow("SELECT webhookurl FROM shops WHERE id = $1", request.PartnerId).Scan(&partnerWebhookURL)
	if err != nil {
		http.Error(w, "Shop not found", http.StatusBadRequest)
		return
	}

	url, err := url.Parse(partnerWebhookURL)
	if err != nil {
		http.Error(w, "Error parsing the partner webhook URL", http.StatusInternalServerError)
		return
	}

	// removes the '/webhook' part
	url.Path = ""

	newURL := url.String() 

	jsonData, err := json.Marshal(request)
    if err != nil {
        http.Error(w, "Failed to create JSON body", http.StatusInternalServerError)
        return
    }

    req, err := http.NewRequest("POST", newURL+"/api/v1/partnerships/process", bytes.NewBuffer(jsonData))
    if err != nil {
        http.Error(w, "Failed to create request", http.StatusInternalServerError)
        return
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+tokenString)
	fmt.Printf("Parsed request: %+v\n", request)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to send partnership request", http.StatusInternalServerError)
		return
	}

	// Success
	fmt.Fprintln(w, "Partnership request successfully sent")
}

func getPublicKeyFromDB(shopId string) (string, error) {
    db := database.DbConn()
    defer db.Close()
    
    var publicKey string
    row := db.QueryRow("SELECT publicKey FROM shops WHERE id = $1", shopId)
    err := row.Scan(&publicKey)
    
    if err != nil {
        return "", fmt.Errorf("error getting public key from DB: %w", err)
    }
    
    return publicKey, nil
}