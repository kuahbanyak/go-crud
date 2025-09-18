package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

type BookingHandler struct {
	bookingUsecase *usecases.BookingUsecase
}

func NewBookingHandler(bookingUsecase *usecases.BookingUsecase) *BookingHandler {
	return &BookingHandler{
		bookingUsecase: bookingUsecase,
	}
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	userID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	booking := &entities.Booking{
		CustomerID:  userID,
		VehicleID:   req.VehicleID,
		ScheduledAt: req.StartDate,
		Notes:       req.Notes,
		Status:      entities.StatusScheduled,
	}

	err := h.bookingUsecase.CreateBooking(r.Context(), booking)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to create booking", err)
		return
	}

	response.Success(w, http.StatusCreated, "Booking created successfully", map[string]interface{}{
		"id":           booking.ID,
		"customer_id":  booking.CustomerID,
		"vehicle_id":   booking.VehicleID,
		"scheduled_at": booking.ScheduledAt,
		"status":       booking.Status,
		"notes":        booking.Notes,
	})
}

func (h *BookingHandler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	userID, ok := r.Context().Value("id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	bookings, err := h.bookingUsecase.GetBookingsByCustomer(r.Context(), userID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get bookings", err)
		return
	}

	response.Success(w, http.StatusOK, "Bookings retrieved successfully", bookings)
}

func (h *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid booking ID", err)
		return
	}

	booking, err := h.bookingUsecase.GetBooking(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Booking not found", err)
		return
	}

	response.Success(w, http.StatusOK, "Booking retrieved successfully", booking)
}

func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid booking ID", err)
		return
	}

	var req dto.UpdateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	userID, ok := r.Context().Value("user_id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	existingBooking, err := h.bookingUsecase.GetBooking(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Booking not found", err)
		return
	}

	if existingBooking.CustomerID != userID {
		response.Error(w, http.StatusForbidden, "Forbidden", "You can only update your own bookings")
		return
	}

	updatedBooking := &entities.Booking{
		ID:          id,
		CustomerID:  userID,
		VehicleID:   req.VehicleID,
		ScheduledAt: req.StartDate,
		Notes:       req.Notes,
	}

	if req.Status != "" {
		switch req.Status {
		case "scheduled":
			updatedBooking.Status = entities.StatusScheduled
		case "in_progress":
			updatedBooking.Status = entities.StatusInProgress
		case "completed":
			updatedBooking.Status = entities.StatusCompleted
		case "cancelled":
			updatedBooking.Status = entities.StatusCanceled
		default:
			updatedBooking.Status = entities.StatusScheduled
		}
	} else {
		updatedBooking.Status = existingBooking.Status
	}

	err = h.bookingUsecase.UpdateBooking(r.Context(), id, updatedBooking)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to update booking", err)
		return
	}

	finalBooking, err := h.bookingUsecase.GetBooking(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve updated booking", err)
		return
	}

	response.Success(w, http.StatusOK, "Booking updated successfully", finalBooking)
}

func (h *BookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid booking ID", err)
		return
	}

	userID, ok := r.Context().Value("user_id").(types.MSSQLUUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	existingBooking, err := h.bookingUsecase.GetBooking(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Booking not found", err)
		return
	}
	if existingBooking.CustomerID != userID {
		response.Error(w, http.StatusForbidden, "Forbidden", "You can only delete your own bookings")
		return
	}

	err = h.bookingUsecase.DeleteBooking(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete booking", err)
		return
	}

	response.Success(w, http.StatusOK, "Booking deleted successfully", nil)
}
