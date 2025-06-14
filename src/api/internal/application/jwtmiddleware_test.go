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
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkViZmFfb1FYMmw0Q1VRNVlSZ3hndSJ9.eyJpc3MiOiJodHRwczovL2ZyYXBwZS1kZXYudXMuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfDY4NDg2ZWI2YmFkMDU5Mzc4N2E4YjA2YSIsImF1ZCI6WyJwZWRpbWUtYXBpIiwiaHR0cHM6Ly9mcmFwcGUtZGV2LnVzLmF1dGgwLmNvbS91c2VyaW5mbyJdLCJpYXQiOjE3NDk1Nzc2MTAsImV4cCI6MTc0OTY2NDAxMCwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsImF6cCI6InZhb28zUkpncncxanFNTXFqS2I2dXoyaUphTUtQSU02In0.LdfqLGg3ka0sNM5AyCjv0RW5bG6GuuT1K_SIpPvNdpCyELClttLFmjx0kZfYkZRjT2id93m-oV0aVoMXbxhM_-3K2WwRb34ZrPRZPXa_rOIWGZvcgmBlJRERSR3qmb5rrarTdU13RSfxRQUF0qPVDCXxGDwdFpgnDgEMbO35qoBV-gkwdpxgBVAYMlTHhXGWZj7b9goAxuNG3YenLmKNhyXNE9ahRw2VToA9A2eB5kFBZtWc76sNDBr_Ey3GCMUzb7aOw27qHdOFRdU6y4we4Z6etH9soFvQdFD4LW2083sGbJlKRNTvHz3MFIDzLyNFD_blwU4qdOfogKcZnu5RCg"

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
