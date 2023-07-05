package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type Shop struct {
	Name       string `json:"name"`
	WebhookURL string `json:"webhookURL"`
	PublicKey string `json:"publicKey"`
	Description string `json:"description"`
}

type Partner struct {
	ShopId string   `json:"shopId"`
	ShopName string   `json:"shopName"`
	Rights struct {
		CanEarnCommission  bool `json:"canEarnCommission"`
		CanShareInventory  bool `json:"canShareInventory"`
		CanShareData       bool `json:"canShareData"`
		CanCoPromote       bool `json:"canCoPromote"`
		CanSell            bool `json:"canSell"`
	} `json:"rights"`
	RequestStatus string `json:"requestStatus"`
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
	Rights struct {
		CanEarnCommission  bool `json:"canEarnCommission"`
		CanShareInventory  bool `json:"canShareInventory"`
		CanShareData       bool `json:"canShareData"`
		CanCoPromote       bool `json:"canCoPromote"`
		CanSell            bool `json:"canSell"`
	} `json:"rights"`
}

type tokenClaims struct {
	ShopId     string   `json:"shopId"`
	PartnerId  string   `json:"partnerId"`
	Rights     []string `json:"rights"`
	jwt.RegisteredClaims
}