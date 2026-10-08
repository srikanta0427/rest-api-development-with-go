package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/srikanta0427/rest_api_design/model"
	"github.com/srikanta0427/rest_api_design/services"
	"github.com/srikanta0427/rest_api_design/util"
	"github.com/srikanta0427/rest_api_design/validation"
)

type OrganizerController struct {
	organizerService services.OrganizerService
}

func NewOrganizerController(
	service services.OrganizerService,
) *OrganizerController {
	return &OrganizerController{
		organizerService: service,
	}
}

// Send OTP and temporarily store organizer details.

func (u *OrganizerController) RegisterOrganizer(
	w http.ResponseWriter,
	r *http.Request,
) {
	var org model.Organizer

	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		util.Error(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	org.ID = 0
	org.Status = "active"
	org.VerificationStatus = "pending"

	if err := validation.OrganizerValidation(org); err != nil {
		util.Error(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	if err := u.organizerService.SendRegistrationOtp(&org); err != nil {
		log.Printf("Sending registration OTP failed: %v", err)
		util.Error(
			w,
			http.StatusInternalServerError,
			"Registration failed",
			"Could not send OTP",
		)
		return
	}

	util.Success(w, http.StatusAccepted, "OTP sent", nil)
}

// Verify OTP, then save organizer to the database.

func (u *OrganizerController) VerifyOrganizer(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.Error(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	if request.Email == "" || request.OTP == "" {
		util.Error(
			w,
			http.StatusBadRequest,
			"Missing fields",
			"Email and OTP are required",
		)
		return
	}

	if err := u.organizerService.VerifyOtp(
		request.Email,
		request.OTP,
	); err != nil {
		util.Error(w, http.StatusBadRequest, "Verification failed", err.Error())
		return
	}

	util.Success(w, http.StatusCreated, "Organizer created", nil)
}
