package app

import (
	"gitlab.com/jebo87/makako-gateway/controllers/listings"
)

// 	// router.HandleFunc("/ads", httputils.ValidateMiddleware(adsHandler)).Methods("GET", "OPTIONS")
// 	// router.HandleFunc("/ads", (gwhandlers.AdsHandler)).Methods("GET", "OPTIONS")
// 	router.HandleFunc("/ads", (gwhandlers.ListingHandler)).Methods("GET", "POST", "OPTIONS")
// 	router.HandleFunc("/ads/{key}", gwhandlers.AdHandler).Methods("GET", "OPTIONS")
// 	router.HandleFunc("/listing/new", httputils.ValidateMiddleware(gwhandlers.NewListingHandler)).Methods("POST", "OPTIONS")
// 	router.HandleFunc("/ad_count", httputils.ValidateMiddleware(gwhandlers.AdsCountHandler)).Methods("GET", "OPTIONS")
// 	router.HandleFunc("/{userId}/listings", httputils.ValidateMiddleware(gwhandlers.UserListingsHandler)).Methods("POST", "OPTIONS")

func MapURLs() {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/listings" /*middleware.ValidateMiddleware(),*/, listings.GetListings)
		v1.GET("/listings/:id" /*middleware.ValidateMiddleware(),*/, listings.GetSingleListing)
	}
}
