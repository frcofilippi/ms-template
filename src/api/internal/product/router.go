package product

import (
	"encoding/json"
	"fmt"
	"frcofilippi/pedimeapp/internal/application"
	"frcofilippi/pedimeapp/shared/logger"
	"log"
	"net/http"
	"strconv"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
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

	userId, err := GetUserIdFromRequest(r)
	if err != nil {
		logger.GetLogger().Error("error parsing user from context", zap.String("handler", "handleGetProduct"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	command := &GetProductByIdCommand{
		ProductId: id,
		UserId:    userId,
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

	userID, err := GetUserIdFromRequest(r)
	if err != nil {
		logger.GetLogger().Error("error parsing user from context", zap.String("handler", "handleCreateProduct"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var request CreateProductRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not able to parse your request"))
		return
	}

	log.Default().Printf("[%s] With values: %v \n", requestId, request)

	command := &CreateNewProductCommand{
		Name:   request.Name,
		Cost:   request.Cost,
		UserId: userID,
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

func GetUserIdFromRequest(r *http.Request) (string, error) {
	token, ok := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok || token == nil {
		return "", fmt.Errorf("invalid token. Not able to get user")
	}
	customClaims, ok := token.CustomClaims.(*application.CustomClaims)
	if !ok || customClaims == nil {
		return "", fmt.Errorf("not able to parse custom claims")
	}
	userID := customClaims.Sub
	return userID, nil
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
