package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/shared/utils"
)

type WaitingListUsecase struct {
	waitingListRepo repositories.WaitingListRepository
	vehicleRepo     repositories.VehicleRepository
	userRepo        repositories.UserRepository
	serviceItemRepo repositories.ServiceItemRepository
	vehicleUsecase  *VehicleUseCase
	jobUsecase      *JobUsecase
}

func NewWaitingListUsecase(
	waitingListRepo repositories.WaitingListRepository,
	vehicleRepo repositories.VehicleRepository,
	userRepo repositories.UserRepository,
	vehicleUsecase *VehicleUseCase,
	serviceItemRepo repositories.ServiceItemRepository,
	jobUsecase *JobUsecase,
) *WaitingListUsecase {
	return &WaitingListUsecase{
		waitingListRepo: waitingListRepo,
		vehicleRepo:     vehicleRepo,
		userRepo:        userRepo,
		serviceItemRepo: serviceItemRepo,
		vehicleUsecase:  vehicleUsecase,
		jobUsecase:      jobUsecase,
	}
}

func (u *WaitingListUsecase) TakeQueueNumber(ctx context.Context, waitingList *entities.WaitingList) error {
	// Check if any job is active (system is accepting tickets)
	if u.jobUsecase != nil {
		hasActiveJob, err := u.jobUsecase.HasAnyActiveJob(ctx)
		if err != nil {
			return fmt.Errorf("failed to check job status: %w", err)
		}
		if !hasActiveJob {
			return errors.New("the ticket system is currently not accepting new bookings. Please contact the administrator to activate the system")
		}
	}

	if u.vehicleRepo != nil {
		_, err := u.vehicleRepo.GetByID(ctx, waitingList.VehicleID)
		if err != nil {
			return errors.New("vehicle not found")
		}
	}
	_, err := u.userRepo.GetByID(ctx, waitingList.CustomerID)
	if err != nil {
		return errors.New("customer not found")
	}

	// If service item is selected, populate estimated time and cost
	if waitingList.ServiceItemID != nil && u.serviceItemRepo != nil {
		serviceItem, err := u.serviceItemRepo.GetByID(ctx, *waitingList.ServiceItemID)
		if err != nil {
			return errors.New("selected service item not found")
		}

		// If not active, reject
		if !serviceItem.IsActive {
			return errors.New("selected service item is not available")
		}

		// Auto-populate estimated time and cost from service item
		if waitingList.EstimatedTime == 0 {
			waitingList.EstimatedTime = serviceItem.EstimatedTime
		}
		waitingList.EstimatedCost = serviceItem.EstimatedCost

		// Use service item name as service type if not provided
		if waitingList.ServiceType == "" {
			waitingList.ServiceType = serviceItem.Name
		}
	}

	// Check if the service date is in a past week
	now := utils.NowWIB()
	if !u.IsInCurrentOrFutureWeek(waitingList.ServiceDate, now) {
		return errors.New("cannot book tickets for previous weeks. Please select a date in the current week or a future week")
	}

	// Check weekly ticket availability
	available, remaining, weekStart, weekEnd, err := u.CheckWeeklyTicketAvailability(ctx, waitingList.ServiceDate)
	if err != nil {
		return fmt.Errorf("failed to check weekly ticket availability: %w", err)
	}
	if !available {
		maxTickets := 20 // Default weekly ticket limit
		return fmt.Errorf("weekly ticket limit reached for week of %s to %s: maximum %d tickets per week (%d remaining)",
			weekStart.Format("Jan 02"), weekEnd.Format("Jan 02"), maxTickets, remaining)
	}

	queueNumber, err := u.waitingListRepo.GetNextQueueNumber(ctx, waitingList.ServiceDate)
	if err != nil {
		return errors.New("failed to generate queue number")
	}

	waitingList.QueueNumber = queueNumber
	waitingList.Status = entities.WaitingListStatusWaiting
	return u.waitingListRepo.Create(ctx, waitingList)
}

