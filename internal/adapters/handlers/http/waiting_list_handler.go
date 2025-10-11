package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

type WaitingListHandler struct {
	waitingListUsecase *usecases.WaitingListUsecase
}

func NewWaitingListHandler(waitingListUsecase *usecases.WaitingListUsecase) *WaitingListHandler {
	return &WaitingListHandler{
		waitingListUsecase: waitingListUsecase,
	}
}

// TakeQueueNumber allows a customer to take a queue number for service
func (h *WaitingListHandler) TakeQueueNumber(w http.ResponseWriter, r *http.Request) {
	var req dto.TakeQueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Get customer ID from context (set by auth middleware)
	customerID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Parse service date from string (YYYY-MM-DD format)
	serviceDate, err := time.Parse("2006-01-02", req.ServiceDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid service_date format. Use YYYY-MM-DD", err)
		return
	}

	waitingList := &entities.WaitingList{
		VehicleID:     req.VehicleID,
		CustomerID:    customerID,
		ServiceDate:   serviceDate,
		ServiceType:   req.ServiceType,
		EstimatedTime: req.EstimatedTime,
		Notes:         req.Notes,
	}

	if err := h.waitingListUsecase.TakeQueueNumber(r.Context(), waitingList); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to take queue number", err)
		return
	}

	resp := dto.WaitingListResponse{
		ID:            waitingList.ID,
		QueueNumber:   waitingList.QueueNumber,
		VehicleID:     waitingList.VehicleID,
		CustomerID:    waitingList.CustomerID,
		ServiceDate:   waitingList.ServiceDate,
		ServiceType:   waitingList.ServiceType,
		EstimatedTime: waitingList.EstimatedTime,
		Status:        string(waitingList.Status),
		Notes:         waitingList.Notes,
		CreatedAt:     waitingList.CreatedAt,
		UpdatedAt:     waitingList.UpdatedAt,
	}

	response.Success(w, http.StatusCreated, "Queue number taken successfully", resp)
}

// GetMyQueue retrieves all queue entries for the authenticated customer
func (h *WaitingListHandler) GetMyQueue(w http.ResponseWriter, r *http.Request) {
	customerID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	waitingLists, err := h.waitingListUsecase.GetCustomerWaitingLists(r.Context(), customerID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get queue entries", map[string]interface{}{
			"error":       err.Error(),
			"customer_id": customerID.String(),
		})
		return
	}

	// If no entries found, return empty array with success
	if len(waitingLists) == 0 {
		response.Success(w, http.StatusOK, "No queue entries found", []interface{}{})
		return
	}

	resp := make([]dto.WaitingListWithDetailsResponse, len(waitingLists))
	for i, wl := range waitingLists {
		resp[i] = h.buildDetailResponse(wl)
	}

	response.Success(w, http.StatusOK, "Queue entries retrieved successfully", resp)
}

// GetTodayQueue retrieves today's queue
func (h *WaitingListHandler) GetTodayQueue(w http.ResponseWriter, r *http.Request) {
	waitingLists, err := h.waitingListUsecase.GetTodayQueue(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get today's queue", err)
		return
	}

	resp := h.buildWaitingListResponse(waitingLists, time.Now())
	response.Success(w, http.StatusOK, "Today's queue retrieved successfully", resp)
}

// GetQueueByDate retrieves queue by specific date
func (h *WaitingListHandler) GetQueueByDate(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		response.Error(w, http.StatusBadRequest, "Date parameter is required", nil)
		return
	}

	serviceDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD", err)
		return
	}

	waitingLists, err := h.waitingListUsecase.GetQueueByDate(r.Context(), serviceDate)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get queue", err)
		return
	}

	resp := h.buildWaitingListResponse(waitingLists, serviceDate)
	response.Success(w, http.StatusOK, "Queue retrieved successfully", resp)
}

