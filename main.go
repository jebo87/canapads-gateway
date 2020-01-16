package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gitlab.com/jebo87/makako-grpc/ads"
	"google.golang.org/grpc"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	yaml "gopkg.in/yaml.v2"
)

//Config struct
type Config struct {
	Gateway struct {
		Port string `yaml:"port"`
	} `yaml:gateway`
	API struct {
		ProdAddress string `yaml:"prod-address"`
		DevAddress  string `yaml:"dev-address"`
		Port        string `yaml:"port"`
	} `yaml:api`
}

type Exception struct {
	Message string `json:"message"`
}
type AdJson struct {
	Count int `json:"count"`
}

var deployedFlag *bool
var conf Config
var conn *grpc.ClientConn
var client ads.AdsClient
var itemsPerPage = 20

//ContextKey used in context
type ContextKey string

//ContextDecodeKey key for the context
const ContextDecodeKey ContextKey = "decoded"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	conf = loadConfig()
	// app := &App{
	// 	AdsHandler: new(AdsHandler),
	// }
	router := mux.NewRouter()
	methodsOk := handlers.AllowedMethods([]string{"GET", "OPTIONS"})
	//allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"})
	originsOK := handlers.AllowedOrigins([]string{"https://www.canapads.ca"})
	optionsOk := handlers.IgnoreOptions()

	log.Println("Launching makako-gateway...")
	log.Println("Version 0.13")
	log.Println("Developed by Makako Labs http://www.makakolabs.ca")
	router.HandleFunc("/ads", (adsHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/ads/{key}", adHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/ad_count", adsCountHandler).Methods("GET", "OPTIONS")
	//router.HandleFunc("/ads", (testEndpoint)).Methods("GET")

	go func() {
		// log.Fatal(http.ListenAndServe(":"+conf.Gateway.Port, handlers.CORS(methodsOk, originsOK)(router)))

		if *deployedFlag {
			log.Println("Starting server in production mode...")
			log.Println("Server started in http://0.0.0.0" + ":" + conf.Gateway.Port + ". Press CTRL+C to exit application")
			log.Fatal(http.ListenAndServe(":"+conf.Gateway.Port, handlers.CORS(optionsOk, methodsOk, originsOK)(router)))

		} else {
			log.Println("Starting server in develpment mode")
			log.Println("Server started in http://localhost" + ":" + conf.Gateway.Port + ". Press CTRL+C to exit application")
			log.Fatal(http.ListenAndServe("localhost:"+conf.Gateway.Port, handlers.CORS(optionsOk, methodsOk, originsOK)(router)))
		}

	}()

	//call gRPC server

	var err error
	if *deployedFlag {
		conn, err = grpc.Dial(conf.API.ProdAddress+":"+conf.API.Port, grpc.WithInsecure())
		log.Println("connecting to GRPC server " + conf.API.ProdAddress)

	} else {
		conn, err = grpc.Dial(conf.API.DevAddress+":"+conf.API.Port, grpc.WithInsecure())

		log.Println("connecting to GRPC server " + conf.API.DevAddress)

	}
	defer conn.Close()
	client = ads.NewAdsClient(conn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %v\n", err)
		os.Exit(1)
	}

	<-c

}

func loadConfig() (conf Config) {
	deployedFlag = flag.Bool("deployed", false, "Defines if absolute paths need to be used for the config files")
	var configFile []byte
	var err error
	flag.Parse()

	if *deployedFlag {
		configFile, err = ioutil.ReadFile("/makako-gateway/bin/config/conf.yaml")
	} else {
		configFile, err = ioutil.ReadFile("config/conf.yaml")
	}
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		panic(err)
	}
	return

}

func adHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		log.Println("Options request")

		res.Header().Add("Access-Control-Allow-Methods", "GET")
		//res.Header().Add("Access-Control-Allow-Headers", "Authorization")
		res.Header().Add("Access-Control-Allow-Origin", "https://www.canapads.ca")
		res.WriteHeader(http.StatusOK)

		return
	}
	requestedAd := req.URL.Path[5:]

	ad, err := adDetail(context.Background(), client, requestedAd)
	//in case we need to see the ad returned, uncomment the three following lines
	// var dat ads.Ad
	// json.Unmarshal(ad, &dat)
	// log.Println(dat)
	if err != nil {
		log.Println(err)

		http.Error(res, "{}", http.StatusNotFound)
	} else {
		res.Write(ad)
	}
}

func adsCountHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		log.Println("Options request")

		w.Header().Add("Access-Control-Allow-Methods", "GET")
		w.Header().Add("Access-Control-Allow-Origin", "https://www.canapads.ca")
		w.WriteHeader(http.StatusOK)

		return
	}
	log.Println("getting ad count from server")
	adCount, err := count(context.Background(), client)

	if err != nil {
		log.Println(err)
		//json.NewEncoder(w).Encode(Exception{Message: err.Error()})
		http.Error(w, "{}", http.StatusNotAcceptable)
		return
	}

	json.NewEncoder(w).Encode(AdJson{Count: adCount})
}

func adsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		log.Println("Options request")

		w.Header().Add("Access-Control-Allow-Methods", "GET")
		w.Header().Add("Access-Control-Allow-Origin", "https://www.canapads.ca")
		w.WriteHeader(http.StatusOK)

		return
	}

	log.Println("loading ads, request from " + req.RemoteAddr)
	var pageCount, from, size int
	var errStr error

	//check if the user is asking to show more listings per page, 100 maximum.
	if len(req.URL.Query()["qty"]) > 0 {
		size, errStr = strconv.Atoi(req.URL.Query()["qty"][0])
		//if the user is requesting unexpected listing quantities, set size to default
		if size > 100 || size < itemsPerPage {
			size = itemsPerPage
		}
		if errStr != nil {
			log.Println("Error trying to parse page number")
		} else {
			itemsPerPage = size
		}

	}
	if len(req.URL.Query()["page"]) > 0 {
		pageCount, errStr = strconv.Atoi(req.URL.Query()["page"][0])
		from = pageCount*itemsPerPage - itemsPerPage
		if errStr != nil {
			from = 0
		}
	} else {
		from = 0
	}

	filter := ads.Filter{
		From: int32(from),
		Size: int32(itemsPerPage)}
	ads, err := list(context.Background(), client, filter)
	if err != nil {
		log.Println(err)
		http.Error(w, "{}", http.StatusNotFound)

	} else {

		log.Println("printing ads in ServeHTTP for the Ads")
		w.Write(ads)
	}
}

//This is the middleware used to protect certain api calls
func validateMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//enableCors(&w)
		if req.Method == "OPTIONS" {
			log.Println("Options request")

			w.Header().Add("Access-Control-Allow-Methods", "GET")
			w.Header().Add("Access-Control-Allow-Headers", "Authorization")
			w.Header().Add("Access-Control-Allow-Origin", "https://www.canapads.ca")
			w.WriteHeader(http.StatusOK)

			return
		}
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				if bearerToken[0] == "Bearer" {
					log.Println(bearerToken)
					if validateToken(bearerToken[1]) {
						//everything is ok, proceed with allow the exectution of the next function
						next(w, req)
					} else {
						json.NewEncoder(w).Encode(Exception{Message: "Token invalid or expired"})
					}
				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
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

func validateToken(token string) bool {
	form := url.Values{}
	form.Add("token", token)
	log.Println("validating token " + token)
	req, erro := http.NewRequest("POST", "http://localhost:4445/oauth2/introspect", strings.NewReader(form.Encode()))
	if erro != nil {
		log.Println(erro.Error())

		return false

	}
	log.Println("adding headers")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Add("Authorization", "Basic "+basicAuth("api1", "TopSecret2"))
	resp, err := netClient.Do(req)

	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	var validationObject ValidationResponse
	json.Unmarshal(bodyBytes, &validationObject)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Validation for the request has returned %v \n", validationObject.Active)

	return validationObject.Active

}

func list(ctx context.Context, client ads.AdsClient, filter ads.Filter) ([]byte, error) {
	// ads, err := client.List(ctx, &ads.Void{})
	ads, err := client.List(ctx, &filter)
	if err != nil {
		return nil, fmt.Errorf("could not fetch ads: %v", err)
	}
	marshalled, err := json.Marshal(ads)
	log.Println(ads)
	return marshalled, err

}

func adDetail(ctx context.Context, client ads.AdsClient, adID string) ([]byte, error) {
	searchParam := &ads.Text{}
	searchParam.Text = adID
	ad, err := client.AdDetail(ctx, searchParam)
	log.Println("Returning ad ", adID)
	adFormatted, _ := json.Marshal(ad)
	return adFormatted, err

}

func count(ctx context.Context, client ads.AdsClient) (int, error) {
	adCount, err := client.Count(ctx, &ads.Void{})
	log.Println("Ad Count", adCount)
	return int(adCount.GetCount()), err
}
