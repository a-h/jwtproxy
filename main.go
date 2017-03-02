package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

var remoteURLFlag = flag.String("remoteURL", "", "The remote host to proxy to.")
var keysFlag = flag.String("keys", "", "The location of the JSON map containing issuers and their public keys.")
var portFlag = flag.String("port", "", "The port for the proxy to listen on.")
var healthCheckFlag = flag.String("health", "/health", "The path to the healthcheck endpoint.")

func main() {
	flag.Parse()

	remoteURL, err := getRemoteURL()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	keys, err := getKeys()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	port, err := getPort()
	if err != nil {
		fmt.Println(err)
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
	auth := jwtmw.Handler(httputil.NewSingleHostReverseProxy(remoteURL))

	// Wrap the authentication in a health check (health checks don't need authentication).
	health := HealthCheckHandler{
		Path: getHealthCheckURI(),
		Next: auth,
	}

	// Wrap the health check in a logger.
	app := NewLoggingHandler(health)

	http.ListenAndServe(":"+port, app)
}

func getPort() (string, error) {
	port := os.Getenv("JWTPROXY_LISTEN_PORT")
	if port == "" {
		port = *portFlag
	}
	if port == "" {
		return "9090", errors.New("JWTPROXY_LISTEN_PORT environment variable or port command line flag not found")
	}
	return port, nil
}

func getRemoteURL() (*url.URL, error) {
	remoteURL := os.Getenv("JWTPROXY_REMOTE_URL")
	if remoteURL == "" {
		remoteURL = *remoteURLFlag
	}
	if remoteURL == "" {
		return nil, errors.New("JWTPROXY_REMOTE_URL environment variable or remoteURL command line flag not found")
	}
	u, err := url.Parse(remoteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remoteURL %s with error %v", remoteURL, err)
	}
	return u, nil
}

func getKeys() (map[string]string, error) {
	keys := make(map[string]string)
	configPath := os.Getenv("JWTPROXY_CONFIG")
	if configPath == "" {
		configPath = *keysFlag
	}
	if configPath == "" {
		return keys, errors.New("JWTPROXY_CONFIG environment variable or config command line flag not found")
	}
	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil {
		return keys, fmt.Errorf("Failed to open file %s with error %v", configPath, err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return keys, fmt.Errorf("Failed to read file %s with error %v", configPath, err)
	}
	err = json.Unmarshal(data, &keys)
	if err != nil {
		return keys, fmt.Errorf("Failed to parse JSON file %s with error %v", configPath, err)
	}
	return keys, nil
}

func getHealthCheckURI() string {
	hc := os.Getenv("JWTPROXY_HEALTHCHECK_URI")
	if hc == "" {
		return *healthCheckFlag
	}
	return hc
}
