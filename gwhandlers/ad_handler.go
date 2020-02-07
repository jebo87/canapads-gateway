package gwhandlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-grpc/ads"
)

//AdHandler handler for a single listing
func AdHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		log.Println("Options request")

		res.Header().Add("Access-Control-Allow-Methods", "GET")
		//res.Header().Add("Access-Control-Allow-Headers", "Authorization")
		res.Header().Add("Access-Control-Allow-Origin", "https://www.canapads.ca")
		res.WriteHeader(http.StatusOK)

		return
	}
	requestedAd := req.URL.Path[5:]

	ad, err := adDetail(context.Background(), structs.ClientGRPC, requestedAd)
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

func adDetail(ctx context.Context, client ads.AdsClient, adID string) ([]byte, error) {
	searchParam := &ads.Text{}
	searchParam.Text = adID
	ad, err := client.AdDetail(ctx, searchParam)
	log.Println("Returning ad ", adID)
	adFormatted, _ := json.Marshal(ad)
	return adFormatted, err

}
