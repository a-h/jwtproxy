package main

import (
	"context"
	"errors"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	golangjwt "github.com/golang-jwt/jwt/v4"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"strings"
	"time"
)

// JWTAuthHandler provides the capability to authenticate incoming HTTP requests.
type JWTAuthHandler struct {
	Keys       map[string]string
	Next       http.Handler
	Now        func() time.Time
	middleware *jwtmiddleware.JWTMiddleware
}

const signatureAlgorithm = validator.RS256

var (
	invalidSingingAlg         = fmt.Sprintf("invalid signing algorithm, expected %q", signatureAlgorithm)
	ErrJWTInvalidCharacter    = errors.New("invalid character")
	ErrJWTIssuerNotValid      = errors.New("iss not valid")
	ErrJWTInvalidExpiry       = errors.New("token expired")
	ErrJWTInvalidSigningAlg   = errors.New(invalidSingingAlg)
	ErrJWTClaimFormatNotValid = errors.New("could not find claims")
	ErrJWTInvalidToken        = errors.New("Authorization header format must be Bearer {token}")
)

// More finegrained ErrorHandler than the default.
func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusUnauthorized)

	switch {
	case errors.Is(err, ErrJWTInvalidCharacter):
		_, _ = w.Write([]byte("invalid character"))
	case errors.Is(err, ErrJWTInvalidSigningAlg):
		_, _ = w.Write([]byte(invalidSingingAlg))
	case errors.Is(err, ErrJWTInvalidExpiry):
		_, _ = w.Write([]byte("token expired"))
	case errors.Is(err, ErrJWTIssuerNotValid):
		_, _ = w.Write([]byte("iss not valid"))
	case errors.Is(err, jwtmiddleware.ErrJWTMissing):
		_, _ = w.Write([]byte("Required authorization token not found"))
	case errors.Is(err, ErrJWTClaimFormatNotValid):
		_, _ = w.Write([]byte("invalid claims"))
	case errors.Is(err, ErrJWTInvalidToken):
		_, _ = w.Write([]byte("Authorization header format must be Bearer {token}"))
	default:
		_, _ = w.Write([]byte("jwt invalid"))
	}
}

// NewJWTAuthHandler creates a new JWTAuthHandler, passing in a map of issuers to public RSA keys, and a
// time provider to allow for variation of the time.
func NewJWTAuthHandler(keys map[string]string, now func() time.Time, authHeader string, next http.Handler) JWTAuthHandler {
	h := JWTAuthHandler{
		Keys: keys,
		Next: next,
		Now:  time.Now,
	}

	h.middleware = jwtmiddleware.New(
		func(ctx context.Context, tokenString string) (interface{}, error) {
			token, err := jwt.ParseSigned(tokenString)

			if err != nil {
				return nil, ErrJWTInvalidCharacter
			}

			if string(signatureAlgorithm) != token.Headers[0].Algorithm {
				return nil, ErrJWTInvalidSigningAlg
			}

			claims := &jwt.Claims{}
			err = token.UnsafeClaimsWithoutVerification(claims)
			if err != nil {
				return nil, ErrJWTClaimFormatNotValid
			}
			if claims.Expiry == nil {
				return nil, ErrJWTInvalidExpiry
			}
			pub, ok := keys[claims.Issuer]
			if !ok {
				return nil, ErrJWTIssuerNotValid
			}

			key, err := golangjwt.ParseRSAPublicKeyFromPEM([]byte(pub))
			if err != nil {
				return nil, ErrJWTInvalidSigningAlg
			}

			expectedClaims := jwt.Expected{
				Time: now(),
			}

			claimDest := []interface{}{&jwt.Claims{}}
			if err := token.Claims(key, claimDest...); err != nil {
				return nil, ErrJWTInvalidExpiry
			}

			registeredClaims := *claimDest[0].(*jwt.Claims)
			if err := registeredClaims.Validate(expectedClaims); err != nil {
				return nil, ErrJWTInvalidExpiry
			}

			return claims, nil

		},
		jwtmiddleware.WithErrorHandler(ErrorHandler),
		jwtmiddleware.WithTokenExtractor(func(r *http.Request) (string, error) {
			authHeader := r.Header.Get(authHeader)
			if authHeader == "" {
				return "", nil // No error, just no JWT.
			}

			authHeaderParts := strings.Fields(authHeader)
			if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
				return "", ErrJWTInvalidToken
			}

			return authHeaderParts[1], nil
		}))

	return h
}

func (jwth JWTAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwth.middleware.CheckJWT(jwth.Next).ServeHTTP(w, r)
}
