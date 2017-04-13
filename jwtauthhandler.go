package main

import (
	"errors"
	"net/http"

	"time"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

type JWTAuthHandler struct {
	Keys       map[string]string
	Next       http.Handler
	Now        func() time.Time
	middleware *jwtmiddleware.JWTMiddleware
}

func NewJWTAuthHandler(keys map[string]string, now func() time.Time, next http.Handler) JWTAuthHandler {
	h := JWTAuthHandler{
		Keys: keys,
		Next: next,
		Now:  time.Now,
	}
	h.middleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Assume standard claims of "iss", "exp" and "iat".
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, errors.New("JWT claims not found")
			}

			if !claims.VerifyExpiresAt(now().Unix(), true) {
				return nil, errors.New("token expired")
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
	return h
}

func (jwth JWTAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwth.middleware.Handler(jwth.Next).ServeHTTP(w, r)
}
