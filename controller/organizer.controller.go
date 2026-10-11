package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/srikanta0427/rest_api_design/middlewre"
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

	ctx := context.Background()
	if err := u.organizerService.VerifyOtp(
		request.Email,
		request.OTP,
		ctx,

	); err != nil {
		util.Error(w, http.StatusBadRequest, "Verification failed", err.Error())
		return
	}
	//util.Success(w, http.StatusCreated, "Organizer created", nil)
}

func (org *OrganizerController) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.Error(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	// Validate input
	if request.Email == "" || request.Password == "" {
		util.Error(w, http.StatusBadRequest,
			"Email and password required", "")
		return
	}

	// Authenticate
	organizer, err := org.organizerService.Login(
		request.Email,
		request.Password,
	)

	if err != nil {
		fmt.Println("Login failed: %v", err)
		util.Error(w, http.StatusUnauthorized,
			"Invalid email or password", map[string]string{"message": "Invalid email or password"})
		return
	}

	// Generate JWT
	token, err := middlewere.GenerateToken(organizer.ID)
	if err != nil {
		util.Error(w, http.StatusInternalServerError,
			"Token generation failed", "")
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	util.Success(w, http.StatusOK, "Login Success", nil)
}

func (or *OrganizerController) UploadLogo(
	w http.ResponseWriter,
	r *http.Request,
) {
	// 1. Get organizer ID from JWT context
	id, ok := r.Context().Value("organizerID").(int64)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "Unauthorized", "")
		return
	}

	// 2. Get organizer from service
	org, err := or.organizerService.UploadLogo(id)
	if err != nil {
		util.Error(w, http.StatusNotFound, "Organizer not found", "")
		return
	}

	// 3. Limit upload size to 5 MiB
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)

	// 4. Get file from form-data
	file, _, err := r.FormFile("logo")
	if err != nil {
		util.Error(w, http.StatusBadRequest, "Invalid logo upload", "")
		return
	}
	defer file.Close()

	// 5. Create destination directory
	const dir = "upload"

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		util.Error(w, 500, "Cannot create upload directory", "")
		return
	}

	// 6. Create filename using organizer ID
	fileName := fmt.Sprintf("organizer_%d_logo.png", org.ID)
	destFilepath := filepath.Join(dir, fileName)

	// 7. Create destination file
	destFile, err := os.Create(destFilepath)
	if err != nil {
		util.Error(w, 500, "Cannot create destination file", "")
		return
	}

	// 8. Create buffer
	buffer := make([]byte, 32*1024)

	// 9. Stream uploaded file to destination
	_, err = io.CopyBuffer(destFile, file, buffer)
	if err != nil {
		destFile.Close()
		util.Error(w, 500, "Error uploading logo", "")
		return
	}

	// 10. Close destination and check error
	if err := destFile.Close(); err != nil {
		util.Error(w, 500, "Error saving logo", "")
		return
	}

	// 11. Success response
	err = or.organizerService.UpdateLogoUrl(org.ID, fileName)
	if err != nil {
		util.Error(w, 500, "Error uploading logo", "")
	}
	util.Success(w, http.StatusOK, "Upload Success", nil)
}

func (org *OrganizerController) SaveOrganizer(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request struct {
		Email    string `json:"email"`
		Accepted bool   `json:"accepted"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.Error(w, http.StatusBadRequest,
			"Invalid JSON", err.Error())
		return
	}

	if request.Email == "" {
		util.Error(w, http.StatusBadRequest,
			"Email is required", "")
		return
	}

	if !request.Accepted {
		util.Error(w, http.StatusBadRequest,
			"Terms and conditions must be accepted", "")
		return
	}

	err := org.organizerService.SaveOrganizer(
		request.Email,
		request.Accepted,
	)

	if err != nil {
		log.Printf("SaveOrganizer failed: %v", err)

		util.Error(w, http.StatusBadRequest,
			"Cannot save organizer", err.Error())
		return
	}

	util.Success(w, http.StatusCreated,
		"Organizer saved successfully", nil)
}
