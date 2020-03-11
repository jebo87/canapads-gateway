package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gitlab.com/jebo87/makako-gateway/gwhandlers"
	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-grpc/ads"
	"google.golang.org/grpc"

	yaml "gopkg.in/yaml.v2"
)

var deployedFlag *bool
var conf structs.Config
var conn *grpc.ClientConn
var clientGRPC ads.AdsClient
var netClient = &http.Client{
	Timeout: time.Second * 10,
}
var router *mux.Router
var es elasticsearch.Client

func main() {
	log.Println("Launching makako-gateway...")
	log.Println("Developed by Makako Labs http://www.makakolabs.ca")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	loadConfig()
	loadHandlers()
	loadElastic()
	go startServer(router)

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
	clientGRPC = ads.NewAdsClient(conn)
	structs.ClientGRPC = clientGRPC

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %v\n", err)
		os.Exit(1)
	}
	log.Println("connected to GRPC server")

	<-c

}

func loadConfig() {
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

func loadHandlers() {
	router = mux.NewRouter()

	// router.HandleFunc("/ads", httputils.ValidateMiddleware(adsHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/ads", (gwhandlers.AdsHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/ads/{key}", gwhandlers.AdHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/ad_count", gwhandlers.AdsCountHandler).Methods("GET", "OPTIONS")
}

func startServer(router *mux.Router) {
	methodsOk := handlers.AllowedMethods([]string{"GET", "OPTIONS"})
	//allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"})
	originsOK := handlers.AllowedOrigins([]string{"https://www.canapads.ca"})
	optionsOk := handlers.IgnoreOptions()
	// log.Fatal(http.ListenAndServe(":"+conf.Gateway.Port, handlers.CORS(methodsOk, originsOK)(router)))

	if *deployedFlag {
		log.Println("Server started in prod mode @ http://0.0.0.0" + ":" + conf.Gateway.Port + ". Press CTRL+C to exit application")
		log.Fatal(http.ListenAndServe(":"+conf.Gateway.Port, handlers.CORS(optionsOk, methodsOk, originsOK)(router)))
	} else {
		log.Println("Server started in dev mode @ http://localhost" + ":" + conf.Gateway.Port + ". Press CTRL+C to exit application")
		log.Fatal(http.ListenAndServe("localhost:"+conf.Gateway.Port, handlers.CORS(optionsOk, methodsOk, originsOK)(router)))
	}
}

func loadElastic() {
	log.Println("connecting to ElasticSearch...")
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://192.168.2.200:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	_, err = es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	log.Println("Connected to ElasticSearch")
}
