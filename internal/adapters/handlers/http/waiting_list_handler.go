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
	vehicleUsecase     *usecases.VehicleUseCase
}

func NewWaitingListHandler(waitingListUsecase *usecases.WaitingListUsecase, vehicleUsecase *usecases.VehicleUseCase) *WaitingListHandler {
	return &WaitingListHandler{
		waitingListUsecase: waitingListUsecase,
		vehicleUsecase:     vehicleUsecase,
	}
}
func (h *WaitingListHandler) TakeQueueNumber(w http.ResponseWriter, r *http.Request) {
	var req dto.TakeQueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	customerID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Check if user has vehicles
	userVehicles, err := h.vehicleUsecase.GetMyVehicles(r.Context(), customerID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get user vehicles", err)
		return
	}

	var vehicleID types.MSSQLUUID

	// Handle vehicle selection or creation
	if req.VehicleID != nil {
		// User selected an existing vehicle
		vehicleID = *req.VehicleID

		// Verify ownership
		_, err := h.vehicleUsecase.GetVehicleByID(r.Context(), customerID, vehicleID)
		if err != nil {
			response.Error(w, http.StatusForbidden, "Vehicle not found or you don't own this vehicle", err)
			return
		}
	} else if req.NewVehicle != nil {
		// User wants to add a new vehicle
		newVehicle, err := h.vehicleUsecase.CreateVehicle(r.Context(), customerID, req.NewVehicle)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Failed to create new vehicle", err)
			return
		}
		vid, _ := types.ParseMSSQLUUID(newVehicle.ID)
		vehicleID = vid
	} else if len(userVehicles) == 0 {
		// No vehicle provided and user has no vehicles
		response.Error(w, http.StatusBadRequest, "No vehicle specified. Please provide vehicle_id or new_vehicle", map[string]interface{}{
			"has_vehicles": false,
			"message":      "You need to add a vehicle first",
		})
		return
	} else {
		// No vehicle selected but user has vehicles
		response.Error(w, http.StatusBadRequest, "Please select a vehicle or add a new one", map[string]interface{}{
			"has_vehicles": true,
			"vehicles":     userVehicles,
		})
		return
	}

	serviceDate, err := time.Parse("2006-01-02", req.ServiceDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid service_date format. Use YYYY-MM-DD", err)
		return
	}

	waitingList := &entities.WaitingList{
		VehicleID:     vehicleID,
		CustomerID:    customerID,
		ServiceDate:   serviceDate,
		ServiceType:   req.ServiceType,
		EstimatedTime: req.EstimatedTime, // Default 0 if not provided, will be updated by mechanic
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
		MechanicNotes: waitingList.MechanicNotes,
		CreatedAt:     waitingList.CreatedAt,
		UpdatedAt:     waitingList.UpdatedAt,
	}

	response.Success(w, http.StatusCreated, "Queue number taken successfully", resp)
}
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
func (h *WaitingListHandler) GetTodayQueue(w http.ResponseWriter, r *http.Request) {
	waitingLists, err := h.waitingListUsecase.GetTodayQueue(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get today's queue", err)
		return
	}
	resp := h.buildWaitingListResponse(waitingLists, time.Now())
	response.Success(w, http.StatusOK, "Today's queue retrieved successfully", resp)
}
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

// UpdateWaitingList allows admin and mechanic to update queue with notes and estimates
func (h *WaitingListHandler) UpdateWaitingList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := types.ParseMSSQLUUID(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	var req dto.UpdateWaitingListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Build update entity from request
	updates := &entities.WaitingList{
		ServiceType:   req.ServiceType,
		EstimatedTime: req.EstimatedTime,
		Notes:         req.Notes,
		MechanicNotes: req.MechanicNotes,
	}

	if req.Status != "" {
		updates.Status = entities.WaitingListStatus(req.Status)
	}

	if err := h.waitingListUsecase.UpdateWaitingList(r.Context(), id, updates); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update waiting list", err)
		return
	}

	// Get updated waiting list to return
	updatedWL, err := h.waitingListUsecase.GetWaitingList(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get updated waiting list", err)
		return
	}

	resp := h.buildDetailResponse(updatedWL)
	response.Success(w, http.StatusOK, "Waiting list updated successfully", resp)
}
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
func (h *WaitingListHandler) GetServiceProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID format", err)
		return
	}
	customerID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	waitingList, err := h.waitingListUsecase.GetWaitingList(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Service ticket not found", err)
		return
	}
	if waitingList.CustomerID != customerID {
		response.Error(w, http.StatusForbidden, "You don't have permission to view this ticket", nil)
		return
	}
	allTickets, err := h.waitingListUsecase.GetQueueByDate(r.Context(), waitingList.ServiceDate)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve queue information", err)
		return
	}
	var currentlyServing int
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
	estimatedWaitMin := waitingAhead * 30
	if waitingList.Status == entities.WaitingListStatusInService ||
		waitingList.Status == entities.WaitingListStatusCompleted {
		estimatedWaitMin = 0
	}
	statusMessage := h.generateStatusMessage(waitingList.Status, waitingAhead, currentlyServing, waitingList.QueueNumber)
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
	if waitingList.Vehicle.ID.String() != "00000000-0000-0000-0000-000000000000" {
		resp.VehicleBrand = waitingList.Vehicle.Brand
		resp.VehicleModel = waitingList.Vehicle.Model
		resp.LicensePlate = waitingList.Vehicle.LicensePlate
	}
	if waitingList.Mechanic != nil && waitingList.Mechanic.ID.String() != "00000000-0000-0000-0000-000000000000" {
		resp.MechanicName = waitingList.Mechanic.Name
	}
	response.Success(w, http.StatusOK, "Service progress retrieved successfully", resp)
}

