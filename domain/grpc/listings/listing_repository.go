package listings

import (
	"context"
	"encoding/json"
	"log"

	"gitlab.com/jebo87/makako-gateway/clients"
	"gitlab.com/jebo87/makako-gateway/utils/errors"
	"gitlab.com/jebo87/makako-grpc/ads"
)

func GetListings(filter *ads.Filter) (*ads.SearchResponse, *errors.RestErr) {
	searchResponse, err := clients.GrpcClient.List(context.Background(), filter)
	if err != nil {
		log.Println("Invalid response from grpc server", err)
		restError := errors.NewServerError("invalid response from server")
		return nil, restError
	}
	marshalled, err := json.Marshal(searchResponse)

	if err != nil {
		log.Println("Error marshalling response from gRPC server")
	}
	result := &ads.SearchResponse{}
	if err := json.Unmarshal(marshalled, &result); err != nil {
		log.Println("Error parsing the search result data from the server", err)
		restError := errors.NewServerError("invalid response from server")
		return nil, restError
	}

	return result, nil
}
