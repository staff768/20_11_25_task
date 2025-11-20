package main

import (
	"19_11_2026_go/internal/app"
	"19_11_2026_go/internal/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Handler struct {
	service app.LinkServicer
}

func NewHandler(s app.LinkServicer) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CheckLinksHandler(w http.ResponseWriter, r *http.Request) {

	var req models.CheckRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Links) == 0 {
		http.Error(w, "Links list cannot be empty", http.StatusBadRequest)
		return
	}

	task, err := h.service.CheckLinks(req.Links)
	if err != nil {
		log.Printf("Service error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (h *Handler) GenerateReportHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GenerateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		http.Error(w, "IDs list cannot be empty", http.StatusBadRequest)
	}
	pdfBytes, err := h.service.GenerateReport(req.IDs)
	if err != nil {
		if errors.Is(err, app.ErrTasksNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Service error when report generation: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="report.pdf"`)
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}
