package usecases
import (
	"context"
	"errors"
	"time"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)
type MaintenanceItemUsecase struct {
	maintenanceItemRepo repositories.MaintenanceItemRepository
	waitingListRepo     repositories.WaitingListRepository
	userRepo            repositories.UserRepository
}
func NewMaintenanceItemUsecase(
	maintenanceItemRepo repositories.MaintenanceItemRepository,
	waitingListRepo repositories.WaitingListRepository,
	userRepo repositories.UserRepository,
) *MaintenanceItemUsecase {
	return &MaintenanceItemUsecase{
		maintenanceItemRepo: maintenanceItemRepo,
		waitingListRepo:     waitingListRepo,
		userRepo:            userRepo,
	}
}
func (u *MaintenanceItemUsecase) CreateInitialItems(ctx context.Context, waitingListID types.MSSQLUUID, requests []dto.CreateMaintenanceItemRequest) error {
	_, err := u.waitingListRepo.GetByID(ctx, waitingListID)
	if err != nil {
		return errors.New("waiting list not found")
	}
	items := make([]*entities.MaintenanceItem, len(requests))
	for i, req := range requests {
		items[i] = &entities.MaintenanceItem{
			WaitingListID:    waitingListID,
			ItemType:         entities.MaintenanceItemTypeInitial,
			Status:           entities.MaintenanceItemStatusPending,
			Category:         req.Category,
			Name:             req.Name,
			Description:      req.Description,
			EstimatedCost:    req.EstimatedCost,
			Priority:         "normal",
			RequiresApproval: false, // Initial items don't need approval
		}
	}
	return u.maintenanceItemRepo.CreateMany(ctx, items)
}
func (u *MaintenanceItemUsecase) AddDiscoveredItem(ctx context.Context, mechanicID types.MSSQLUUID, req dto.AddDiscoveredItemRequest) (*entities.MaintenanceItem, error) {
	waitingList, err := u.waitingListRepo.GetByID(ctx, req.WaitingListID)
	if err != nil {
		return nil, errors.New("waiting list not found")
	}
	if waitingList.Status != entities.WaitingListStatusInService {
		return nil, errors.New("service must be in progress to add discovered items")
	}
	_, err = u.userRepo.GetByID(ctx, mechanicID)
	if err != nil {
		return nil, errors.New("mechanic not found")
	}
	now := time.Now()
	item := &entities.MaintenanceItem{
		WaitingListID:    req.WaitingListID,
		MechanicID:       &mechanicID,
		ItemType:         entities.MaintenanceItemTypeDiscovered,
		Status:           entities.MaintenanceItemStatusInspected,
		Category:         req.Category,
		Name:             req.Name,
		Description:      req.Description,
		Priority:         req.Priority,
		EstimatedCost:    req.EstimatedCost,
		LaborHours:       req.LaborHours,
		RequiresApproval: req.RequiresApproval,
		ImageURL:         req.ImageURL,
		Notes:            req.Notes,
		InspectedAt:      &now,
	}
	err = u.maintenanceItemRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}
