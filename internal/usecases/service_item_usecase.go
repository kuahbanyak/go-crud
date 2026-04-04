package usecases

import (
	"context"
	"errors"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/google/uuid"
)

type ServiceItemUsecase struct {
	serviceItemRepo repositories.ServiceItemRepository
}

func NewServiceItemUsecase(serviceItemRepo repositories.ServiceItemRepository) *ServiceItemUsecase {
	return &ServiceItemUsecase{
		serviceItemRepo: serviceItemRepo,
	}
}

func (u *ServiceItemUsecase) CreateServiceItem(ctx context.Context, req *dto.CreateServiceItemRequest) (*entities.ServiceItem, error) {
	serviceItem := &entities.ServiceItem{
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		EstimatedTime:   req.EstimatedTime,
		EstimatedCost:   req.EstimatedCost,
		DisplayOrder:    req.DisplayOrder,
		RequiresBooking: req.RequiresBooking,
		IsActive:        true,
	}

	if err := u.serviceItemRepo.Create(ctx, serviceItem); err != nil {
		return nil, err
	}

	return serviceItem, nil
}

func (u *ServiceItemUsecase) GetServiceItem(ctx context.Context, id uuid.UUID) (*entities.ServiceItem, error) {
	return u.serviceItemRepo.GetByID(ctx, id)
}

func (u *ServiceItemUsecase) GetAllServiceItems(ctx context.Context) ([]*entities.ServiceItem, error) {
	return u.serviceItemRepo.GetAll(ctx)
}

func (u *ServiceItemUsecase) GetActiveServiceItems(ctx context.Context) ([]*entities.ServiceItem, error) {
	return u.serviceItemRepo.GetActive(ctx)
}

func (u *ServiceItemUsecase) GetServiceItemsByCategory(ctx context.Context, category string) ([]*entities.ServiceItem, error) {
	return u.serviceItemRepo.GetByCategory(ctx, category)
}

func (u *ServiceItemUsecase) GetServiceItemsGroupedByCategory(ctx context.Context) (map[string][]dto.ServiceItemResponse, error) {
	items, err := u.serviceItemRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]dto.ServiceItemResponse)
	for _, item := range items {
		grouped[item.Category] = append(grouped[item.Category], dto.ServiceItemResponse{
			ID:              item.ID.String(),
			Name:            item.Name,
			Description:     item.Description,
			Category:        item.Category,
			EstimatedTime:   item.EstimatedTime,
			EstimatedCost:   item.EstimatedCost,
			DisplayOrder:    item.DisplayOrder,
			RequiresBooking: item.RequiresBooking,
		})
	}

	return grouped, nil
}

func (u *ServiceItemUsecase) UpdateServiceItem(ctx context.Context, id uuid.UUID, req *dto.UpdateServiceItemRequest) error {
	serviceItem, err := u.serviceItemRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("service item not found")
	}

	if req.Name != nil {
		serviceItem.Name = *req.Name
	}
	if req.Description != nil {
		serviceItem.Description = *req.Description
	}
	if req.Category != nil {
		serviceItem.Category = *req.Category
	}
	if req.EstimatedTime != nil {
		serviceItem.EstimatedTime = *req.EstimatedTime
	}
	if req.EstimatedCost != nil {
		serviceItem.EstimatedCost = *req.EstimatedCost
	}
	if req.IsActive != nil {
		serviceItem.IsActive = *req.IsActive
	}
	if req.DisplayOrder != nil {
		serviceItem.DisplayOrder = *req.DisplayOrder
	}
	if req.RequiresBooking != nil {
		serviceItem.RequiresBooking = *req.RequiresBooking
	}

	return u.serviceItemRepo.Update(ctx, serviceItem)
}

func (u *ServiceItemUsecase) DeleteServiceItem(ctx context.Context, id uuid.UUID) error {
	return u.serviceItemRepo.Delete(ctx, id)
}