// GetQueueByNumber retrieves a queue entry by queue number
func (h *WaitingListHandler) GetQueueByNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queueNumber, err := strconv.Atoi(vars["number"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid queue number", err)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	serviceDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD", err)
		return
	}

	waitingList, err := h.waitingListUsecase.GetByQueueNumber(r.Context(), queueNumber, serviceDate)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Queue entry not found", err)
		return
	}

	resp := h.buildDetailResponse(waitingList)
	response.Success(w, http.StatusOK, "Queue entry retrieved successfully", resp)
}

// CallCustomer marks a customer as called for service
func (h *WaitingListHandler) CallCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	if err := h.waitingListUsecase.CallCustomer(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to call customer", err)
		return
	}

	response.Success(w, http.StatusOK, "Customer called successfully", nil)
}

// StartService marks the start of service for a customer
func (h *WaitingListHandler) StartService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	if err := h.waitingListUsecase.StartService(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to start service", err)
		return
	}

	response.Success(w, http.StatusOK, "Service started successfully", nil)
}

// CompleteService marks the completion of service for a customer
func (h *WaitingListHandler) CompleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	if err := h.waitingListUsecase.CompleteService(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to complete service", err)
		return
	}

	response.Success(w, http.StatusOK, "Service completed successfully", nil)
}

// CancelQueue cancels a queue entry
func (h *WaitingListHandler) CancelQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	if err := h.waitingListUsecase.CancelQueue(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to cancel queue", err)
		return
	}

	response.Success(w, http.StatusOK, "Queue cancelled successfully", nil)
}

// MarkNoShow marks a customer as no-show
func (h *WaitingListHandler) MarkNoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	if err := h.waitingListUsecase.MarkNoShow(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to mark no-show", err)
		return
	}

	response.Success(w, http.StatusOK, "Marked as no-show successfully", nil)
}

// CheckAvailability checks ticket availability for a given date
func (h *WaitingListHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	var serviceDate time.Time
	var err error

	if dateStr == "" {
		serviceDate = time.Now()
	} else {
		serviceDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD", err)
			return
		}
	}

	available, remaining, err := h.waitingListUsecase.CheckTicketAvailability(r.Context(), serviceDate)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to check availability", err)
		return
	}

	resp := map[string]interface{}{
		"date":              serviceDate.Format("2006-01-02"),
		"available":         available,
		"remaining_tickets": remaining,
		"max_tickets":       10,
		"message": func() string {
			if available {
				return fmt.Sprintf("%d tickets remaining for this date", remaining)
			}
			return "No tickets available for this date (limit reached)"
		}(),
	}

	response.Success(w, http.StatusOK, "Availability checked successfully", resp)
}

// GetServiceProgress allows customers to track the real-time progress of their car service
func (h *WaitingListHandler) GetServiceProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	// Get customer ID from context
	customerID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Get the waiting list entry
	waitingList, err := h.waitingListUsecase.GetWaitingList(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Service ticket not found", err)
		return
	}

	// Verify that this ticket belongs to the customer
	if waitingList.CustomerID != customerID {
		response.Error(w, http.StatusForbidden, "You don't have permission to view this ticket", nil)
		return
	}

	// Get all tickets for the same service date to calculate queue position
	allTickets, err := h.waitingListUsecase.GetQueueByDate(r.Context(), waitingList.ServiceDate)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve queue information", err)
		return
	}

	// Calculate progress metrics
	var currentlyServing int
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

	// Calculate estimated wait time (30 minutes per service on average)
	estimatedWaitMin := waitingAhead * 30
	if waitingList.Status == entities.WaitingListStatusInService ||
		waitingList.Status == entities.WaitingListStatusCompleted {
		estimatedWaitMin = 0
	}

	// Generate status message
	statusMessage := h.generateStatusMessage(waitingList.Status, waitingAhead, currentlyServing, waitingList.QueueNumber)

	// Build response with detailed progress information
	resp := dto.ServiceProgressResponse{
		ID:            waitingList.ID,
		QueueNumber:   waitingList.QueueNumber,
		Status:        string(waitingList.Status),
		StatusMessage: statusMessage,
		ServiceType:   waitingList.ServiceType,
		ServiceDate:   waitingList.ServiceDate,
		EstimatedTime: waitingList.EstimatedTime,
		QueuePosition: waitingList.QueueNumber,
		PeopleAhead:   waitingAhead,
		EstimatedWait: estimatedWaitMin,
		Timeline: dto.Timeline{
			QueueTakenAt:   waitingList.CreatedAt,
			CalledAt:       waitingList.CalledAt,
			ServiceStartAt: waitingList.ServiceStartAt,
			ServiceEndAt:   waitingList.ServiceEndAt,
		},
		Notes: waitingList.Notes,
	}

	// Add vehicle details if available
	if waitingList.Vehicle.ID.String() != "00000000-0000-0000-0000-000000000000" {
		resp.VehicleBrand = waitingList.Vehicle.Brand
		resp.VehicleModel = waitingList.Vehicle.Model
		resp.LicensePlate = waitingList.Vehicle.LicensePlate
	}

	response.Success(w, http.StatusOK, "Service progress retrieved successfully", resp)
}

