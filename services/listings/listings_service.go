package listings

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-gateway/utils/errors"
	"gitlab.com/jebo87/makako-gateway/utils/utils_http"
	"gitlab.com/jebo87/makako-grpc/ads"
)

type listingsServiceInterface interface {
	GetListings(c *gin.Context, filter *ads.Filter) (*ads.AdList, *errors.RestErr)
	GetSingleListing(c *gin.Context, requestedID string) (*ads.Ad, *errors.RestErr)
}

type listingsService struct {
}

var (
	ListingsService listingsServiceInterface
)

func init() {
	ListingsService = &listingsService{}
}

func (s *listingsService) GetListings(c *gin.Context, filter *ads.Filter) (*ads.AdList, *errors.RestErr) {

	utils_http.AppendIPSourceToRequest(c)

	searchResponse, err := structs.GrpcClient.List(c, filter)
	if err != nil {
		log.Println("Invalid response from grpc server", err)
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

func (s *listingsService) GetSingleListing(c *gin.Context, requestedID string) (*ads.Ad, *errors.RestErr) {
	searchResponse, err := structs.GrpcClient.AdDetail(context.Background(), &ads.Text{Text: requestedID})
	if err != nil {
		log.Println("Invalid response from grpc server", err)
		restError := errors.NewServerError("invalid response from server")
		return nil, restError
	}
	return searchResponse, nil
}
