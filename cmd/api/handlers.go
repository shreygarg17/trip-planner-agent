package api

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	serviceinterfaces "github.com/shreygarg/trip-planner-agent/service/interfaces"
	"github.com/shreygarg/trip-planner-agent/utils"
	validationinterfaces "github.com/shreygarg/trip-planner-agent/validations/interfaces"
)

var (
	handlerOnce     sync.Once
	handlerInstance *APIHandler
)

// APIHandler wraps agent services into HTTP handlers.
type APIHandler struct {
	agent     serviceinterfaces.Agent
	validator validationinterfaces.Validator
}

// NewAPIHandler returns an initialized APIHandler (singleton).
func NewAPIHandler(agent serviceinterfaces.Agent, validator validationinterfaces.Validator) *APIHandler {
	handlerOnce.Do(func() {
		handlerInstance = &APIHandler{
			agent:     agent,
			validator: validator,
		}
	})
	return handlerInstance
}

// PlanTripHandler handles requests to recommend trip plans.
func (h *APIHandler) PlanTripHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	reqID := utils.GenerateRequestID()
	ctx := utils.WithRequestID(r.Context(), reqID)
	ctx, cancel := context.WithTimeout(ctx, constants.APITimeout)
	defer cancel()

	if r.Method != http.MethodPost {
		log.Printf("[%s] [WARN] Invalid method: %s %s - Responded 405", reqID, r.Method, r.URL.Path)
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": constants.ErrMethodNotAllowed})
		return
	}

	var req models.PlanRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		log.Printf("[%s] [WARN] Failed to decode request body: %v - Responded 400", reqID, err)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": constants.ErrInvalidJSON})
		return
	}

	log.Printf("[%s] [INFO] POST /api/v1/trips/plan - Prompt: %q", reqID, req.Prompt)
	if err := h.validator.ValidatePlanRequest(ctx, req); err != nil {
		log.Printf("[%s] [WARN] Request validation failed: %v - Responded 400", reqID, err)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := h.agent.PlanTrip(ctx, req.Prompt)
	if err != nil {
		log.Printf("[%s] [ERROR] Agent failed to plan trip: %v - Responded 500", reqID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("[%s] [INFO] POST /api/v1/trips/plan - Responded 200 OK - Duration: %v", reqID, time.Since(startTime))
	utils.WriteJSON(w, http.StatusOK, models.PlanResponse{Response: res})
}

// HealthHandler responds with a simple health status code 200 OK.
// This is used for backend status verification and frontend integration checks.
func (h *APIHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

