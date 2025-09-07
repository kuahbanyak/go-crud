package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
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
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	createdProduct, err := h.productUsecase.CreateProduct(r.Context(), &product)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to create product", err)
		return
	}

	response.Success(w, http.StatusCreated, "Product created successfully", createdProduct)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	product, err := h.productUsecase.GetProductByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Product not found", err)
		return
	}

	if product == nil {
		response.Error(w, http.StatusNotFound, "Product not found", nil)
		return
	}

	response.Success(w, http.StatusOK, "Product retrieved successfully", product)
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	filter := &entities.ProductFilter{}

	// Parse query parameters
	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name = name
	}

	if category := r.URL.Query().Get("category"); category != "" {
		filter.Category = category
	}

	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = minPrice
		}
	}

	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = maxPrice
		}
	}

	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filter.IsActive = &isActive
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	} else {
		filter.Limit = 10 // default limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	products, err := h.productUsecase.GetProducts(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get products", err)
		return
	}

	response.Success(w, http.StatusOK, "Products retrieved successfully", products)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	var product entities.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	product.ID = id

	updatedProduct, err := h.productUsecase.UpdateProduct(r.Context(), id, &product)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to update product", err)
		return
	}

	response.Success(w, http.StatusOK, "Product updated successfully", updatedProduct)
}

func (h *ProductHandler) UpdateProductStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	id, err := types.ParseMSSQLUUID(idStr)
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
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		response.Error(w, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	err = h.productUsecase.DeleteProduct(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete product", err)
		return
	}

	response.Success(w, http.StatusOK, "Product deleted successfully", nil)
}
