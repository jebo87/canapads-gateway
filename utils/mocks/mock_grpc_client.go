package mocks

import (
	"context"

	"gitlab.com/jebo87/makako-grpc/ads"
	"google.golang.org/grpc"
)

type MockGrpcClient struct{}

var (
	ListFunc         func(ctx context.Context, in *ads.Filter, opts ...grpc.CallOption) (*ads.SearchResponse, error)
	AdDetailFunc     func(ctx context.Context, in *ads.Text, opts ...grpc.CallOption) (*ads.Ad, error)
	CountFunc        func(ctx context.Context, in *ads.Void, opts ...grpc.CallOption) (*ads.AdCount, error)
	AddListingFunc   func(ctx context.Context, in *ads.Ad, opts ...grpc.CallOption) (*ads.ListingID, error)
	UserListingsFunc func(ctx context.Context, in *ads.UserID, opts ...grpc.CallOption) (*ads.AdList, error)
)

func (m *MockGrpcClient) List(ctx context.Context, in *ads.Filter, opts ...grpc.CallOption) (*ads.SearchResponse, error) {
	return ListFunc(ctx, in, opts...)
}

func (m *MockGrpcClient) AdDetail(ctx context.Context, in *ads.Text, opts ...grpc.CallOption) (*ads.Ad, error) {
	return AdDetailFunc(ctx, in, opts...)
}
func (m *MockGrpcClient) Count(ctx context.Context, in *ads.Void, opts ...grpc.CallOption) (*ads.AdCount, error) {
	return CountFunc(ctx, in, opts...)
}

func (m *MockGrpcClient) AddListing(ctx context.Context, in *ads.Ad, opts ...grpc.CallOption) (*ads.ListingID, error) {
	return AddListingFunc(ctx, in, opts...)
}
func (m *MockGrpcClient) UserListings(ctx context.Context, in *ads.UserID, opts ...grpc.CallOption) (*ads.AdList, error) {
	return UserListingsFunc(ctx, in, opts...)
}
