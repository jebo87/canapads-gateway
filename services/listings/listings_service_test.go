package listings

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/jebo87/makako-gateway/clients"
	"gitlab.com/jebo87/makako-gateway/utils/mocks"
	"gitlab.com/jebo87/makako-grpc/ads"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	clients.GrpcClient = &mocks.MockGrpcClient{}
	os.Exit(m.Run())
}

func TestGetSingleListingInvalidResponse(t *testing.T) {
	mocks.AdDetailFunc = func(ctx context.Context, in *ads.Text, opts ...grpc.CallOption) (*ads.Ad, error) {
		return nil, errors.New(`invalid response from grpc server rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial tcp 127.0.0.1:7777: connect: connection refused"`)
	}
	result, err := ListingsService.GetSingleListing(context.Background(), "41")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, 500, err.Status)
	assert.EqualValues(t, "internal_server_error", err.Error)
	assert.EqualValues(t, `invalid response from grpc server rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial tcp 127.0.0.1:7777: connect: connection refused"`, err.Message)
}

func TestGetSingleListingOK(t *testing.T) {
	mocks.AdDetailFunc = func(ctx context.Context, in *ads.Text, opts ...grpc.CallOption) (*ads.Ad, error) {
		listing := &ads.Ad{
			Id:          41,
			Title:       "Test title",
			Description: "Test Description",
			City:        "Montreal",
		}
		return listing, nil
	}
	result, err := ListingsService.GetSingleListing(context.Background(), "41")
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
