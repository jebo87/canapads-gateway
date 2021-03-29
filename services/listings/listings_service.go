package listings

import (
	"context"
	"encoding/json"
	"log"

	"gitlab.com/jebo87/makako-gateway/clients"
	"gitlab.com/jebo87/makako-gateway/utils/errors"
	"gitlab.com/jebo87/makako-grpc/ads"
)

type listingsServiceInterface interface {
	GetListings(c context.Context, filter *ads.Filter) (*ads.AdList, *errors.RestErr)
	GetSingleListing(c context.Context, requestedID string) (*ads.Ad, *errors.RestErr)
}

type listingsService struct {
}

var (
	ListingsService listingsServiceInterface
)

func init() {
	ListingsService = &listingsService{}
}

func (s *listingsService) GetListings(c context.Context, filter *ads.Filter) (*ads.AdList, *errors.RestErr) {

	//utils_http.AppendIPSourceToRequest(c)

	searchResponse, err := clients.GrpcClient.List(c, filter)
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

func (s *listingsService) GetSingleListing(c context.Context, requestedID string) (*ads.Ad, *errors.RestErr) {
	searchResponse, err := clients.GrpcClient.AdDetail(c, &ads.Text{Text: requestedID})
	if err != nil {
		message := err.Error()
		restError := errors.NewServerError(message)
		return nil, restError
	}
	return searchResponse, nil
}
