package handlers

type CreateProductRequest struct {
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}
