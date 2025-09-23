package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"paperwork-service/internal/services"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// PaperworkHandler handles PDF generation requests
type PaperworkHandler struct {
	logger       *zap.Logger
	eventService *services.EventService
	pdfService   *services.PaperworkPDFService
}

// NewPaperworkHandler creates a new paperwork handler
func NewPaperworkHandler(
	logger *zap.Logger,
	eventService *services.EventService,
	pdfService *services.PaperworkPDFService,
) *PaperworkHandler {
	return &PaperworkHandler{
		logger:       logger,
		eventService: eventService,
		pdfService:   pdfService,
	}
}

// GenerateEventPaperwork generates a PDF for the given event EID
func (h *PaperworkHandler) GenerateEventPaperwork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract EID from URL path
	vars := mux.Vars(r)
	eid, exists := vars["eid"]
	if !exists || eid == "" {
		h.respondWithError(w, http.StatusBadRequest, "Event EID is required")
		return
	}

	h.logger.Info("Generating paperwork for event", zap.String("eid", eid))

	// Fetch all required data from the edge function
	data, err := h.eventService.GetEventPaperworkData(ctx, eid)
	if err != nil {
		h.logger.Error("Failed to fetch event data",
			zap.String("eid", eid),
			zap.Error(err))

		if err.Error() == fmt.Sprintf("event not found: %s", eid) {
			h.respondWithError(w, http.StatusNotFound, "Event not found")
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch event data")
		}
		return
	}

	// Check if we have any artists
	if len(data.Artists) == 0 {
		h.logger.Warn("No artists found for event", zap.String("eid", eid))
		h.respondWithError(w, http.StatusNotFound, "No artists found for this event")
		return
	}

	// Generate the PDF
	pdfData, err := h.pdfService.GenerateEventPaperwork(&data.Event, data.Artists, data.AuctionLots)
	if err != nil {
		h.logger.Error("Failed to generate PDF",
			zap.String("eid", eid),
			zap.Error(err))
		h.respondWithError(w, http.StatusInternalServerError, "Failed to generate PDF")
		return
	}

	// Set response headers for PDF download
	filename := fmt.Sprintf("artbattle_%s_paperwork.pdf", eid)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfData)))

	// Write PDF data to response
	if _, err := w.Write(pdfData); err != nil {
		h.logger.Error("Failed to write PDF response",
			zap.String("eid", eid),
			zap.Error(err))
		return
	}

	h.logger.Info("Successfully generated paperwork PDF",
		zap.String("eid", eid),
		zap.String("event_name", data.Event.Name),
		zap.Int("pdf_size_bytes", len(pdfData)),
		zap.Int("artist_count", len(data.Artists)),
		zap.Int("auction_lots", len(data.AuctionLots)))
}

// HealthCheck provides a health check endpoint
func (h *PaperworkHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "ok",
		"service": "paperwork-service",
		"version": "1.0.0",
		"time":    r.Context().Value("request_time"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// respondWithError sends an error response
func (h *PaperworkHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.logger.Error("HTTP error response",
		zap.Int("status_code", code),
		zap.String("message", message))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}