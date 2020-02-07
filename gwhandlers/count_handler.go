package gwhandlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-grpc/ads"
)

//AdsCountHandler handler for ad_count
func AdsCountHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		log.Println("Options request")

		w.Header().Add("Access-Control-Allow-Methods", "GET")
		w.Header().Add("Access-Control-Allow-Origin", "https://www.canapads.ca")
		w.WriteHeader(http.StatusOK)

		return
	}
	log.Println("getting ad count from server")
	adCount, err := count(context.Background(), structs.ClientGRPC)

	if err != nil {
		log.Println(err)
		//json.NewEncoder(w).Encode(Exception{Message: err.Error()})
		http.Error(w, "{}", http.StatusNotAcceptable)
		return
	}

	json.NewEncoder(w).Encode(structs.AdJson{Count: adCount})
}

func count(ctx context.Context, client ads.AdsClient) (int, error) {
	adCount, err := client.Count(ctx, &ads.Void{})
	log.Println("Ad Count", adCount)
	return int(adCount.GetCount()), err
}
