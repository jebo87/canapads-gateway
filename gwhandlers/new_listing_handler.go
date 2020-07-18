package gwhandlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"gitlab.com/jebo87/makako-gateway/httputils"
	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-grpc/ads"
)

func NewListingHandler(w http.ResponseWriter, req *http.Request) {
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

	req.Body = http.MaxBytesReader(w, req.Body, 524288)
	var listing ads.Ad
	err := httputils.DecodeJSONFromRequest(w, req, &listing, originAll)
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
	result, err := addListing(context.Background(), structs.ClientGRPC, listing)
	json.NewEncoder(w).Encode(structs.ListingID{ListingID: result})

}

func addListing(ctx context.Context, client ads.AdsClient, listing ads.Ad) (int, error) {

	result, err := client.AddListing(ctx, &listing)

	return int(result.ListingID), err

}
