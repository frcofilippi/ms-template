package product

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestProductRouter_Routes_Registered(t *testing.T) {
	router := NewProductRouter(nil)
	handler := router.Routes()

	getReq := httptest.NewRequest("GET", "/123", nil)
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)
	if getRec.Code == http.StatusNotFound {
		t.Errorf("GET /{product-id} route not registered got %d", getRec.Code)
	}

	postReq := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{"name":"foo","cost":1}`)))
	postRec := httptest.NewRecorder()
	handler.ServeHTTP(postRec, postReq)
	if postRec.Code == http.StatusNotFound {
		t.Errorf("GET /{product-id} route not registered got %d", postRec.Code)
	}
}

func TestProductRouter_UnregisteredRoute_Returns404(t *testing.T) {
	const expectedStatusCode = http.StatusBadRequest
	router := NewProductRouter(nil)
	handler := router.Routes()

	mux := chi.NewMux()
	mux.Mount("/api/v1/product", handler)

	server := httptest.NewServer(mux)

	url := fmt.Sprintf("%s/api/v1/product/notfound", server.URL)
	response, err := http.Get(url)

	if err != nil {
		t.Errorf("expected to be able to create the request but got error %s", err.Error())
	}

	if response.StatusCode != expectedStatusCode {
		t.Errorf("expected %d for unregistered route, got %d", expectedStatusCode, response.StatusCode)
	}
}

func TestProductRouterGetProductById_InvalidIdReturns400(t *testing.T) {
	router := NewProductRouter(nil)
	const expectedStatusCode = http.StatusBadRequest
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/dslakjdlkasj", nil)

	router.Routes().ServeHTTP(rr, req)

	if rr.Result().StatusCode != expectedStatusCode {
		t.Errorf("expected %d status code and got %d", expectedStatusCode, rr.Result().StatusCode)
	}

}
