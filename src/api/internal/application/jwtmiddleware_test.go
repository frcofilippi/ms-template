package application

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEnsureValidToken_HappyPath(t *testing.T) {

	os.Setenv("AUTH0_DOMAIN", "frappe-dev.us.auth0.com")
	os.Setenv("AUTH0_AUDIENCE", "pedime-api")

	token := "***REMOVED***"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	wrapped := EnsureValidToken(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", rr.Code, rr.Body.String())
	}
}