// GetAllServiceProgress allows admin to view progress of all queues
func (h *WaitingListHandler) GetAllServiceProgress(w http.ResponseWriter, r *http.Request) {
	// Get date from query parameter, default to today
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

	// Get all tickets for the date
	allTickets, err := h.waitingListUsecase.GetQueueByDate(r.Context(), serviceDate)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve queue information", err)
		return
	}

	if len(allTickets) == 0 {
		response.Success(w, http.StatusOK, "No queue entries found for this date", []interface{}{})
		return
	}

	// Calculate current serving queue
	var currentlyServing int
	for _, ticket := range allTickets {
		if ticket.Status == entities.WaitingListStatusInService {
			currentlyServing = ticket.QueueNumber
			break
		}
	}

	// Build progress response for each ticket
	progressList := make([]dto.ServiceProgressResponse, 0, len(allTickets))

	for _, waitingList := range allTickets {
		// Calculate waiting ahead for this ticket
		var waitingAhead int
		for _, ticket := range allTickets {
			if ticket.QueueNumber < waitingList.QueueNumber &&
				(ticket.Status == entities.WaitingListStatusWaiting || ticket.Status == entities.WaitingListStatusCalled) {
				waitingAhead++
			}
		}

		estimatedWaitMin := waitingAhead * 30
		if waitingList.Status == entities.WaitingListStatusInService ||
			waitingList.Status == entities.WaitingListStatusCompleted {
			estimatedWaitMin = 0
		}

		statusMessage := h.generateStatusMessage(waitingList.Status, waitingAhead, currentlyServing, waitingList.QueueNumber)

		progress := dto.ServiceProgressResponse{
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

		// Add vehicle info if available
		if waitingList.Vehicle.ID.String() != "00000000-0000-0000-0000-000000000000" {
			progress.VehicleBrand = waitingList.Vehicle.Brand
			progress.VehicleModel = waitingList.Vehicle.Model
			progress.LicensePlate = waitingList.Vehicle.LicensePlate
		}

		// Add customer info
		if waitingList.Customer.ID.String() != "00000000-0000-0000-0000-000000000000" {
			progress.CustomerName = waitingList.Customer.Name
			progress.CustomerPhone = waitingList.Customer.Phone
		}

		// Add mechanic info
		if waitingList.Mechanic != nil && waitingList.Mechanic.ID.String() != "00000000-0000-0000-0000-000000000000" {
			progress.MechanicName = waitingList.Mechanic.Name
		}

		progressList = append(progressList, progress)
	}

	responseData := map[string]interface{}{
		"date":              serviceDate.Format("2006-01-02"),
		"total_queues":      len(progressList),
		"currently_serving": currentlyServing,
		"progress_list":     progressList,
	}

	response.Success(w, http.StatusOK, "All service progress retrieved successfully", responseData)
}

// GetAvailableQueues returns list of queues available for service (waiting or called status)
func (h *WaitingListHandler) GetAvailableQueues(w http.ResponseWriter, r *http.Request) {
	// Get date from query parameter, default to today
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

	// Get available queues
	availableQueues, err := h.waitingListUsecase.GetAvailableQueues(r.Context(), serviceDate)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve available queues", err)
		return
	}

	if len(availableQueues) == 0 {
		response.Success(w, http.StatusOK, "No queues available for service", []interface{}{})
		return
	}

	// Build response with full details
	queueList := make([]dto.WaitingListWithDetailsResponse, 0, len(availableQueues))
	for _, queue := range availableQueues {
		queueList = append(queueList, h.buildDetailResponse(queue))
	}

	responseData := map[string]interface{}{
		"date":             serviceDate.Format("2006-01-02"),
		"total_queues":     len(queueList),
		"available_queues": queueList,
	}

	response.Success(w, http.StatusOK, "Available queues retrieved successfully", responseData)
}

// AssignMechanicToQueue allows mechanic to assign themselves to service a queue
func (h *WaitingListHandler) AssignMechanicToQueue(w http.ResponseWriter, r *http.Request) {
	mechanicID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req dto.AssignMechanicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Assign mechanic to queue
	if err := h.waitingListUsecase.AssignMechanicToQueue(r.Context(), req.QueueID, mechanicID); err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to assign mechanic to queue", err)
		return
	}

	// Get updated queue to return
	updatedQueue, err := h.waitingListUsecase.GetWaitingList(r.Context(), req.QueueID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get updated queue", err)
		return
	}

	resp := h.buildDetailResponse(updatedQueue)
	response.Success(w, http.StatusOK, "Mechanic assigned to queue successfully", resp)
}

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
		MechanicID:     wl.MechanicID,
		ServiceDate:    wl.ServiceDate,
		ServiceType:    wl.ServiceType,
		EstimatedTime:  wl.EstimatedTime,
		Status:         string(wl.Status),
		CalledAt:       wl.CalledAt,
		ServiceStartAt: wl.ServiceStartAt,
		ServiceEndAt:   wl.ServiceEndAt,
		Notes:          wl.Notes,
		MechanicNotes:  wl.MechanicNotes,
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
	if wl.Mechanic != nil && wl.Mechanic.ID.String() != "00000000-0000-0000-0000-000000000000" {
		resp.MechanicName = wl.Mechanic.Name
	}
	return resp
}
