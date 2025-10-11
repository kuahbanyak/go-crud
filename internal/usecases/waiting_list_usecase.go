package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
)

type WaitingListUsecase struct {
	waitingListRepo repositories.WaitingListRepository
	vehicleRepo     repositories.VehicleRepository
	userRepo        repositories.UserRepository
	settingUsecase  *SettingUsecase
}

func NewWaitingListUsecase(
	waitingListRepo repositories.WaitingListRepository,
	vehicleRepo repositories.VehicleRepository,
	userRepo repositories.UserRepository,
	settingUsecase *SettingUsecase,
) *WaitingListUsecase {
	return &WaitingListUsecase{
		waitingListRepo: waitingListRepo,
		vehicleRepo:     vehicleRepo,
		userRepo:        userRepo,
		settingUsecase:  settingUsecase,
	}
}

func (u *WaitingListUsecase) TakeQueueNumber(ctx context.Context, waitingList *entities.WaitingList) error {
	// Validate vehicle exists
	if u.vehicleRepo != nil {
		_, err := u.vehicleRepo.GetByID(ctx, waitingList.VehicleID)
		if err != nil {
			return errors.New("vehicle not found")
		}
	}

	// Validate customer exists
	_, err := u.userRepo.GetByID(ctx, waitingList.CustomerID)
	if err != nil {
		return errors.New("customer not found")
	}

	// Check if ticket limit has been reached for the service date
	available, _, err := u.CheckTicketAvailability(ctx, waitingList.ServiceDate)
	if err != nil {
		return fmt.Errorf("failed to check ticket availability: %w", err)
	}

	if !available {
		maxTickets := u.settingUsecase.GetMaxTicketsPerDay(ctx)
		return fmt.Errorf("daily ticket limit reached: maximum %d tickets per day (0 remaining)", maxTickets)
	}

	// Get next queue number for the service date
	queueNumber, err := u.waitingListRepo.GetNextQueueNumber(ctx, waitingList.ServiceDate)
	if err != nil {
		return errors.New("failed to generate queue number")
	}

	// Enforce maximum queue number
	maxTickets := u.settingUsecase.GetMaxTicketsPerDay(ctx)
	if queueNumber > maxTickets {
		return fmt.Errorf("cannot create ticket: queue number %d exceeds daily limit of %d tickets", queueNumber, maxTickets)
	}

	waitingList.QueueNumber = queueNumber
	waitingList.Status = entities.WaitingListStatusWaiting

	return u.waitingListRepo.Create(ctx, waitingList)
}

// CheckTicketAvailability checks if new tickets can be created for the given date
func (u *WaitingListUsecase) CheckTicketAvailability(ctx context.Context, serviceDate time.Time) (bool, int, error) {
	// Get all entries for the service date
	entries, err := u.waitingListRepo.GetByServiceDate(ctx, serviceDate)
	if err != nil {
		return false, 0, fmt.Errorf("failed to get entries for date: %w", err)
	}

	// Count only waiting and active entries (not canceled, completed, or no-show)
	activeCount := 0
	for _, entry := range entries {
		if entry.Status == entities.WaitingListStatusWaiting ||
			entry.Status == entities.WaitingListStatusCalled ||
			entry.Status == entities.WaitingListStatusInService {
			activeCount++
		}
	}

	maxTickets := u.settingUsecase.GetMaxTicketsPerDay(ctx)
	available := activeCount < maxTickets
	remaining := maxTickets - activeCount
	if remaining < 0 {
		remaining = 0
	}

	return available, remaining, nil
}

func (u *WaitingListUsecase) GetWaitingList(ctx context.Context, id types.MSSQLUUID) (*entities.WaitingList, error) {
	return u.waitingListRepo.GetByID(ctx, id)
}

func (u *WaitingListUsecase) GetByQueueNumber(ctx context.Context, queueNumber int, serviceDate time.Time) (*entities.WaitingList, error) {
	return u.waitingListRepo.GetByQueueNumber(ctx, queueNumber, serviceDate)
}

func (u *WaitingListUsecase) GetCustomerWaitingLists(ctx context.Context, customerID types.MSSQLUUID) ([]*entities.WaitingList, error) {
	return u.waitingListRepo.GetByCustomerID(ctx, customerID)
}

func (u *WaitingListUsecase) GetTodayQueue(ctx context.Context) ([]*entities.WaitingList, error) {
	today := time.Now()
	return u.waitingListRepo.GetByServiceDate(ctx, today)
}

func (u *WaitingListUsecase) GetQueueByDate(ctx context.Context, serviceDate time.Time) ([]*entities.WaitingList, error) {
	return u.waitingListRepo.GetByServiceDate(ctx, serviceDate)
}

