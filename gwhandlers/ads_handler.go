package gwhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-grpc/ads"
)

var itemsPerPage = 20

// ads?page=1&size=100

// AdsHandler handler for searches
func AdsHandler(w http.ResponseWriter, req *http.Request) {
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
	ads, err := list(context.Background(), structs.ClientGRPC, filter)
	if err != nil {
		log.Println(err)
		http.Error(w, "{}", http.StatusNotFound)

	} else {

		log.Println("printing ads in ServeHTTP for the Ads")
		w.Write(ads)
	}
}

func list(ctx context.Context, client ads.AdsClient, filter ads.Filter) ([]byte, error) {
	// ads, err := client.List(ctx, &ads.Void{})
	ads, err := client.List(ctx, &filter)
	if err != nil {
		return nil, fmt.Errorf("could not fetch ads: %v", err)
	}
	marshalled, err := json.Marshal(ads)
	//log.Println(ads)
	return marshalled, err

}
