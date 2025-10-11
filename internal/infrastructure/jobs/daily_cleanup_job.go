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
	settingUsecase  *usecases.SettingUsecase
}

func NewDailyCleanupJob(waitingListRepo repositories.WaitingListRepository, settingUsecase *usecases.SettingUsecase) *DailyCleanupJob {
	return &DailyCleanupJob{
		waitingListRepo: waitingListRepo,
		settingUsecase:  settingUsecase,
	}
}

func (j *DailyCleanupJob) Name() string {
	return "DailyWaitingListCleanup"
}

func (j *DailyCleanupJob) Schedule() string {
	// Get schedule from settings
	if j.settingUsecase != nil {
		ctx := context.Background()
		schedule := j.settingUsecase.GetJobSchedule(ctx)
		if schedule != "" {
			return schedule
		}
	}
	// Default: Run every day at midnight (00:00)
	return "0 0 * * *"
}

func (j *DailyCleanupJob) Run(ctx context.Context) error {
	// Check if job is enabled
	if j.settingUsecase != nil && !j.settingUsecase.IsCleanupJobEnabled(ctx) {
		logger.Info("Daily cleanup job is disabled in settings, skipping...")
		return nil
	}

	logger.Info("Running daily waiting list cleanup job...")

	// Get today's date
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	// 1. Clean up old completed/canceled entries from previous days
	if err := j.cleanupOldEntries(ctx, startOfDay); err != nil {
		logger.Error(fmt.Sprintf("Failed to cleanup old entries: %v", err))
		return err
	}

	// 2. Cancel any excess waiting entries beyond the limit
	if err := j.enforceTicketLimit(ctx, today); err != nil {
		logger.Error(fmt.Sprintf("Failed to enforce ticket limit: %v", err))
		return err
	}

	logger.Info("Daily cleanup completed successfully")
	return nil
}

func (j *DailyCleanupJob) cleanupOldEntries(ctx context.Context, today time.Time) error {
	// Get retention days from settings
	retentionDays := 7
	if j.settingUsecase != nil {
		retentionDays = j.settingUsecase.GetCleanupRetentionDays(ctx)
	}

	oldDate := today.AddDate(0, 0, -retentionDays)

	logger.Info(fmt.Sprintf("Cleaning up entries older than %s (%d days)", oldDate.Format("2006-01-02"), retentionDays))

	// Get old completed entries
	completedEntries, err := j.waitingListRepo.GetByStatus(ctx, entities.WaitingListStatusCompleted, oldDate)
	if err != nil {
		return fmt.Errorf("failed to get completed entries: %w", err)
	}

	// Get old canceled entries
	canceledEntries, err := j.waitingListRepo.GetByStatus(ctx, entities.WaitingListStatusCanceled, oldDate)
	if err != nil {
		return fmt.Errorf("failed to get canceled entries: %w", err)
	}

	// Get old no-show entries
	noShowEntries, err := j.waitingListRepo.GetByStatus(ctx, entities.WaitingListStatusNoShow, oldDate)
	if err != nil {
		return fmt.Errorf("failed to get no-show entries: %w", err)
	}

	totalCleaned := 0

	// Delete old entries (soft delete)
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
	// Get max tickets from settings
	maxTickets := 10
	if j.settingUsecase != nil {
		maxTickets = j.settingUsecase.GetMaxTicketsPerDay(ctx)
	}

	// Get all waiting entries for today
	todayEntries, err := j.waitingListRepo.GetByServiceDate(ctx, today)
	if err != nil {
		return fmt.Errorf("failed to get today's entries: %w", err)
	}

	// Count only waiting entries
	waitingCount := 0
	var waitingEntries []*entities.WaitingList
	for _, entry := range todayEntries {
		if entry.Status == entities.WaitingListStatusWaiting {
			waitingEntries = append(waitingEntries, entry)
			waitingCount++
		}
	}

	// If more than max waiting tickets, cancel the excess (oldest created first)
	if waitingCount > maxTickets {
		excessCount := waitingCount - maxTickets
		logger.Info(fmt.Sprintf("Found %d waiting tickets, canceling %d excess tickets", waitingCount, excessCount))

		// Sort by creation time and cancel the oldest excess ones
		for i := 0; i < excessCount && i < len(waitingEntries); i++ {
			entry := waitingEntries[i]
			entry.Status = entities.WaitingListStatusCanceled
			entry.Notes = fmt.Sprintf("%s [Auto-canceled: Daily limit of %d tickets exceeded]", entry.Notes, maxTickets)

			if err := j.waitingListRepo.Update(ctx, entry); err != nil {
				logger.Error(fmt.Sprintf("Failed to cancel excess entry %s: %v", entry.ID, err))
			} else {
				logger.Info(fmt.Sprintf("Canceled excess ticket #%d for customer %s", entry.QueueNumber, entry.CustomerID))
			}
		}
	} else {
		logger.Info(fmt.Sprintf("Current waiting tickets: %d/%d (within limit)", waitingCount, maxTickets))
	}

	return nil
}
