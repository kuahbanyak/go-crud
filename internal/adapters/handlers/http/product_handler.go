package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/adapters/handlers/http/helpers"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

type ProductHandler struct {
	productUsecase *usecases.ProductUsecase
}

func NewProductHandler(productUsecase *usecases.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
	}
}
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product entities.Product
	if err := helpers.DecodeJSONBody(r, &product); err != nil {
		helpers.RespondBadRequest(w, r, "Invalid JSON format", err)
		return
	}

	createdProduct, err := h.productUsecase.CreateProduct(r.Context(), &product)
	if err != nil {
		helpers.RespondBadRequest(w, r, "Failed to create product", err)
		return
	}

	helpers.RespondSuccess(w, r, http.StatusCreated, "Product created successfully", createdProduct)
}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetUUIDFromPath(r, "id")
	if err != nil {
		helpers.RespondBadRequest(w, r, "Invalid product ID", err)
		return
	}

	product, err := h.productUsecase.GetProductByID(r.Context(), id)
	if err != nil {
		helpers.RespondNotFound(w, r, "Product not found")
		return
	}

	if product == nil {
		helpers.RespondNotFound(w, r, "Product not found")
		return
	}

	helpers.RespondSuccess(w, r, http.StatusOK, "Product retrieved successfully", product)
}
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	filter := &entities.ProductFilter{
		Name:     helpers.GetQueryParam(r, "name"),
		Category: helpers.GetQueryParam(r, "category"),
		MinPrice: helpers.ParseFloatQueryParam(r, "min_price", 0),
		MaxPrice: helpers.ParseFloatQueryParam(r, "max_price", 0),
		IsActive: helpers.ParseBoolQueryParam(r, "is_active"),
	}

	filter.Limit, filter.Offset = helpers.ParsePaginationParams(r, 10, 100)

	products, err := h.productUsecase.GetProducts(r.Context(), filter)
	if err != nil {
		helpers.RespondInternalError(w, r, "Failed to get products", err)
		return
	}

	helpers.RespondSuccess(w, r, http.StatusOK, "Products retrieved successfully", products)
}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetUUIDFromPath(r, "id")
	if err != nil {
		helpers.RespondBadRequest(w, r, "Invalid product ID", err)
		return
	}

	var product entities.Product
	if err := helpers.DecodeJSONBody(r, &product); err != nil {
		helpers.RespondBadRequest(w, r, "Invalid JSON format", err)
		return
	}

	product.ID = id
	updatedProduct, err := h.productUsecase.UpdateProduct(r.Context(), id, &product)
	if err != nil {
		helpers.RespondBadRequest(w, r, "Failed to update product", err)
		return
	}

	helpers.RespondSuccess(w, r, http.StatusOK, "Product updated successfully", updatedProduct)
}
func (h *ProductHandler) UpdateProductStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "Product ID is required", nil)
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid product ID", err)
		return
	}
	var req struct {
		Stock int `json:"stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}
	err = h.productUsecase.UpdateProductStock(r.Context(), id, req.Stock)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to update product stock", err)
		return
	}
	response.Success(w, http.StatusOK, "Product stock updated successfully", map[string]interface{}{
		"id":    id,
		"stock": req.Stock,
	})
}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetUUIDFromPath(r, "id")
	if err != nil {
		helpers.RespondBadRequest(w, r, "Invalid product ID", err)
		return
	}

	err = h.productUsecase.DeleteProduct(r.Context(), id)
	if err != nil {
		helpers.RespondInternalError(w, r, "Failed to delete product", err)
		return
	}

	helpers.RespondSuccess(w, r, http.StatusOK, "Product deleted successfully", nil)
}
