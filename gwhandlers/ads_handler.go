package gwhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/ptypes/wrappers"
	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-grpc/ads"
)

var itemsPerPage = 20

// ads?page=1&qty=100

// AdsHandler handler for searches
func AdsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		log.Println("Options request")

		w.Header().Add("Access-Control-Allow-Methods", "GET")
		w.Header().Add("Access-Control-Allow-Origin", "https://www.canapads.ca")
		w.WriteHeader(http.StatusOK)

		return
	}
	log.Println("loading listings, request from ", req.RemoteAddr, " requesting ", req.URL)

	var pageCount, from, size int
	var errStr error
	filter := ads.Filter{}

	//https://gw.canapads.ca/ads?from=6&size=2&gym=true&furnished=true&pool=true&city=sutamarcha&gym=false&rent_by_owner=true&country=CAMBODIA&property_type=apartment&province=qc&neighborhood=la%20salle

	//check if the user is asking to show more listings per page, 100 maximum.
	if len(req.URL.Query()["qty"]) > 0 {
		size, errStr = strconv.Atoi(req.URL.Query()["qty"][0])
		//if the user is requesting unexpected listing quantities, set size to default
		if size > 100 || size < itemsPerPage {
			size = itemsPerPage
		}
		if errStr != nil {

			log.Println("Error trying to parse quantity of listings to show per page: ", errStr)
		} else {
			log.Println("Going back to default size")
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

	filter.From = &wrappers.Int32Value{Value: int32(from)}
	filter.Size = &wrappers.Int32Value{Value: int32(itemsPerPage)}

	if len(req.URL.Query()["gym"]) > 0 {
		gym, errStr := strconv.ParseBool(req.URL.Query()["gym"][0])
		if errStr != nil {
			filter.Gym = nil
			log.Println("Error trying to parse gym: ", errStr)
		} else {
			filter.Gym = &wrappers.BoolValue{Value: gym}
		}
	}
	if len(req.URL.Query()["pets"]) > 0 {
		pets, errStr := strconv.Atoi(req.URL.Query()["pets"][0])
		if errStr != nil || pets > 4 || pets < 0 {
			filter.Pets = nil
			log.Println("Error trying to parse gym: ", errStr)
		} else {
			filter.Pets = &wrappers.Int32Value{Value: int32(pets)}
		}
	}
	if len(req.URL.Query()["pool"]) > 0 {
		pool, errStr := strconv.ParseBool(req.URL.Query()["pool"][0])
		if errStr != nil {
			filter.Pool = nil
			log.Println("Error trying to parse pool: ", errStr)
		} else {
			filter.Pool = &wrappers.BoolValue{Value: pool}
		}
	}
	if len(req.URL.Query()["city"]) > 0 {
		filter.City = &wrappers.StringValue{Value: req.URL.Query()["city"][0]}
	}
	if len(req.URL.Query()["country"]) > 0 {
		//At this moment we default to Canada
		filter.Country = &wrappers.StringValue{Value: "Canada"}
		// filter.Country = &wrappers.StringValue{Value: req.URL.Query()["country"][0]}
	}
	if len(req.URL.Query()["property_type"]) > 0 {
		filter.PropertyType = &wrappers.StringValue{Value: req.URL.Query()["property_type"][0]}
	}
	if len(req.URL.Query()["furnished"]) > 0 {
		furnished, errStr := strconv.ParseBool(req.URL.Query()["furnished"][0])
		if errStr != nil {
			filter.Furnished = nil
			log.Println("Error trying to parse furnished: ", errStr)
		} else {
			filter.Furnished = &wrappers.BoolValue{Value: furnished}
		}
	}
	if len(req.URL.Query()["rent_by_owner"]) > 0 {
		rentByOwner, errStr := strconv.ParseBool(req.URL.Query()["rent_by_owner"][0])
		if errStr != nil {
			filter.RentByOwner = nil
			log.Println("Error trying to parse furnished: ", errStr)
		} else {
			filter.RentByOwner = &wrappers.BoolValue{Value: rentByOwner}
		}
	}
	if len(req.URL.Query()["province"]) > 0 {
		filter.StateProvince = &wrappers.StringValue{Value: req.URL.Query()["province"][0]}
	}
	if len(req.URL.Query()["neighborhood"]) > 0 {
		filter.Neighborhood = &wrappers.StringValue{Value: req.URL.Query()["neighborhood"][0]}
	}
	log.Println(req.URL.Query())
	if len(req.URL.Query()["search_param"]) > 0 {
		filter.SearchParam = &wrappers.StringValue{Value: req.URL.Query()["search_param"][0]}
	}

	//if k := keys[]

	listings, err := list(context.Background(), structs.ClientGRPC, filter)
	if err != nil {
		log.Println(err)
		http.Error(w, "{}", http.StatusNotFound)

	} else {

		log.Println("Success! - Returning listings to remote client")
		w.Write(listings)
	}
}

func list(ctx context.Context, client ads.AdsClient, filter ads.Filter) ([]byte, error) {
	// ads, err := client.List(ctx, &ads.Void{})
	ads, err := client.List(ctx, &filter)
	if err != nil {
		return nil, fmt.Errorf("could not fetch listings: %v", err)
	}
	marshalled, err := json.Marshal(ads)
	//log.Println(ads)
	return marshalled, err

}
