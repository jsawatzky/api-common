package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/spf13/viper"

	"github.com/jsawatzky/go-common/api"
)

const userProperty = "auth0-user"

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

var jwksCache = ttlcache.NewCache()

const (
	cacheKey     = "jwks"
	cacheTimeout = 6 * time.Hour
)

func init() {
	jwksCache.SetLoaderFunction(func(key string) (data interface{}, ttl time.Duration, err error) {
		resp, err := http.Get(fmt.Sprintf("%s/.well-known/jwks.json", viper.GetString("auth0_tenent")))

		if err != nil {
			return nil, cacheTimeout, err
		}
		defer resp.Body.Close()

		var jwks = Jwks{}
		err = json.NewDecoder(resp.Body).Decode(&jwks)
		if err != nil {
			return nil, cacheTimeout, err
		}

		return jwks, cacheTimeout, nil
	})
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""

	jwks, err := jwksCache.Get(cacheKey)
	if err != nil {
		return cert, err
	}

	for k := range jwks.(Jwks).Keys {
		if token.Header["kid"] == jwks.(Jwks).Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.(Jwks).Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		return cert, errors.New("could not find appropriate key")
	}

	return cert, nil
}

func NewAuth0Middleware(ip IdentityProvider) func(http.Handler) http.Handler {
	if !(viper.IsSet("auth0_audience") && viper.IsSet("auth0_issuer") && viper.IsSet("auth0_tenent")) {
		logger.Fatal("Missing required config")
	}

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(t *jwt.Token) (interface{}, error) {
			aud := viper.GetString("auth0_audience")
			if !t.Claims.(jwt.MapClaims).VerifyAudience(aud, false) {
				return t, errors.New("invalid audience")
			}

			iss := viper.GetString("auth0_issuer")
			if !t.Claims.(jwt.MapClaims).VerifyIssuer(iss, false) {
				return t, errors.New("invalid issuer")
			}

			cert, err := getPemCert(t)
			if err != nil {
				return t, err
			}

			result, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, err
		},
		SigningMethod: jwt.SigningMethodRS256,
		UserProperty:  userProperty,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			api.EncodeResponse(w, http.StatusForbidden, AuthError(err))
		},
	})

	return func(h http.Handler) http.Handler {
		return jwtMiddleware.Handler(
			http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				token, ok := r.Context().Value(userProperty).(*jwt.Token)
				if !ok {
					logger.Error("Could not retreive JWT")
					api.EncodeResponse(rw, http.StatusForbidden, AuthError("could not retreive JWT"))
				}
				subject := token.Claims.(jwt.MapClaims)["sub"].(string)
				user, err := ip.Get(r.Context(), subject)
				if err != nil {
					if errors.Is(err, ErrPermissionDenied) {
						logger.Warn("Permission denied for user \"%s\"", subject)
					} else {
						logger.Error("Error retreiving user identity: %v", err)
					}
					api.EncodeResponse(rw, http.StatusForbidden, AuthError(err.Error()))
					return
				}
				ctx := context.WithValue(r.Context(), userKey, user)
				h.ServeHTTP(rw, r.WithContext(ctx))
			}),
		)
	}
}