// CheckWeeklyTicketAvailability checks if there are available tickets for the week containing serviceDate
func (u *WaitingListUsecase) CheckWeeklyTicketAvailability(ctx context.Context, serviceDate time.Time) (available bool, remaining int, weekStart time.Time, weekEnd time.Time, err error) {
	// Get week boundaries
	weekStart = u.GetWeekStart(serviceDate)
	weekEnd = u.GetWeekEnd(serviceDate)

	// Get all entries for the week
	entries, err := u.waitingListRepo.GetByWeekRange(ctx, weekStart, weekEnd)
	if err != nil {
		return false, 0, weekStart, weekEnd, fmt.Errorf("failed to get entries for week: %w", err)
	}

	// Count active tickets in the week
	activeCount := 0
	for _, entry := range entries {
		if entry.Status == entities.WaitingListStatusWaiting ||
			entry.Status == entities.WaitingListStatusCalled ||
			entry.Status == entities.WaitingListStatusInService {
			activeCount++
		}
	}

	maxTickets := 20 // Default weekly ticket limit
	available = activeCount < maxTickets
	remaining = maxTickets - activeCount
	if remaining < 0 {
		remaining = 0
	}

	return available, remaining, weekStart, weekEnd, nil
}

// GetWeekStart returns the start of the week (Monday 00:00:00) for the given date
func (u *WaitingListUsecase) GetWeekStart(t time.Time) time.Time {
	weekday := int(t.Weekday())
	daysToMonday := (weekday + 6) % 7
	monday := t.AddDate(0, 0, -daysToMonday)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, t.Location())
}

// GetWeekEnd returns the end of the week (Sunday 23:59:59) for the given date
func (u *WaitingListUsecase) GetWeekEnd(t time.Time) time.Time {
	weekStart := u.GetWeekStart(t)
	sunday := weekStart.AddDate(0, 0, 6)
	return time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 999999999, t.Location())
}

// IsInCurrentOrFutureWeek checks if the given date is in the current week or a future week
func (u *WaitingListUsecase) IsInCurrentOrFutureWeek(checkDate time.Time, referenceDate time.Time) bool {
	checkWeekStart := u.GetWeekStart(checkDate)
	referenceWeekStart := u.GetWeekStart(referenceDate)
	return !checkWeekStart.Before(referenceWeekStart)
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
	today := utils.NowWIB()
	return u.waitingListRepo.GetByServiceDate(ctx, today)
}
func (u *WaitingListUsecase) GetQueueByDate(ctx context.Context, serviceDate time.Time) ([]*entities.WaitingList, error) {
	return u.waitingListRepo.GetByServiceDate(ctx, serviceDate)
}
func (u *WaitingListUsecase) CheckServiceProgress(ctx context.Context, id types.MSSQLUUID) (*ServiceProgressResponse, error) {
	waitingList, err := u.waitingListRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("ticket not found")
	}
	allTickets, err := u.waitingListRepo.GetByServiceDate(ctx, waitingList.ServiceDate)
	if err != nil {
		return nil, errors.New("failed to retrieve queue information")
	}
	var currentlyServing int
	var _ int
	var waitingAhead int
	for _, ticket := range allTickets {
		if ticket.Status == entities.WaitingListStatusInService {
			currentlyServing = ticket.QueueNumber
		}
		if ticket.QueueNumber < waitingList.QueueNumber &&
			(ticket.Status == entities.WaitingListStatusWaiting || ticket.Status == entities.WaitingListStatusCalled) {
			waitingAhead++
		}
	}
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
	now := utils.NowWIB()
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
	now := utils.NowWIB()
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
	now := utils.NowWIB()
	waitingList.Status = entities.WaitingListStatusCompleted
	waitingList.ServiceEndAt = &now
	return u.waitingListRepo.Update(ctx, waitingList)
}
func (u *WaitingListUsecase) CancelQueue(ctx context.Context, id types.MSSQLUUID) error {
	waitingList, err := u.waitingListRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if the ticket can be cancelled
	if waitingList.Status == entities.WaitingListStatusCompleted {
		return errors.New("cannot cancel completed service")
	}
	if waitingList.Status == entities.WaitingListStatusCanceled {
		return errors.New("ticket is already cancelled")
	}

	// Store the cancelled queue number for returning the ticket
	cancelledQueueNumber := waitingList.QueueNumber
	serviceDate := waitingList.ServiceDate

	// Mark the ticket as cancelled
	waitingList.Status = entities.WaitingListStatusCanceled
	if err := u.waitingListRepo.Update(ctx, waitingList); err != nil {
		return fmt.Errorf("failed to cancel ticket: %w", err)
	}

	// Return the ticket by reordering queue numbers
	// Get all waiting and called tickets after the cancelled one
	allTickets, err := u.waitingListRepo.GetByServiceDate(ctx, serviceDate)
	if err != nil {
		return nil // Ticket is already cancelled, just log this error
	}

	// Reorder queue numbers for tickets that come after the cancelled one
	for _, ticket := range allTickets {
		// Only reorder tickets that are waiting or called and have a higher queue number
		if (ticket.Status == entities.WaitingListStatusWaiting || ticket.Status == entities.WaitingListStatusCalled) &&
			ticket.QueueNumber > cancelledQueueNumber {
			ticket.QueueNumber--
			if err := u.waitingListRepo.Update(ctx, ticket); err != nil {
				// Log error but continue with other tickets
				continue
			}
		}
	}

	return nil
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
	if err != nil || existing == nil {
		return errors.New("waiting list entry not found")
	}

	// Update only provided fields
	if updates.ServiceType != "" {
		existing.ServiceType = updates.ServiceType
	}
	if updates.EstimatedTime > 0 {
		existing.EstimatedTime = updates.EstimatedTime
	}
	if updates.Notes != "" {
		existing.Notes = updates.Notes
	}
	if updates.MechanicNotes != "" {
		existing.MechanicNotes = updates.MechanicNotes
	}
	if updates.Status != "" {
		existing.Status = updates.Status
	}

	return u.waitingListRepo.Update(ctx, existing)
}

