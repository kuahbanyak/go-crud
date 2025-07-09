package handler

import (
	"go-crud/internal/domain/entity"
	"go-crud/internal/usecase"
	"go-crud/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

// NewProductHandler creates a new product handler
func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided information
// @Tags products
// @Accept json
// @Produce json
// @Param product body entity.CreateProductRequest true "Product creation data"
// @Success 201 {object} response.Response{data=entity.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req entity.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	product, err := h.productUsecase.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create product", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Product created successfully", product)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response{data=entity.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	product, err := h.productUsecase.GetProductByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			response.Error(c, http.StatusNotFound, "Product not found", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get product", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product retrieved successfully", product)
}

// GetProducts godoc
// @Summary Get all products
// @Description Get all products with pagination
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.PaginatedResponse{data=[]entity.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	products, total, err := h.productUsecase.GetProducts(c.Request.Context(), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get products", err.Error())
		return
	}

	response.Paginated(c, http.StatusOK, "Products retrieved successfully", products, total, limit, offset)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body entity.UpdateProductRequest true "Product update data"
// @Success 200 {object} response.Response{data=entity.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	var req entity.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	product, err := h.productUsecase.UpdateProduct(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err.Error() == "product not found" {
			response.Error(c, http.StatusNotFound, "Product not found", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to update product", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product updated successfully", product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	err = h.productUsecase.DeleteProduct(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			response.Error(c, http.StatusNotFound, "Product not found", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete product", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product deleted successfully", nil)
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Description Get products filtered by category with pagination
// @Tags products
// @Accept json
// @Produce json
// @Param category path string true "Product Category"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response{data=[]entity.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/category/{category} [get]
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	products, err := h.productUsecase.GetProductsByCategory(c.Request.Context(), category, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get products by category", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", products)
}
