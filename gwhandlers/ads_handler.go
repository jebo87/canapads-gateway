package gwhandlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/jebo87/makako-gateway/httputils"
	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-grpc/ads"
)

var originAll string

//example call
// ads?page=1&qty=5&gym=true&furnished=true&pool=true&city=montreal&gym=false
//     &rent_by_owner=true&country=canada&property_type=apartment&province=qc
//     &neighborhood=la%20salle&price_low=0&price_high=2000&search_param=metallica

//ListingHandler handler for searches
func ListingHandler(w http.ResponseWriter, req *http.Request) {
	httputils.LogDivider()
	originAll = httputils.GetIP(req)
	if req.Method == "OPTIONS" {
		log.Printf("[%v] Options request", originAll)
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		//w.Header().Add("Access-Control-Allow-Origin", "https://www.canapads.ca")
		w.Header().Add("Access-Control-Allow-Origin", os.Getenv("ALLOWED_DOMAIN"))
		w.WriteHeader(http.StatusOK)

		return
	}

	//set maximum size for the request
	req.Body = http.MaxBytesReader(w, req.Body, 524288)

	var filter ads.Filter

	log.Printf("[%v] Client connected", originAll)

	err := httputils.DecodeJSONFromRequest(w, req, &filter, originAll)
	if err != nil {
		var mr *httputils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Printf("[%v] %v", originAll, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}

	listings, err := list(context.Background(), structs.ClientGRPC, filter)
	if err != nil {
		log.Println(err)
		http.Error(w, "{}", http.StatusBadRequest)

	} else {

		log.Printf("[%v] Success! - Returning listings to remote client", originAll)
		log.Printf("[%v] finished", originAll)
		httputils.LogDivider()
		w.Write(listings)
	}

}

func list(ctx context.Context, client ads.AdsClient, filter ads.Filter) ([]byte, error) {
	// ads, err := client.List(ctx, &ads.Void{})
	ads, err := client.List(ctx, &filter)
	if err != nil {
		return nil, fmt.Errorf("[%v]could not fetch listings: %v", originAll, err)
	}
	marshalled, err := json.Marshal(ads)
	//log.Println(strconv.Itoa(int(ads.GetList().Ads)))
	return marshalled, err

}
