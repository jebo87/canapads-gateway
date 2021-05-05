package app

import (
	"gitlab.com/jebo87/makako-gateway/controllers/listings"
	"gitlab.com/jebo87/makako-gateway/controllers/middleware"
)

// 	// router.HandleFunc("/ads", httputils.ValidateMiddleware(adsHandler)).Methods("GET", "OPTIONS")
// 	// router.HandleFunc("/ads", (gwhandlers.AdsHandler)).Methods("GET", "OPTIONS")
// 	router.HandleFunc("/ads", (gwhandlers.ListingHandler)).Methods("GET", "POST", "OPTIONS")
// 	router.HandleFunc("/ads/{key}", gwhandlers.AdHandler).Methods("GET", "OPTIONS")
// 	router.HandleFunc("/listing/new", httputils.ValidateMiddleware(gwhandlers.NewListingHandler)).Methods("POST", "OPTIONS")
// 	router.HandleFunc("/ad_count", httputils.ValidateMiddleware(gwhandlers.AdsCountHandler)).Methods("GET", "OPTIONS")
// 	router.HandleFunc("/{userId}/listings", httputils.ValidateMiddleware(gwhandlers.UserListingsHandler)).Methods("POST", "OPTIONS")

func MapURLs() {
	router.POST("api/v1/listings", middleware.ValidateMiddleware(), listings.GetListings)
	router.GET("api/v1/listings/:id", listings.GetSingleListing)
}
