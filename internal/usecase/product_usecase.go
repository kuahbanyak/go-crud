package usecase

import (
	"context"
	"errors"
	"go-crud/internal/domain/entity"
	"go-crud/internal/domain/repository"
	"gorm.io/gorm"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.ProductResponse, error)
	GetProductByID(ctx context.Context, id uint) (*entity.ProductResponse, error)
	GetProducts(ctx context.Context, limit, offset int) ([]*entity.ProductResponse, int64, error)
	UpdateProduct(ctx context.Context, id uint, req *entity.UpdateProductRequest) (*entity.ProductResponse, error)
	DeleteProduct(ctx context.Context, id uint) error
	GetProductsByCategory(ctx context.Context, category string, limit, offset int) ([]*entity.ProductResponse, error)
}

type productUsecase struct {
	productRepo repository.ProductRepository
}

func NewProductUsecase(productRepo repository.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (u *productUsecase) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.ProductResponse, error) {

	existingProduct, err := u.productRepo.GetBySKU(ctx, req.SKU)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingProduct != nil {
		return nil, errors.New("product with this SKU already exists")
	}

	product := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Category:    req.Category,
		SKU:         req.SKU,
		IsActive:    true,
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return u.mapToResponse(product), nil
}

func (u *productUsecase) GetProductByID(ctx context.Context, id uint) (*entity.ProductResponse, error) {
	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return u.mapToResponse(product), nil
}

func (u *productUsecase) GetProducts(ctx context.Context, limit, offset int) ([]*entity.ProductResponse, int64, error) {
	products, err := u.productRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.productRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*entity.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = u.mapToResponse(product)
	}

	return responses, total, nil
}

func (u *productUsecase) UpdateProduct(ctx context.Context, id uint, req *entity.UpdateProductRequest) (*entity.ProductResponse, error) {

	_, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	updateProduct := &entity.Product{}
	if req.Name != nil {
		updateProduct.Name = *req.Name
	}
	if req.Description != nil {
		updateProduct.Description = *req.Description
	}
	if req.Price != nil {
		updateProduct.Price = *req.Price
	}
	if req.Quantity != nil {
		updateProduct.Quantity = *req.Quantity
	}
	if req.Category != nil {
		updateProduct.Category = *req.Category
	}
	if req.IsActive != nil {
		updateProduct.IsActive = *req.IsActive
	}

	if err := u.productRepo.Update(ctx, id, updateProduct); err != nil {
		return nil, err
	}

	updatedProduct, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u.mapToResponse(updatedProduct), nil
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id uint) error {

	_, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	return u.productRepo.Delete(ctx, id)
}

func (u *productUsecase) GetProductsByCategory(ctx context.Context, category string, limit, offset int) ([]*entity.ProductResponse, error) {
	products, err := u.productRepo.GetByCategory(ctx, category, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*entity.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = u.mapToResponse(product)
	}

	return responses, nil
}

func (u *productUsecase) mapToResponse(product *entity.Product) *entity.ProductResponse {
	return &entity.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		Category:    product.Category,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
