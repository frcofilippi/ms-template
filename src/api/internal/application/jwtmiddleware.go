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

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
	Sub   string `json:"sub"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	if c.Sub == "" {
		return fmt.Errorf("not able to parse userId from request")
	}
	return nil
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
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

	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		logger.GetLogger().Info("JWT middleware invoked")
	// 		// Run the JWT middleware
	// 		m.ServeHTTP(w, r)

	// 		// Extract the validated claims from the context using jwtmiddleware.ContextKey{}
	// 		token, ok := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	// 		if !ok || token == nil {
	// 			return // JWT not valid, error handler already ran
	// 		}
	// 		customClaims, ok := token.CustomClaims.(*CustomClaims)
	// 		logger.GetLogger().Info("Custom claims parsed")
	// 		if !ok || customClaims == nil {
	// 			return // No custom claims, error handler already ran
	// 		}

	// 		next.ServeHTTP(w, r)
	// 		logger.GetLogger().Info("After m.ServeHTTP")
	// 	})
}

// HasScope checks whether our claims have a specific scope.
func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}
