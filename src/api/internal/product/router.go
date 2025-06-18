package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"frcofilippi/pedimeapp/internal/application"
	"frcofilippi/pedimeapp/shared/utils"
	"net/http"
	"strconv"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
)

type ProductRouter struct {
	productService ProductService
}

func (ph *ProductRouter) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	prodId := chi.URLParam(r, "product-id")

	if prodId == "" {
		utils.SendErrorResponse(w, r, http.StatusBadRequest, "productId parameter is empty", nil)
		return
	}

	id, err := strconv.ParseInt(prodId, 0, 64)

	if err != nil {
		utils.SendErrorResponse(w, r, http.StatusBadRequest, "invalid productId", err)
		return
	}

	userId, err := GetUserIdFromRequest(r)
	if err != nil {
		utils.SendErrorResponse(w, r, http.StatusUnauthorized, "not authorized to see product information", err)
		return
	}

	command := &GetProductByIdCommand{
		ProductId: id,
		UserId:    userId,
	}

	result, err := ph.productService.GetProductById(*command)

	if err != nil {
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, "error processing your request", err)
		return
	}

	if result == nil {
		utils.SendApiResponse(w, r, http.StatusNotFound, "product not found", nil)
		return
	}

	utils.SendApiResponse(w, r, http.StatusOK, "product found", result)

}

func (ph *ProductRouter) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {

	userID, err := GetUserIdFromRequest(r)
	if err != nil {
		utils.SendErrorResponse(w, r, http.StatusUnauthorized, "not authorized to perform the operation", errors.New("Missing UserId."))
		return
	}

	var request CreateProductRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.SendErrorResponse(w, r, http.StatusBadRequest, "not able to parse the request", err)
		return
	}

	command := &CreateNewProductCommand{
		Name:   request.Name,
		Cost:   request.Cost,
		UserId: userID,
	}

	id, err := ph.productService.CreateNewProduct(*command)
	if err != nil {
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SendApiResponse(w, r, http.StatusCreated, "prduct created", id)
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
