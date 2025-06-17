package application

import (
	"context"
	"fmt"
	"frcofilippi/pedimeapp/shared/logger"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"go.uber.org/zap"
)

type CustomClaims struct {
	Scope string `json:"scope"`
	Sub   string `json:"sub"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	if c.Sub == "" {
		return fmt.Errorf("not able to parse userId from request")
	}
	return nil
}

func EnsureValidToken(next http.Handler) http.Handler {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	logger.GetLogger().Info(
		"setting up jwt validation",
		zap.String("issuerUrl", issuerURL.String()),
	)
	if err != nil {
		logger.GetLogger().Error("error constructing the issues URL", zap.Error(err))
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		logger.GetLogger().Error("error setting up jwt validator", zap.Error(err))
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		logger.GetLogger().Error("error parsing the token", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return middleware.CheckJWT(next)
}

func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}
