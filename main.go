package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"io/ioutil"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

var remoteURLFlag = flag.String("remoteURL", "http://localhost:8080", "The remote host to proxy to.")
var keysFlag = flag.String("keys", "config.json", "The location of the JSON map containing issuers and their public keys.")
var portFlag = flag.Int("port", 9090, "The port for the proxy to listen on.")

func main() {
	flag.Parse()

	u, err := url.Parse(*remoteURLFlag)
	if err != nil {
		fmt.Printf("Failed to parse remoteURL %s with error %v", *remoteURLFlag, err)
		os.Exit(-1)
	}

	file, err := os.Open(*keysFlag)
	if err != nil {
		fmt.Printf("Failed to open file %s with error %v", *keysFlag, err)
		os.Exit(-1)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Failed to read file %s with error %v", *keysFlag, err)
		os.Exit(-1)
	}

	keys := make(map[string]string)
	err = json.Unmarshal(data, &keys)
	if err != nil {
		fmt.Printf("Failed to parse JSON file %s with error %v", *keysFlag, err)
		os.Exit(-1)
	}

	jwtmw := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Assume standard claims of "iss", "exp" and "iat".
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, errors.New("JWT claims not found")
			}

			// Find the public key to match the issuer.
			issuerClaim, ok := claims["iss"]
			if !ok {
				return nil, errors.New("iss not found")
			}

			issuer, ok := issuerClaim.(string)
			if !ok {
				return nil, errors.New("iss was not in correct format")
			}

			pub, ok := keys[issuer]
			if !ok {
				return nil, errors.New("iss not valid")
			}

			return jwt.ParseRSAPublicKeyFromPEM([]byte(pub))
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodRS256,
	})

	// Wrap the proxy in authentication.
	app := jwtmw.Handler(httputil.NewSingleHostReverseProxy(u))

	http.ListenAndServe(":"+strconv.Itoa(*portFlag), app)
}
