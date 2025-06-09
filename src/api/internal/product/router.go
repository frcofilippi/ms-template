package product

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ProductRouter struct {
	productService ProductService
}

func (ph *ProductRouter) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	requestId, _ := r.Context().Value(middleware.RequestIDKey).(string)
	prodId := chi.URLParam(r, "product-id")
	if prodId == "" {
		log.Default().Printf("[%s]ProdId not received. Url: %s", requestId, r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(prodId, 0, 64)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO: Parse the customerid from the jsonweb token

	command := &GetProductByIdCommand{
		ProductId:  id,
		CustomerId: 2,
	}

	result, err := ph.productService.GetProductById(*command)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	jresult, _ := json.MarshalIndent(result, "", "  ")

	w.Write([]byte(jresult))

	w.WriteHeader(http.StatusOK)
	log.Default().Printf("[%s]Successfully processed. Url: %s", requestId, r.URL.String())
}

func (ph *ProductRouter) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	requestId, _ := r.Context().Value(middleware.RequestIDKey).(string)
	var request CreateProductRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not able to parse your request"))
		return
	}

	log.Default().Printf("[%s] With values: %v \n", requestId, request)

	command := &CreateNewProductCommand{
		Name:       request.Name,
		Cost:       request.Cost,
		CustomerId: 2,
	}

	id, err := ph.productService.CreateNewProduct(*command)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	result := strconv.FormatInt(id, 10)
	w.Write([]byte(result))
	log.Default().Printf("[%s]Successfully processed. Url: %s", requestId, r.URL.String())
}

func (ph *ProductRouter) Routes() http.Handler {
	router := chi.NewRouter()
	router.Get("/{product-id}", ph.HandleGetProduct)
	router.Post("/", ph.HandleCreateProduct)
	return router
}

func NewProductRouter(service ProductService) *ProductRouter {
	return &ProductRouter{
		productService: service,
	}
}
