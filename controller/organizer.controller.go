package controller

import (
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

	if err := u.organizerService.VerifyOtp(
		request.Email,
		request.OTP,
	); err != nil {
		util.Error(w, http.StatusBadRequest, "Verification failed", err.Error())
		return
	}

	util.Success(w, http.StatusCreated, "Organizer created", nil)
}

func (org *OrganizerController) Login(w http.ResponseWriter, r *http.Request) {
	// get the email and password
	// decode json to model
	var orgi model.Organizer
	if err := json.NewDecoder(r.Body).Decode(&orgi); err != nil {
		util.Error(w, http.StatusBadRequest, "Invalid JSON", err.Error())
	}
	fmt.Println(orgi.Email)
	orgi, err := org.organizerService.Login(orgi.Email, orgi.Password)
	if err != nil {
		fmt.Println("e rroror")
		util.Error(w, http.StatusBadRequest, "Login failed", err.Error())
	}
	token, err := middlewere.GenerateToken(orgi.ID)
	if err != nil {
		fmt.Println(err)
		util.Error(w, http.StatusInternalServerError, "Error generating token", err.Error())
	}
	// set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
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