// CheckServiceProgress checks the current progress and status of a service ticket
func (u *WaitingListUsecase) CheckServiceProgress(ctx context.Context, id types.MSSQLUUID) (*ServiceProgressResponse, error) {
	// Get the waiting list entry
	waitingList, err := u.waitingListRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	// Get all tickets for the same service date
	allTickets, err := u.waitingListRepo.GetByServiceDate(ctx, waitingList.ServiceDate)
	if err != nil {
		return nil, errors.New("failed to retrieve queue information")
	}

	// Calculate position and progress
	var currentlyServing int
	var _ int
	var waitingAhead int

	for _, ticket := range allTickets {
		// Find currently being serviced ticket
		if ticket.Status == entities.WaitingListStatusInService {
			currentlyServing = ticket.QueueNumber
		}

		// Count how many are ahead in queue (waiting or called status)
		if ticket.QueueNumber < waitingList.QueueNumber &&
			(ticket.Status == entities.WaitingListStatusWaiting || ticket.Status == entities.WaitingListStatusCalled) {
			waitingAhead++
		}
	}

	// Calculate estimated wait time
	estimatedWaitMin := waitingAhead * 30 // Assume 30 minutes per service on average

	response := &ServiceProgressResponse{
		ID:                   waitingList.ID,
		QueueNumber:          waitingList.QueueNumber,
		Status:               string(waitingList.Status),
		ServiceDate:          waitingList.ServiceDate,
		ServiceType:          waitingList.ServiceType,
		CurrentlyServing:     currentlyServing,
		QueuePosition:        waitingList.QueueNumber,
		WaitingAhead:         waitingAhead,
		EstimatedWaitMinutes: estimatedWaitMin,
		CalledAt:             waitingList.CalledAt,
		ServiceStartAt:       waitingList.ServiceStartAt,
		ServiceEndAt:         waitingList.ServiceEndAt,
		Message:              u.generateProgressMessage(waitingList.Status, waitingAhead, currentlyServing),
	}

	return response, nil
}

// generateProgressMessage creates a user-friendly message based on current status
func (u *WaitingListUsecase) generateProgressMessage(status entities.WaitingListStatus, waitingAhead, currentlyServing int) string {
	switch status {
	case entities.WaitingListStatusWaiting:
		if waitingAhead == 0 {
			return "You're next! Please be ready."
		}
		return fmt.Sprintf("%d customer(s) ahead of you. Currently serving #%d", waitingAhead, currentlyServing)
	case entities.WaitingListStatusCalled:
		return "You've been called! Please proceed to the service area."
	case entities.WaitingListStatusInService:
		return "Your service is currently in progress."
	case entities.WaitingListStatusCompleted:
		return "Your service has been completed. Thank you!"
	case entities.WaitingListStatusCanceled:
		return "This ticket has been canceled."
	case entities.WaitingListStatusNoShow:
		return "Marked as no-show. Please contact us to reschedule."
	default:
		return "Status unknown"
	}
}

// ServiceProgressResponse represents the detailed progress of a service ticket
type ServiceProgressResponse struct {
	ID                   types.MSSQLUUID `json:"id"`
	QueueNumber          int             `json:"queue_number"`
	Status               string          `json:"status"`
	ServiceDate          time.Time       `json:"service_date"`
	ServiceType          string          `json:"service_type"`
	CurrentlyServing     int             `json:"currently_serving"`
	QueuePosition        int             `json:"queue_position"`
	WaitingAhead         int             `json:"waiting_ahead"`
	EstimatedWaitMinutes int             `json:"estimated_wait_minutes"`
	CalledAt             *time.Time      `json:"called_at,omitempty"`
	ServiceStartAt       *time.Time      `json:"service_start_at,omitempty"`
	ServiceEndAt         *time.Time      `json:"service_end_at,omitempty"`
	Message              string          `json:"message"`
}

func (u *WaitingListUsecase) CallCustomer(ctx context.Context, id types.MSSQLUUID) error {
	waitingList, err := u.waitingListRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if waitingList.Status != entities.WaitingListStatusWaiting {
		return errors.New("can only call customers in waiting status")
	}

	now := time.Now()
	waitingList.Status = entities.WaitingListStatusCalled
	waitingList.CalledAt = &now

	return u.waitingListRepo.Update(ctx, waitingList)
}

func (u *WaitingListUsecase) StartService(ctx context.Context, id types.MSSQLUUID) error {
	waitingList, err := u.waitingListRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if waitingList.Status != entities.WaitingListStatusCalled {
		return errors.New("customer must be called before starting service")
	}

	now := time.Now()
	waitingList.Status = entities.WaitingListStatusInService
	waitingList.ServiceStartAt = &now

	return u.waitingListRepo.Update(ctx, waitingList)
}

func (u *WaitingListUsecase) CompleteService(ctx context.Context, id types.MSSQLUUID) error {
	waitingList, err := u.waitingListRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if waitingList.Status != entities.WaitingListStatusInService {
		return errors.New("service must be in progress to complete")
	}

	now := time.Now()
	waitingList.Status = entities.WaitingListStatusCompleted
	waitingList.ServiceEndAt = &now

	return u.waitingListRepo.Update(ctx, waitingList)
}

func (u *WaitingListUsecase) CancelQueue(ctx context.Context, id types.MSSQLUUID) error {
	waitingList, err := u.waitingListRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if waitingList.Status == entities.WaitingListStatusCompleted {
		return errors.New("cannot cancel completed service")
	}

	waitingList.Status = entities.WaitingListStatusCanceled
	return u.waitingListRepo.Update(ctx, waitingList)
}

func (u *WaitingListUsecase) MarkNoShow(ctx context.Context, id types.MSSQLUUID) error {
	waitingList, err := u.waitingListRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if waitingList.Status != entities.WaitingListStatusCalled {
		return errors.New("can only mark no-show for called customers")
	}

	waitingList.Status = entities.WaitingListStatusNoShow
	return u.waitingListRepo.Update(ctx, waitingList)
}

func (u *WaitingListUsecase) GetWaitingCount(ctx context.Context, serviceDate time.Time) (int, error) {
	waitingLists, err := u.waitingListRepo.GetByStatus(ctx, entities.WaitingListStatusWaiting, serviceDate)
	if err != nil {
		return 0, err
	}
	return len(waitingLists), nil
}

func (u *WaitingListUsecase) UpdateWaitingList(ctx context.Context, id types.MSSQLUUID, updates *entities.WaitingList) error {
	existing, err := u.waitingListRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("waiting list entry not found")
	}

	updates.ID = id
	return u.waitingListRepo.Update(ctx, updates)
}