// AssignMechanicToQueue assigns a mechanic to service a queue entry
func (u *WaitingListUsecase) AssignMechanicToQueue(ctx context.Context, queueID types.MSSQLUUID, mechanicID types.MSSQLUUID) error {
	// Get the waiting list entry
	waitingList, err := u.waitingListRepo.GetByID(ctx, queueID)
	if err != nil {
		return errors.New("queue entry not found")
	}

	// Check if mechanic exists and has mechanic role
	mechanic, err := u.userRepo.GetByID(ctx, mechanicID)
	if err != nil {
		return errors.New("mechanic not found")
	}

	// Check if user is a mechanic or admin (assuming roles are checked elsewhere)
	_ = mechanic // Use mechanic variable to avoid unused error

	// Check if queue is in a valid state for assignment
	if waitingList.Status != entities.WaitingListStatusWaiting &&
		waitingList.Status != entities.WaitingListStatusCalled {
		return errors.New("can only assign mechanic to waiting or called queues")
	}

	// Assign mechanic
	waitingList.MechanicID = &mechanicID

	return u.waitingListRepo.Update(ctx, waitingList)
}

// GetAvailableQueues returns queues that are waiting to be serviced (for mechanics to choose from)
func (u *WaitingListUsecase) GetAvailableQueues(ctx context.Context, serviceDate time.Time) ([]*entities.WaitingList, error) {
	// Get all queues for the date
	allQueues, err := u.waitingListRepo.GetByServiceDate(ctx, serviceDate)
	if err != nil {
		return nil, err
	}

	// Filter for queues that are waiting or called (available for service)
	availableQueues := make([]*entities.WaitingList, 0)
	for _, queue := range allQueues {
		if queue.Status == entities.WaitingListStatusWaiting ||
			queue.Status == entities.WaitingListStatusCalled {
			availableQueues = append(availableQueues, queue)
		}
	}

	return availableQueues, nil
}
