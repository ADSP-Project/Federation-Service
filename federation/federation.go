package federation

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	mathrand "math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ADSP-Project/Federation-Service/database"
	"github.com/ADSP-Project/Federation-Service/globals"
	"github.com/ADSP-Project/Federation-Service/types"
	_ "github.com/lib/pq"
)

func ExportPublicKeyAsPemStr(pubkey *rsa.PublicKey) string {
	PublicKey := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(pubkey)}))
	return PublicKey
}


func ExportPrivateKeyAsPemStr(privatekey *rsa.PrivateKey) string {
	privatekey_pem := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}))
	return privatekey_pem
}

func JoinFederation(shopName string, shopDescription string) *rsa.PrivateKey {

	seed, err := strconv.Atoi(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

	randomSource := mathrand.New(mathrand.NewSource(int64(seed)))

	privKey, err := rsa.GenerateKey(randomSource, 2048)
	// ExportPrivateKeyAsPemStr(privKey)
	PublicKey := ExportPublicKeyAsPemStr(&privKey.PublicKey)

	newShop := types.Shop{
		Name:        shopName, 
		WebhookURL:  fmt.Sprintf("http://localhost:%s/webhook", os.Args[1]), 
		PublicKey:   PublicKey, 
		Description: shopDescription,
	}

	// log.Printf("New Shop Private Key is \n %s", privatekey_pem)
	// log.Printf("New Shop Public key is \n %s", newShop.PublicKey)

	resp, err := http.PostForm(globals.AuthServer+"/login", url.Values{"name": {shopName}, "webhookURL": {newShop.WebhookURL}, "publicKey": {newShop.PublicKey}})
	if err != nil {
		log.Fatal("Failed to authenticate with auth server")
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	accessToken := result["access_token"]

	jsonData, _ := json.Marshal(newShop)
	req, err := http.NewRequest("POST", globals.FederationServer+"/shops", bytes.NewBuffer(jsonData))
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

func PollFederationServer() {
	db := database.DbConn()
	defer db.Close()

	for {
		time.Sleep(10 * time.Second)

		resp, err := http.Get(globals.FederationServer + "/shops")
		if err != nil {
			log.Printf("Failed to poll federation server: %v\n", err)
			continue
		}

		var shops []types.Shop
		json.NewDecoder(resp.Body).Decode(&shops)

		var shopsDisplay []types.ShopDisplay
		for _, shop := range shops {
			insForm, err := db.Prepare("INSERT INTO shops(name, webhookURL, publicKey, description) VALUES($1,$2,$3,$4)") // modify this line
			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(shop.Name, shop.WebhookURL, shop.PublicKey, shop.Description)
			shopsDisplay = append(shopsDisplay, types.ShopDisplay{Name: shop.Name, WebhookURL: shop.WebhookURL})
		}

		// shopsDisplayJSON, _ := json.MarshalIndent(shopsDisplay, "", "    ")
		// fmt.Printf("Current shops in the federation: \n%s\n", string(shopsDisplayJSON))
	}
}
