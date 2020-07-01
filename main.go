package main

import (
	"flag"
	"fmt"
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
	"gitlab.com/jebo87/makako-gateway/httputils"
	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-grpc/ads"
	"google.golang.org/grpc"
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

//Environment variables to set:
// os.Getenv("API_ADDRESS_PROD"), os.Getenv("API_PORT_PROD")
// os.Getenv("API_ADDRESS_DEV"), os.Getenv("API_PORT_DEV")
// os.Getenv("PORT")
// os.Getenv("ELASTIC_ADDRESS"), os.Getenv("ELASTIC_PORT"))
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
		conn, err = grpc.Dial(os.Getenv("API_ADDRESS_PROD")+":"+os.Getenv("API_PORT_PROD"), grpc.WithInsecure())
		log.Printf("connecting to GRPC server in %v:%v ", os.Getenv("API_ADDRESS_PROD"), os.Getenv("API_PORT_PROD"))

	} else {
		conn, err = grpc.Dial(os.Getenv("API_ADDRESS_DEV")+":"+os.Getenv("API_PORT_DEV"), grpc.WithInsecure())

		log.Println("connecting to GRPC server " + os.Getenv("API_ADDRESS_DEV"))

	}

	defer conn.Close()
	clientGRPC = ads.NewAdsClient(conn)
	structs.ClientGRPC = clientGRPC

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %v\n", err)
		os.Exit(1)
	}
	//TODO: add a function to call on the gRPC server to make sure it is working.

	<-c

}

func loadConfig() {
	//TODO: check if we really need to do this
	deployedFlag = flag.Bool("deployed", false, "Defines if absolute paths need to be used for the config files")
	flag.Parse()
}

func loadHandlers() {
	router = mux.NewRouter()

	// router.HandleFunc("/ads", httputils.ValidateMiddleware(adsHandler)).Methods("GET", "OPTIONS")
	// router.HandleFunc("/ads", (gwhandlers.AdsHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/ads", (gwhandlers.ListingHandler)).Methods("GET", "POST", "OPTIONS")
	router.HandleFunc("/ads/{key}", gwhandlers.AdHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/ad_count", httputils.ValidateMiddleware(gwhandlers.AdsCountHandler)).Methods("GET", "OPTIONS")
}

func startServer(router *mux.Router) {
	methodsOk := handlers.AllowedMethods([]string{"GET", "OPTIONS"})
	//allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"})
	originsOK := handlers.AllowedOrigins([]string{"https://www.canapads.ca", "http://192.168.2.201:30030"})
	optionsOk := handlers.IgnoreOptions()
	// log.Fatal(http.ListenAndServe(":"+conf.Gateway.Port, handlers.CORS(methodsOk, originsOK)(router)))

	if *deployedFlag {
		log.Println("Server started in prod mode @ http://0.0.0.0" + ":" + os.Getenv("PORT") + ". Press CTRL+C to exit application")
		log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handlers.CORS(optionsOk, methodsOk, originsOK)(router)))
	} else {
		log.Println("Server started in dev mode @ http://localhost" + ":" + os.Getenv("PORT") + ". Press CTRL+C to exit application")
		log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("PORT"), handlers.CORS(optionsOk, methodsOk, originsOK)(router)))
	}
}

func loadElastic() {
	log.Printf("connecting to ElasticSearch in %v:%v", os.Getenv("ELASTIC_ADDRESS"), os.Getenv("ELASTIC_PORT"))
	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%v:%v", os.Getenv("ELASTIC_ADDRESS"), os.Getenv("ELASTIC_PORT")),
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
