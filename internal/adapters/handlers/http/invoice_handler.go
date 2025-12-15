package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

type InvoiceHandler struct {
	usecase *usecases.InvoiceUsecase
}

func NewInvoiceHandler(usecase *usecases.InvoiceUsecase) *InvoiceHandler {
	return &InvoiceHandler{usecase: usecase}
}

func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	invoice, err := h.usecase.CreateInvoice(r.Context(), &req)
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, "Failed to create invoice", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusCreated, "Invoice created successfully", invoice)
}

func (h *InvoiceHandler) GetInvoices(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	status := r.URL.Query().Get("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	invoices, err := h.usecase.ListInvoices(r.Context(), page, pageSize, status)
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, "Failed to retrieve invoices", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Invoices retrieved successfully", invoices)
}

func (h *InvoiceHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, "Invalid invoice ID", err.Error())
		return
	}

	invoice, err := h.usecase.GetInvoice(r.Context(), id)
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusNotFound, "Invoice not found", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Invoice retrieved successfully", invoice)
}

func (h *InvoiceHandler) UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, "Invalid invoice ID", err.Error())
		return
	}

	var req dto.UpdateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	invoice, err := h.usecase.UpdateInvoice(r.Context(), id, &req)
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, "Failed to update invoice", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Invoice updated successfully", invoice)
}

func (h *InvoiceHandler) DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, "Invalid invoice ID", err.Error())
		return
	}

	if err := h.usecase.DeleteInvoice(r.Context(), id); err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, "Failed to delete invoice", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Invoice deleted successfully", nil)
}

func (h *InvoiceHandler) GetInvoiceByWaitingList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, "Invalid waiting list ID", err.Error())
		return
	}

	invoices, err := h.usecase.GetInvoicesByWaitingList(r.Context(), id)
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusInternalServerError, "Failed to retrieve invoices", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Invoices retrieved successfully", invoices)
}

func (h *InvoiceHandler) PayInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, "Invalid invoice ID", err.Error())
		return
	}

	var req dto.PayInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	invoice, err := h.usecase.PayInvoice(r.Context(), id, &req)
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, "Failed to pay invoice", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Invoice paid successfully", invoice)
}

func (h *InvoiceHandler) DownloadInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusBadRequest, "Invalid invoice ID", err.Error())
		return
	}

	invoice, err := h.usecase.GetInvoice(r.Context(), id)
	if err != nil {
		response.ErrorWithContext(r.Context(), w, http.StatusNotFound, "Invoice not found", err.Error())
		return
	}

	response.SuccessWithContext(r.Context(), w, http.StatusOK, "Invoice PDF download (pending implementation)", invoice)
}
