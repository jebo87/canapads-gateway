package gwhandlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-grpc/ads"
)

type responseUL struct {
	Response string
}

//UserListingsHandler handler used to return listings for a particular user
func UserListingsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		log.Printf("[%v] Options request", originAll)
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		//w.Header().Add("Access-Control-Allow-Origin", "https://www.canapads.ca")
		w.Header().Add("Access-Control-Allow-Origin", os.Getenv("ALLOWED_DOMAIN"))
		w.WriteHeader(http.StatusOK)

		return
	}
	log.Println(req)
	var userid []byte
	userid, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	log.Println(string(userid))
	userListings, err := getUserListings(context.Background(), structs.ClientGRPC, &ads.UserID{UserID: string(userid)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	listings, err := json.Marshal(userListings)
	//log.Println(httputils.JSONPrettyPrint(string(listings)))
	w.Write(listings)
}

func getUserListings(ctx context.Context, client ads.AdsClient, userID *ads.UserID) (*ads.AdList, error) {
	return client.UserListings(ctx, userID)

}
