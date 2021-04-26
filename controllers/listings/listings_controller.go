package listings

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/jebo87/makako-gateway/services/listings"
	"gitlab.com/jebo87/makako-gateway/utils/errors"
	"gitlab.com/jebo87/makako-gateway/utils/utils_http"
	"gitlab.com/jebo87/makako-grpc/ads"
)

var originAll string

//example call
// ads?page=1&qty=5&gym=true&furnished=true&pool=true&city=montreal&gym=false
//     &rent_by_owner=true&country=canada&property_type=apartment&province=qc
//     &neighborhood=la%20salle&price_low=0&price_high=2000&search_param=metallica

//GetListings handler for searches
func GetListings(c *gin.Context) {
	var filter = &ads.Filter{}
	utils_http.LogDivider()

	// if utils_http.IsPreflight(c) {
	// 	c.String(http.StatusOK, "OK")
	// 	return
	// }
	//set maximum size for the request
	utils_http.SetMaxRqSize(c, 524288)

	log.Printf("[%v] Client connected", originAll)

	if err := c.ShouldBindJSON(filter); err != nil {
		log.Println(err)
		restError := errors.NewBadRequestError("invalid json body for filter")
		c.AbortWithStatusJSON(restError.Status, restError)
		return
	}

	a, _ := json.Marshal(filter)
	log.Println("listings_controller - Filter: ", utils_http.JSONPrettyPrint(string(a)))

	result, err := listings.ListingsService.GetListings(c, filter)

	if err != nil {
		c.AbortWithStatusJSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, result)

}

func logSuccess(result *ads.AdList) {
	log.Printf("[%v] Success! - Returning listings to remote client", originAll)
	log.Printf("[%v] finished", originAll)
	utils_http.LogDivider()
	log.Println(result)
}
