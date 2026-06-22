package api

import (
	"net/http"

	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	serviceinterfaces "github.com/shreygarg/trip-planner-agent/service/interfaces"
	"github.com/shreygarg/trip-planner-agent/utils"
	validationinterfaces "github.com/shreygarg/trip-planner-agent/validations/interfaces"
)

// APIHandler wraps agent services into HTTP handlers.
type APIHandler struct {
	agent     serviceinterfaces.Agent
	validator validationinterfaces.Validator
}

// NewAPIHandler returns an initialized APIHandler.
func NewAPIHandler(agent serviceinterfaces.Agent, validator validationinterfaces.Validator) *APIHandler {
	return &APIHandler{
		agent:     agent,
		validator: validator,
	}
}

// PlanTripHandler handles requests to recommend trip plans.
func (h *APIHandler) PlanTripHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": constants.ErrMethodNotAllowed})
		return
	}

	var req models.PlanRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": constants.ErrInvalidJSON})
		return
	}

	if err := h.validator.ValidatePlanRequest(req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := h.agent.PlanTrip(req.Prompt)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, models.PlanResponse{Response: res})
}
