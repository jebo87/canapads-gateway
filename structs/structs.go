package structs

import "gitlab.com/jebo87/makako-grpc/ads"

//Config struct
type Config struct {
	Gateway struct {
		Port string `yaml:"port"`
	} `yaml:"gateway"`
	API struct {
		ProdAddress string `yaml:"prod-address"`
		DevAddress  string `yaml:"dev-address"`
		Port        string `yaml:"port"`
	} `yaml:"api"`
}

//Exception struct
type Exception struct {
	Message string `json:"message"`
}

//ListingsCount struct
type ListingsCount struct {
	Count int `json:"count"`
}

//ListingsCount struct
type ListingID struct {
	ListingID int `json:"listingID"`
}

//ValidationResponse struct for the token validation
type ValidationResponse struct {
	Iss      string `json:"iss"`
	Nbf      int    `json:"nbf"`
	Exp      int    `json:"exp"`
	Aud      string `json:"api"`
	ClientID string `json:"client_id"`
	Sub      string `json:"sub"`
	AuthTime int    `json:"auth_time"`
	Idp      string `json:"idp"`
	Amr      string `json:"amr"`
	Active   bool   `json:"active"`
	Scope    string `json:"scope"`
}

//GrpcClient reusable client for GRPC connections
var GrpcClient ads.AdsClient
