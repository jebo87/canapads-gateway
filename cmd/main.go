package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"bitbucket.org/jebo87/makako-grpc/ads"
	"google.golang.org/grpc"

	yaml "gopkg.in/yaml.v2"
)

//App struct
type App struct {
	AdsHandler *AdsHandler
}

//AdsHandler is the handler for all ads requests
type AdsHandler struct {
	AdHandler *AdHandler
}

//AdHandler is the handler for all ads requests
type AdHandler struct {
}

//Config struct
type Config struct {
	Gateway struct {
		Port string `yaml:"port"`
		URL  string `yaml:"url"`
	}
}

func main() {
	log.Println("Launching makako-gateway...")
	log.Println("Version 0.1")
	log.Println("Developed by Makako Labs http://www.makakolabs.ca")
	conf := loadConfig()
	app := &App{
		AdsHandler: new(AdsHandler),
	}

	go func() {
		http.ListenAndServe(conf.Gateway.URL+":"+conf.Gateway.Port, app)

	}()
	log.Println("Server started in http://" + conf.Gateway.URL + ":" + conf.Gateway.Port + ". Press CTRL+C to exit application")
	select {}

}

func loadConfig() (conf Config) {
	configFile, err := ioutil.ReadFile("../config/conf.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		panic(err)
	}
	return

}

func (h *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	enableCors(&res)
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	if head == "ads" {

		h.AdsHandler.ServeHTTP(res, req)
		return
	}
	http.Error(res, "Not Found", http.StatusNotFound)
}

//ServeHTTP for the Ads
func (h *AdsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	enableCors(&res)
	var head string
	head, _ = ShiftPath(req.URL.Path)
	//validate if there is an actual id
	if id, _ := strconv.Atoi(head); id != 0 {
		//if there is and ID then the AdHandler
		//should take care of bringing that specific ad
		h.AdHandler.ServeHTTP(res, req)
		return
	}

	//check if there is an offset and a limit in the query parameters.
	// offset, errOffset := strconv.Atoi(req.URL.Query().Get("offset"))
	// limit, errLimit := strconv.Atoi(req.URL.Query().Get("limit"))

	// //default to zero if offset or limit are not set
	// if errOffset != nil {
	// 	offset = 0
	// }
	// if errLimit != nil {
	// 	limit = 0
	// }

	switch req.Method {
	case "GET":
		fmt.Println("loading ads, request from " + req.RemoteAddr)
		//call gRPC server
		conn, err := grpc.Dial("localhost:7777", grpc.WithInsecure())
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not connect to backend: %v\n", err)
			os.Exit(1)
		}
		client := ads.NewAdsClient(conn)
		ads, err := list(context.Background(), client)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		log.Println("printing ads in ServeHTTP for the Ads")
		log.Println(string(ads))
		res.Write(ads)

	default:
		http.Error(res, "Only GET is allowed", http.StatusMethodNotAllowed)

	}
	return

}

func list(ctx context.Context, client ads.AdsClient) ([]byte, error) {
	ads, err := client.List(ctx, &ads.Void{})
	if err != nil {
		return nil, fmt.Errorf("could not fetch ads: %v", err)
	}
	log.Println("returning ad list")
	return json.Marshal(ads)

}

//ServeHTTP for one Ad
func (h *AdHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	switch req.Method {
	case "GET":
		fmt.Println("loading ad " + head + " request from " + req.RemoteAddr)
		//call gRPC server
		conn, err := grpc.Dial("localhost:7777", grpc.WithInsecure())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not connect to backend: %v\n", err)
			os.Exit(1) //TODO find a way to retry the connection with a timer
		}
		client := ads.NewAdsClient(conn)
		ad, err := adDetail(context.Background(), client, head)
		if err != nil {
			http.Error(res, "{}", http.StatusNotFound)
		} else {
			res.Write(ad)
		}
	default:
		http.Error(res, "Only GET is allowed", http.StatusMethodNotAllowed)

	}
	return
}

func adDetail(ctx context.Context, client ads.AdsClient, adId string) ([]byte, error) {
	searchParam := &ads.Text{}
	searchParam.Text = adId
	ad, err := client.AdDetail(ctx, searchParam)
	log.Println("Returning ad ", adId)
	adFormatted, _ := json.Marshal(ad)
	return adFormatted, err

}

//ShiftPath returns the head of the URL without initial slash '/' and the rest of the URL
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
}
