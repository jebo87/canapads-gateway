package listings

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-gateway/utils"
	"gitlab.com/jebo87/makako-gateway/utils/errors"
	"gitlab.com/jebo87/makako-grpc/ads"
	"google.golang.org/grpc/metadata"
)

type listingsServiceInterface interface {
	GetListings(c *gin.Context, filter ads.Filter) (*ads.AdList, *errors.RestErr)
}

type listingsService struct {
}

var (
	ListingsService listingsServiceInterface
)

func init() {
	ListingsService = &listingsService{}
}

func (s *listingsService) GetListings(c *gin.Context, filter ads.Filter) (*ads.AdList, *errors.RestErr) {

	metadata.AppendToOutgoingContext(c, "remote-addr", utils.GetIP(c.Request))

	searchResponse, err := structs.GrpcClient.List(c, &filter)
	if err != nil {
		log.Println("Error getting the listings from grpc server", err)
		restError := errors.NewServerError("invalid response from server")
		return nil, restError
	}
	marshalled, err := json.Marshal(searchResponse.List)
	result := &ads.AdList{}
	if err := json.Unmarshal(marshalled, &result); err != nil {
		log.Println("Error parsing the search result data from the server", err)
		restError := errors.NewServerError("invalid response from server")
		return nil, restError
	}
	return result, nil
}
