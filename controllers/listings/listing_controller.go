package listings

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/jebo87/makako-gateway/services/listings"
	"gitlab.com/jebo87/makako-gateway/utils/errors"
	"gitlab.com/jebo87/makako-gateway/utils/utils_http"
	"gitlab.com/jebo87/makako-grpc/ads"
)

var origin string

//AdHandler handler for a single listing
func GetSingleListing(c *gin.Context) {

	requestedAd, found := c.Params.Get("id")

	if !found {
		restErr := errors.NewBadRequestError("id for listing required")
		c.AbortWithStatusJSON(restErr.Status, restErr)
		return
	}

	origin = utils_http.GetIP(c.Request)
	log.Println(fmt.Sprintf("[%v] requesting ad %v", origin, requestedAd))

	listingID, err := GetListingID(requestedAd)

	if err != nil {
		log.Println(err)
		restErr := errors.NewBadRequestError("Bad Request")
		c.AbortWithStatusJSON(restErr.Status, restErr)
		return
	}

	ad, err := listings.ListingsService.GetSingleListing(c, listingID)
	if err != nil {
		log.Println(err)
		restErr := errors.NewServerError("Internal server error")
		c.AbortWithStatusJSON(restErr.Status, restErr)
		return
	} else {
		c.JSON(http.StatusOK, ad)
	}
}

func adDetail(ctx context.Context, client ads.AdsClient, adID string) ([]byte, error) {
	searchParam := &ads.Text{}
	searchParam.Text = adID
	ad, err := client.AdDetail(ctx, searchParam)
	log.Println(fmt.Sprintf("[%v] Returning ad %v ", origin, adID))
	adFormatted, _ := json.Marshal(ad)
	return adFormatted, err

}

func GetListingID(id string) (int64, *errors.RestErr) {
	listingID, userErr := strconv.ParseInt(id, 10, 64)
	if userErr != nil {
		return 0, errors.NewBadRequestError("invalid listing ID")
	}

	return listingID, nil
}
