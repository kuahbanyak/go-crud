package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/domain/repositories"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/logger"
	"github.com/kuahbanyak/go-crud/internal/usecases"
)

type DailyCleanupJob struct {
	waitingListRepo repositories.WaitingListRepository
	jobUsecase      *usecases.JobUsecase
}

func NewDailyCleanupJob(waitingListRepo repositories.WaitingListRepository, jobUsecase *usecases.JobUsecase) *DailyCleanupJob {
	return &DailyCleanupJob{
		waitingListRepo: waitingListRepo,
		jobUsecase:      jobUsecase,
	}
}
func (j *DailyCleanupJob) Name() string {
	return "DailyWaitingListCleanup"
}

func (j *DailyCleanupJob) Schedule() string {
	if j.jobUsecase != nil {
		ctx := context.Background()
		schedule := j.jobUsecase.GetJobSchedule(ctx, "DailyWaitingListCleanup")
		if schedule != "" {
			return schedule
		}
	}
	return "0 0 * * *"
}

func (j *DailyCleanupJob) Run(ctx context.Context) error {
	// Check if job is enabled in database
	if j.jobUsecase != nil && !j.jobUsecase.IsJobEnabled(ctx, "DailyWaitingListCleanup") {
		logger.Info("Daily cleanup job is disabled in database, skipping...")
		return nil
	}

	logger.Info("Running daily waiting list cleanup job...")
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	if err := j.cleanupOldEntries(ctx, startOfDay); err != nil {
		logger.Error(fmt.Sprintf("Failed to cleanup old entries: %v", err))
		return err
	}
	if err := j.enforceTicketLimit(ctx, today); err != nil {
		logger.Error(fmt.Sprintf("Failed to enforce ticket limit: %v", err))
		return err
	}
	logger.Info("Daily cleanup completed successfully")
	return nil
}
func (j *DailyCleanupJob) cleanupOldEntries(ctx context.Context, today time.Time) error {
	retentionDays := 7

	oldDate := today.AddDate(0, 0, -retentionDays)
	logger.Info(fmt.Sprintf("Cleaning up entries older than %s (%d days)", oldDate.Format("2006-01-02"), retentionDays))
	completedEntries, err := j.waitingListRepo.GetByStatus(ctx, entities.WaitingListStatusCompleted, oldDate)
	if err != nil {
		return fmt.Errorf("failed to get completed entries: %w", err)
	}
	canceledEntries, err := j.waitingListRepo.GetByStatus(ctx, entities.WaitingListStatusCanceled, oldDate)
	if err != nil {
		return fmt.Errorf("failed to get canceled entries: %w", err)
	}
	noShowEntries, err := j.waitingListRepo.GetByStatus(ctx, entities.WaitingListStatusNoShow, oldDate)
	if err != nil {
		return fmt.Errorf("failed to get no-show entries: %w", err)
	}
	totalCleaned := 0
	for _, entry := range append(append(completedEntries, canceledEntries...), noShowEntries...) {
		if entry.ServiceDate.Before(oldDate) {
			if err := j.waitingListRepo.Delete(ctx, entry.ID); err != nil {
				logger.Error(fmt.Sprintf("Failed to delete entry %s: %v", entry.ID, err))
			} else {
				totalCleaned++
			}
		}
	}
	logger.Info(fmt.Sprintf("Cleaned up %d old entries", totalCleaned))
	return nil
}
func (j *DailyCleanupJob) enforceTicketLimit(ctx context.Context, today time.Time) error {
	// Get weekly limit - hardcoded default
	maxTicketsPerWeek := 20

	// Calculate week boundaries
	weekStart := j.getWeekStart(today)
	weekEnd := j.getWeekEnd(today)

	// Get all entries for the week using GetByWeekRange
	weekEntries, err := j.waitingListRepo.GetByWeekRange(ctx, weekStart, weekEnd)
	if err != nil {
		return fmt.Errorf("failed to get week's entries: %w", err)
	}

	// Count active tickets in the week
	activeCount := 0
	var activeEntries []*entities.WaitingList
	for _, entry := range weekEntries {
		if entry.Status == entities.WaitingListStatusWaiting ||
			entry.Status == entities.WaitingListStatusCalled ||
			entry.Status == entities.WaitingListStatusInService {
			activeEntries = append(activeEntries, entry)
			activeCount++
		}
	}

	// Check if weekly limit is exceeded
	if activeCount > maxTicketsPerWeek {
		excessCount := activeCount - maxTicketsPerWeek
		logger.Info(fmt.Sprintf("Found %d active tickets in week %s to %s, canceling %d excess tickets",
			activeCount, weekStart.Format("Jan 02"), weekEnd.Format("Jan 02"), excessCount))

		// Cancel excess tickets (oldest waiting tickets first)
		for i := 0; i < excessCount && i < len(activeEntries); i++ {
			entry := activeEntries[i]
			if entry.Status == entities.WaitingListStatusWaiting {
				entry.Status = entities.WaitingListStatusCanceled
				entry.Notes = fmt.Sprintf("%s [Auto-canceled: Weekly limit of %d tickets exceeded]", entry.Notes, maxTicketsPerWeek)
				if err := j.waitingListRepo.Update(ctx, entry); err != nil {
					logger.Error(fmt.Sprintf("Failed to cancel excess entry %s: %v", entry.ID, err))
				} else {
					logger.Info(fmt.Sprintf("Canceled excess ticket #%d for customer %s", entry.QueueNumber, entry.CustomerID))
				}
			}
		}
	} else {
		logger.Info(fmt.Sprintf("Current active tickets for week: %d/%d (within limit)", activeCount, maxTicketsPerWeek))
	}
	return nil
}

// getWeekStart returns the start of the week (Monday 00:00:00)
func (j *DailyCleanupJob) getWeekStart(t time.Time) time.Time {
	weekday := int(t.Weekday())
	daysToMonday := (weekday + 6) % 7
	monday := t.AddDate(0, 0, -daysToMonday)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, t.Location())
}

// getWeekEnd returns the end of the week (Sunday 23:59:59)
func (j *DailyCleanupJob) getWeekEnd(t time.Time) time.Time {
	weekStart := j.getWeekStart(t)
	sunday := weekStart.AddDate(0, 0, 6)
	return time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 999999999, t.Location())
}