func (u *MaintenanceItemUsecase) GetItemsByWaitingList(ctx context.Context, waitingListID types.MSSQLUUID) (*dto.MaintenanceItemListResponse, error) {
	items, err := u.maintenanceItemRepo.GetByWaitingListID(ctx, waitingListID)
	if err != nil {
		return nil, err
	}
	estimated, actual, err := u.maintenanceItemRepo.GetTotalCost(ctx, waitingListID)
	if err != nil {
		return nil, err
	}
	counts, err := u.maintenanceItemRepo.CountByStatus(ctx, waitingListID)
	if err != nil {
		return nil, err
	}
	response := &dto.MaintenanceItemListResponse{
		Items:                u.buildItemResponses(items),
		Total:                len(items),
		TotalEstimatedCost:   estimated,
		TotalActualCost:      actual,
		PendingApprovalCount: counts["inspected"],
		CompletedCount:       counts["completed"],
	}
	return response, nil
}
func (u *MaintenanceItemUsecase) GetInspectionSummary(ctx context.Context, waitingListID types.MSSQLUUID, customerID types.MSSQLUUID) (*dto.MaintenanceInspectionSummary, error) {
	waitingList, err := u.waitingListRepo.GetByID(ctx, waitingListID)
	if err != nil {
		return nil, errors.New("waiting list not found")
	}
	if waitingList.CustomerID != customerID {
		return nil, errors.New("unauthorized: not your service ticket")
	}
	initialItems, err := u.maintenanceItemRepo.GetInitialItems(ctx, waitingListID)
	if err != nil {
		return nil, err
	}
	discoveredItems, err := u.maintenanceItemRepo.GetDiscoveredItems(ctx, waitingListID)
	if err != nil {
		return nil, err
	}
	estimated, _, err := u.maintenanceItemRepo.GetTotalCost(ctx, waitingListID)
	if err != nil {
		return nil, err
	}
	requiresApproval := false
	for _, item := range discoveredItems {
		if item.RequiresApproval && item.Status == entities.MaintenanceItemStatusInspected {
			requiresApproval = true
			break
		}
	}
	summary := &dto.MaintenanceInspectionSummary{
		WaitingListID:      waitingListID,
		QueueNumber:        waitingList.QueueNumber,
		VehicleBrand:       waitingList.Vehicle.Brand,
		VehicleModel:       waitingList.Vehicle.Model,
		LicensePlate:       waitingList.Vehicle.LicensePlate,
		InitialItems:       u.buildItemResponses(initialItems),
		DiscoveredItems:    u.buildItemResponses(discoveredItems),
		TotalEstimatedCost: estimated,
		RequiresApproval:   requiresApproval,
		InspectedAt:        time.Now(),
	}
	return summary, nil
}
func (u *MaintenanceItemUsecase) ApproveItems(ctx context.Context, customerID types.MSSQLUUID, req dto.ApproveMaintenanceItemRequest) error {
	for _, itemID := range req.ItemIDs {
		item, err := u.maintenanceItemRepo.GetByID(ctx, itemID)
		if err != nil {
			return errors.New("item not found")
		}
		waitingList, err := u.waitingListRepo.GetByID(ctx, item.WaitingListID)
		if err != nil {
			return errors.New("waiting list not found")
		}
		if waitingList.CustomerID != customerID {
			return errors.New("unauthorized: not your maintenance item")
		}
		if item.Status != entities.MaintenanceItemStatusInspected {
			return errors.New("item is not in inspected status")
		}
	}
	if req.Approve {
		return u.maintenanceItemRepo.ApproveItems(ctx, req.ItemIDs)
	}
	return u.maintenanceItemRepo.RejectItems(ctx, req.ItemIDs)
}
func (u *MaintenanceItemUsecase) UpdateItem(ctx context.Context, itemID types.MSSQLUUID, req dto.UpdateMaintenanceItemRequest) error {
	item, err := u.maintenanceItemRepo.GetByID(ctx, itemID)
	if err != nil {
		return errors.New("item not found")
	}
	if req.Status != "" {
		item.Status = entities.MaintenanceItemStatus(req.Status)
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.EstimatedCost != nil {
		item.EstimatedCost = *req.EstimatedCost
	}
	if req.ActualCost != nil {
		item.ActualCost = *req.ActualCost
	}
	if req.LaborHours != nil {
		item.LaborHours = *req.LaborHours
	}
	if req.Priority != "" {
		item.Priority = req.Priority
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}
	return u.maintenanceItemRepo.Update(ctx, item)
}
func (u *MaintenanceItemUsecase) CompleteItem(ctx context.Context, itemID types.MSSQLUUID, actualCost float64) error {
	item, err := u.maintenanceItemRepo.GetByID(ctx, itemID)
	if err != nil {
		return errors.New("item not found")
	}
	if item.Status != entities.MaintenanceItemStatusApproved && item.Status != entities.MaintenanceItemStatusPending {
		return errors.New("item must be approved or pending to complete")
	}
	now := time.Now()
	item.Status = entities.MaintenanceItemStatusCompleted
	item.ActualCost = actualCost
	item.CompletedAt = &now
	return u.maintenanceItemRepo.Update(ctx, item)
}
func (u *MaintenanceItemUsecase) DeleteItem(ctx context.Context, itemID types.MSSQLUUID) error {
	return u.maintenanceItemRepo.Delete(ctx, itemID)
}
func (u *MaintenanceItemUsecase) buildItemResponses(items []*entities.MaintenanceItem) []dto.MaintenanceItemResponse {
	responses := make([]dto.MaintenanceItemResponse, len(items))
	for i, item := range items {
		responses[i] = dto.MaintenanceItemResponse{
			ID:               item.ID,
			WaitingListID:    item.WaitingListID,
			MechanicID:       item.MechanicID,
			ItemType:         string(item.ItemType),
			Status:           string(item.Status),
			Category:         item.Category,
			Name:             item.Name,
			Description:      item.Description,
			Priority:         item.Priority,
			EstimatedCost:    item.EstimatedCost,
			ActualCost:       item.ActualCost,
			LaborHours:       item.LaborHours,
			RequiresApproval: item.RequiresApproval,
			ImageURL:         item.ImageURL,
			Notes:            item.Notes,
			InspectedAt:      item.InspectedAt,
			ApprovedAt:       item.ApprovedAt,
			CompletedAt:      item.CompletedAt,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
		}
		if item.Mechanic != nil {
			responses[i].MechanicName = item.Mechanic.Name
		}
	}
	return responses
}

