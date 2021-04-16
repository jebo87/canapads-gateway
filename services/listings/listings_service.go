package listings

import (
	"context"
	"strconv"

	"gitlab.com/jebo87/makako-gateway/clients"
	"gitlab.com/jebo87/makako-gateway/domain/grpc/listings"
	"gitlab.com/jebo87/makako-gateway/utils/errors"
	"gitlab.com/jebo87/makako-grpc/ads"
)

type listingsServiceInterface interface {
	GetListings(c context.Context, filter *ads.Filter) (*ads.AdList, *errors.RestErr)
	GetSingleListing(c context.Context, requestedID int64) (*ads.Ad, *errors.RestErr)
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

	return listings.GetListings(filter)

}

func (s *listingsService) GetSingleListing(c context.Context, requestedID int64) (*ads.Ad, *errors.RestErr) {

	searchResponse, err := clients.GrpcClient.AdDetail(c, &ads.Text{Text: strconv.FormatInt(requestedID, 10)})
	if err != nil {
		message := err.Error()
		restError := errors.NewServerError(message)
		return nil, restError
	}
	return searchResponse, nil
}
