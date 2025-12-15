package usecases

import (
	"context"
	"database/sql"
	"time"

	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
)

type AnalyticsUsecase struct {
	db *sql.DB
}

func NewAnalyticsUsecase(db *sql.DB) *AnalyticsUsecase {
	return &AnalyticsUsecase{db: db}
}

func (u *AnalyticsUsecase) GetOverview(ctx context.Context) (*dto.AnalyticsOverviewResponse, error) {
	overview := &dto.AnalyticsOverviewResponse{
		Timestamp: time.Now(),
	}

	todayStart := time.Now().Truncate(24 * time.Hour)
	query := `SELECT COALESCE(SUM(total_amount), 0) FROM invoices WHERE created_at >= @p1 AND status = @p2 AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", todayStart), sql.Named("p2", entities.InvoiceStatusPaid)).Scan(&overview.TodayRevenue)

	query = `SELECT COALESCE(SUM(total_amount), 0) FROM invoices WHERE status = @p1 AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", entities.InvoiceStatusPaid)).Scan(&overview.TotalRevenue)

	query = `SELECT COUNT(*) FROM users WHERE role = 'customer' AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query).Scan(&overview.TotalCustomers)

	query = `SELECT COUNT(*) FROM waiting_lists WHERE status = 'in_service' AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query).Scan(&overview.ActiveServices)

	query = `SELECT COUNT(*) FROM waiting_lists WHERE status = 'completed' AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query).Scan(&overview.CompletedServices)

	query = `SELECT COUNT(*) FROM invoices WHERE status = 'pending' AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query).Scan(&overview.PendingInvoices)

	query = `SELECT COUNT(*) FROM waiting_lists WHERE CAST(service_date AS DATE) = CAST(@p1 AS DATE) AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", time.Now())).Scan(&overview.TodayQueue)

	query = `SELECT AVG(DATEDIFF(MINUTE, created_at, called_at)) FROM waiting_lists WHERE called_at IS NOT NULL AND created_at >= @p1 AND deleted_at IS NULL`
	var avgWait sql.NullFloat64
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", todayStart)).Scan(&avgWait)
	if avgWait.Valid {
		overview.AverageWaitTime = avgWait.Float64
	}

	return overview, nil
}

func (u *AnalyticsUsecase) GetRevenueStats(ctx context.Context, period string) (*dto.RevenueStatsResponse, error) {
	stats := &dto.RevenueStatsResponse{
		Period: period,
		Data:   []dto.RevenueDataPoint{},
	}

	var startDate time.Time
	switch period {
	case "daily":
		startDate = time.Now().AddDate(0, 0, -30)
	case "weekly":
		startDate = time.Now().AddDate(0, 0, -90)
	case "monthly":
		startDate = time.Now().AddDate(0, -12, 0)
	default:
		startDate = time.Now().AddDate(0, 0, -30)
	}

	query := `
		SELECT CONVERT(VARCHAR, created_at, 23) as date, COALESCE(SUM(total_amount), 0) as amount, COUNT(*) as count
		FROM invoices
		WHERE status = @p1 AND created_at >= @p2 AND deleted_at IS NULL
		GROUP BY CONVERT(VARCHAR, created_at, 23)
		ORDER BY date DESC
	`

	rows, err := u.db.QueryContext(ctx, query, sql.Named("p1", entities.InvoiceStatusPaid), sql.Named("p2", startDate))
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var dp dto.RevenueDataPoint
		if err := rows.Scan(&dp.Date, &dp.Amount, &dp.Count); err != nil {
			continue
		}
		stats.Data = append(stats.Data, dp)
		stats.TotalCount += dp.Count
	}

	return stats, nil
}

func (u *AnalyticsUsecase) GetServiceStats(ctx context.Context) (*dto.ServiceStatsResponse, error) {
	stats := &dto.ServiceStatsResponse{
		StatusBreakdown: make(map[string]int),
		ServicesByType:  []dto.ServiceTypeCount{},
	}

	query := `SELECT COUNT(*) FROM waiting_lists WHERE deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query).Scan(&stats.TotalServices)

	query = `SELECT status, COUNT(*) as count FROM waiting_lists WHERE deleted_at IS NULL GROUP BY status`
	rows, err := u.db.QueryContext(ctx, query)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var status string
			var count int
			if err := rows.Scan(&status, &count); err == nil {
				stats.StatusBreakdown[status] = count
			}
		}
	}

	query = `SELECT AVG(DATEDIFF(HOUR, service_start_at, service_end_at)) FROM waiting_lists WHERE service_end_at IS NOT NULL AND service_start_at IS NOT NULL AND deleted_at IS NULL`
	var avgCompletion sql.NullFloat64
	_ = u.db.QueryRowContext(ctx, query).Scan(&avgCompletion)
	if avgCompletion.Valid {
		stats.AverageCompletion = avgCompletion.Float64
	}

	return stats, nil
}

