package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/shared/utils"
)

type ProductUsecase struct {
	productRepo repositories.ProductRepository
	validator   *utils.Validator
}

func NewProductUsecase(productRepo repositories.ProductRepository, validator *utils.Validator) *ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
		validator:   validator,
	}
}

func (uc *ProductUsecase) CreateProduct(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	if err := uc.validateProduct(product); err != nil {
		return nil, err
	}

	if product.SKU == "" {
		product.SKU = uc.generateSKU(product.Name, product.Category)
	}

	existing, err := uc.productRepo.GetBySKU(ctx, product.SKU)
	if err == nil && existing != nil {
		return nil, errors.New("product with this SKU already exists")
	}

	return uc.productRepo.Create(ctx, product)
}

func (uc *ProductUsecase) GetProductByID(ctx context.Context, id types.MSSQLUUID) (*entities.Product, error) {
	if id.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, errors.New("invalid product ID")
	}

	return uc.productRepo.GetByID(ctx, id)
}

func (uc *ProductUsecase) GetProducts(ctx context.Context, filter *entities.ProductFilter) ([]*entities.Product, error) {
	if filter == nil {
		filter = &entities.ProductFilter{Limit: 10, Offset: 0}
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	if filter.Limit > 100 {
		filter.Limit = 100
	}

	return uc.productRepo.GetAll(ctx, filter)
}

func (uc *ProductUsecase) UpdateProduct(ctx context.Context, id types.MSSQLUUID, product *entities.Product) (*entities.Product, error) {
	if id.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, errors.New("invalid product ID")
	}

	existing, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("product not found")
	}

	if err := uc.validateProduct(product); err != nil {
		return nil, err
	}

	return uc.productRepo.Update(ctx, id, product)
}

func (uc *ProductUsecase) DeleteProduct(ctx context.Context, id types.MSSQLUUID) error {
	if id.String() == "00000000-0000-0000-0000-000000000000" {
		return errors.New("invalid product ID")
	}

	existing, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("product not found")
	}

	return uc.productRepo.Delete(ctx, id)
}

func (uc *ProductUsecase) UpdateProductStock(ctx context.Context, id types.MSSQLUUID, stock int) error {
	if id.String() == "00000000-0000-0000-0000-000000000000" {
		return errors.New("invalid product ID")
	}

	if stock < 0 {
		return errors.New("stock cannot be negative")
	}

	existing, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("product not found")
	}

	return uc.productRepo.UpdateStock(ctx, id, stock)
}

func (uc *ProductUsecase) validateProduct(product *entities.Product) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}

	if !uc.validator.IsNotEmpty(product.Name) {
		return errors.New("product name is required")
	}

	if !uc.validator.IsValidLength(product.Name, 1, 255) {
		return errors.New("product name must be between 1 and 255 characters")
	}

	if product.Price < 0 {
		return errors.New("product price cannot be negative")
	}

	if product.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	return nil
}

func (uc *ProductUsecase) generateSKU(name, category string) string {
	categoryCode := strings.ToUpper(category)
	if len(categoryCode) > 3 {
		categoryCode = categoryCode[:3]
	}

	nameCode := strings.ToUpper(strings.ReplaceAll(name, " ", ""))
	if len(nameCode) > 3 {
		nameCode = nameCode[:3]
	}

	return fmt.Sprintf("%s%s%d", categoryCode, nameCode, len(name)*len(category))
}