// generateStatusMessage creates a user-friendly message based on current status
func (h *WaitingListHandler) generateStatusMessage(status entities.WaitingListStatus, waitingAhead, currentlyServing, queueNumber int) string {
	switch status {
	case entities.WaitingListStatusWaiting:
		if waitingAhead == 0 {
			return "🎉 You're next! Please be ready to bring your vehicle to the service area."
		}
		if currentlyServing > 0 {
			return fmt.Sprintf("⏳ %d customer(s) ahead of you. Currently serving queue #%d", waitingAhead, currentlyServing)
		}
		return fmt.Sprintf("⏳ %d customer(s) ahead of you in the queue", waitingAhead)
	case entities.WaitingListStatusCalled:
		return "📢 You've been called! Please proceed to the service area immediately."
	case entities.WaitingListStatusInService:
		return "🔧 Your vehicle is currently being serviced. Please wait in the customer lounge."
	case entities.WaitingListStatusCompleted:
		return "✅ Your service has been completed! Thank you for choosing our service."
	case entities.WaitingListStatusCanceled:
		return "❌ This service ticket has been canceled."
	case entities.WaitingListStatusNoShow:
		return "⚠️ You were marked as no-show. Please contact us to reschedule your service."
	default:
		return "Status information not available"
	}
}

// Helper functions
func (h *WaitingListHandler) buildWaitingListResponse(waitingLists []*entities.WaitingList, serviceDate time.Time) dto.WaitingListListResponse {
	resp := make([]dto.WaitingListWithDetailsResponse, len(waitingLists))
	for i, wl := range waitingLists {
		resp[i] = h.buildDetailResponse(wl)
	}

	return dto.WaitingListListResponse{
		WaitingLists: resp,
		Total:        len(waitingLists),
		Date:         serviceDate.Format("2006-01-02"),
	}
}

func (h *WaitingListHandler) buildDetailResponse(wl *entities.WaitingList) dto.WaitingListWithDetailsResponse {
	resp := dto.WaitingListWithDetailsResponse{
		ID:             wl.ID,
		QueueNumber:    wl.QueueNumber,
		VehicleID:      wl.VehicleID,
		CustomerID:     wl.CustomerID,
		ServiceDate:    wl.ServiceDate,
		ServiceType:    wl.ServiceType,
		EstimatedTime:  wl.EstimatedTime,
		Status:         string(wl.Status),
		CalledAt:       wl.CalledAt,
		ServiceStartAt: wl.ServiceStartAt,
		ServiceEndAt:   wl.ServiceEndAt,
		Notes:          wl.Notes,
		CreatedAt:      wl.CreatedAt,
		UpdatedAt:      wl.UpdatedAt,
	}

	if wl.Vehicle.ID.String() != "00000000-0000-0000-0000-000000000000" {
		resp.VehicleBrand = wl.Vehicle.Brand
		resp.VehicleModel = wl.Vehicle.Model
		resp.LicensePlate = wl.Vehicle.LicensePlate
	}

	if wl.Customer.ID.String() != "00000000-0000-0000-0000-000000000000" {
		resp.CustomerName = wl.Customer.Name
		resp.CustomerPhone = wl.Customer.Phone
	}

	return resp
}