func (u *AnalyticsUsecase) GetQueueStats(ctx context.Context) (*dto.QueueStatsResponse, error) {
	stats := &dto.QueueStatsResponse{}
	today := time.Now().Truncate(24 * time.Hour)

	query := `SELECT COUNT(*) FROM waiting_lists WHERE CAST(service_date AS DATE) = CAST(@p1 AS DATE) AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", today)).Scan(&stats.TodayTotal)

	query = `SELECT COUNT(*) FROM waiting_lists WHERE status = 'waiting' AND CAST(service_date AS DATE) = CAST(@p1 AS DATE) AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", today)).Scan(&stats.CurrentWaiting)

	query = `SELECT COUNT(*) FROM waiting_lists WHERE status = 'in_service' AND CAST(service_date AS DATE) = CAST(@p1 AS DATE) AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", today)).Scan(&stats.CurrentInService)

	query = `SELECT COUNT(*) FROM waiting_lists WHERE status = 'completed' AND CAST(service_date AS DATE) = CAST(@p1 AS DATE) AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", today)).Scan(&stats.Completed)

	query = `SELECT COUNT(*) FROM waiting_lists WHERE status = 'no_show' AND CAST(service_date AS DATE) = CAST(@p1 AS DATE) AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", today)).Scan(&stats.NoShow)

	query = `SELECT COUNT(*) FROM waiting_lists WHERE status = 'canceled' AND CAST(service_date AS DATE) = CAST(@p1 AS DATE) AND deleted_at IS NULL`
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", today)).Scan(&stats.Cancelled)

	query = `SELECT AVG(DATEDIFF(MINUTE, created_at, called_at)) FROM waiting_lists WHERE called_at IS NOT NULL AND CAST(service_date AS DATE) = CAST(@p1 AS DATE) AND deleted_at IS NULL`
	var avgWait sql.NullFloat64
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", today)).Scan(&avgWait)
	if avgWait.Valid {
		stats.AverageWaitTime = avgWait.Float64
	}

	query = `SELECT TOP 1 DATEPART(HOUR, created_at) as hour FROM waiting_lists WHERE CAST(service_date AS DATE) = CAST(@p1 AS DATE) AND deleted_at IS NULL GROUP BY DATEPART(HOUR, created_at) ORDER BY COUNT(*) DESC`
	var peakHour sql.NullInt64
	_ = u.db.QueryRowContext(ctx, query, sql.Named("p1", today)).Scan(&peakHour)
	if peakHour.Valid {
		stats.PeakHour = int(peakHour.Int64)
	}

	return stats, nil
}

func (u *AnalyticsUsecase) GetMechanicPerformance(ctx context.Context) ([]dto.MechanicPerformanceResponse, error) {
	performances := []dto.MechanicPerformanceResponse{}

	// Note: Current schema doesn't have mechanic_id in waiting_lists
	// This query returns mechanics with zero services until the schema is updated
	query := `
		SELECT u.id, u.name,
		       0 as total_services,
		       0 as completed_services,
		       0.0 as avg_completion
		FROM users u
		WHERE u.role = 'mechanic' AND u.deleted_at IS NULL
		GROUP BY u.id, u.name
	`

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		return performances, err
	}
	defer rows.Close()

	for rows.Next() {
		var perf dto.MechanicPerformanceResponse
		var avgCompletion sql.NullFloat64

		if err := rows.Scan(&perf.MechanicID, &perf.MechanicName, &perf.TotalServices, &perf.CompletedServices, &avgCompletion); err != nil {
			continue
		}

		if avgCompletion.Valid {
			perf.AverageCompletion = avgCompletion.Float64
		}

		if perf.TotalServices > 0 {
			perf.Efficiency = float64(perf.CompletedServices) / float64(perf.TotalServices) * 100
		}

		perf.CustomerRating = 0.0
		performances = append(performances, perf)
	}

	return performances, nil
}
