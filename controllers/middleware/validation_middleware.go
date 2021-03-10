package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/jebo87/makako-gateway/structs"
	"gitlab.com/jebo87/makako-gateway/utils/errors"
	"gitlab.com/jebo87/makako-gateway/utils/utils_http"
)

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

//ValidateMiddleware This is the middleware used to protect certain api calls
func ValidateMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		//enableCors(&w)
		if utils_http.IsPreflight(c) {
			return
		}
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				if bearerToken[0] == "Bearer" {
					log.Println(bearerToken)
					if validateToken(bearerToken[1]) {
						//everything is ok, proceed with allow the exectution of the next function
						c.Next()
					} else {
						respErr := errors.NewUnauthorizedError("Missing / expired auth token")
						c.AbortWithStatusJSON(respErr.Status, respErr)
					}
				} else {
					respErr := errors.NewUnauthorizedError("Invalid authorization token")
					c.AbortWithStatusJSON(respErr.Status, respErr)
				}
			}
		} else {
			respErr := errors.NewUnauthorizedError("Authorization required")
			c.AbortWithStatusJSON(respErr.Status, respErr)
		}
	}
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
