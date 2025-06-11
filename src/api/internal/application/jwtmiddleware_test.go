package application

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEnsureValidToken_HappyPath(t *testing.T) {
	// Set required env vars for the middleware
	os.Setenv("AUTH0_DOMAIN", "frappe-dev.us.auth0.com")
	os.Setenv("AUTH0_AUDIENCE", "pedime-api")

	// The provided valid JWT token
	// userId auth0|68486eb6bad0593787a8b06a
	token := "***REMOVED***"

	// Dummy handler to check if request passes through middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Wrap with middleware
	wrapped := EnsureValidToken(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", rr.Code, rr.Body.String())
	}
}
