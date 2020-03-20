package httputils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.com/jebo87/makako-gateway/structs"
)

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

//ValidateMiddleware This is the middleware used to protect certain api calls
func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//enableCors(&w)
		if req.Method == "OPTIONS" {
			log.Println("Options request")

			w.Header().Add("Access-Control-Allow-Methods", "GET")
			w.Header().Add("Access-Control-Allow-Headers", "Authorization")
			w.Header().Add("Access-Control-Allow-Origin", "https://www.canapads.ca")
			w.WriteHeader(http.StatusOK)

			return
		}
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				if bearerToken[0] == "Bearer" {
					log.Println(bearerToken)
					if validateToken(bearerToken[1]) {
						//everything is ok, proceed with allow the exectution of the next function
						next(w, req)
					} else {
						// json.NewEncoder(w).Encode(structs.Exception{Message: "Token invalid or expired"})
						http.Error(w, `{Message: "Token invalid or expired"}`, http.StatusUnauthorized)

					}
				} else {
					http.Error(w, `{Message: "Invalid authorization token"}`, http.StatusUnauthorized)

				}
			}
		} else {
			http.Error(w, `{Message: "Authorization required"}`, http.StatusUnauthorized)
		}
	})
}
func validateToken(token string) bool {
	form := url.Values{}
	form.Add("token", token)
	log.Println("validating token " + token)
	req, erro := http.NewRequest("POST", "http://localhost:4445/oauth2/introspect", strings.NewReader(form.Encode()))
	if erro != nil {
		log.Println(erro.Error())

		return false

	}
	log.Println("adding headers")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Add("Authorization", "Basic "+basicAuth("api1", "TopSecret2"))
	resp, err := netClient.Do(req)

	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	var validationObject structs.ValidationResponse
	json.Unmarshal(bodyBytes, &validationObject)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Validation for the request has returned %v \n", validationObject.Active)

	return validationObject.Active

}
