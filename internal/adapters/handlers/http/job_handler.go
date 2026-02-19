package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kuahbanyak/go-crud/internal/domain/entities"
	"github.com/kuahbanyak/go-crud/internal/shared/dto"
	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/internal/usecases"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

type JobHandler struct {
	jobUsecase *usecases.JobUsecase
}

func NewJobHandler(jobUsecase *usecases.JobUsecase) *JobHandler {
	return &JobHandler{
		jobUsecase: jobUsecase,
	}
}

func (h *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.jobUsecase.GetAllJobs(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get jobs", err)
		return
	}

	resp := make([]dto.JobResponse, len(jobs))
	for i, j := range jobs {
		resp[i] = h.buildResponse(j)
	}

	response.Success(w, http.StatusOK, "Jobs retrieved successfully", dto.JobsListResponse{
		Jobs:  resp,
		Total: len(resp),
	})
}

func (h *JobHandler) GetActiveJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.jobUsecase.GetActiveJobs(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get active jobs", err)
		return
	}

	resp := make([]dto.JobResponse, len(jobs))
	for i, j := range jobs {
		resp[i] = h.buildResponse(j)
	}

	response.Success(w, http.StatusOK, "Active jobs retrieved successfully", dto.JobsListResponse{
		Jobs:  resp,
		Total: len(resp),
	})
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	job, err := h.jobUsecase.GetJob(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get job", err)
		return
	}

	if job == nil {
		response.Error(w, http.StatusNotFound, "Job not found", nil)
		return
	}

	resp := h.buildResponse(job)
	response.Success(w, http.StatusOK, "Job retrieved successfully", resp)
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	job := &entities.Job{
		CronJob: req.CronJob,
		Task:    req.Task,
		Status:  req.Status,
	}

	if err := h.jobUsecase.CreateJob(r.Context(), job); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create job", err)
		return
	}

	resp := h.buildResponse(job)
	response.Success(w, http.StatusCreated, "Job created successfully", resp)
}

func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	var req dto.UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.jobUsecase.UpdateJob(r.Context(), id, req.CronJob, req.Task, req.Status); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update job", err)
		return
	}

	response.Success(w, http.StatusOK, "Job updated successfully", nil)
}

func (h *JobHandler) UpdateJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	var req dto.UpdateJobStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.jobUsecase.UpdateJobStatus(r.Context(), id, req.Status); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update job status", err)
		return
	}

	message := "Job disabled successfully"
	if req.Status {
		message = "Job enabled successfully"
	}

	response.Success(w, http.StatusOK, message, nil)
}

func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := types.ParseMSSQLUUID(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	if err := h.jobUsecase.DeleteJob(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete job", err)
		return
	}

	response.Success(w, http.StatusOK, "Job deleted successfully", nil)
}

func (h *JobHandler) buildResponse(j *entities.Job) dto.JobResponse {
	return dto.JobResponse{
		ID:        j.ID,
		CronJob:   j.CronJob,
		Task:      j.Task,
		Status:    j.Status,
		CreatedAt: j.CreatedAt,
		UpdatedAt: j.UpdatedAt,
	}
}
