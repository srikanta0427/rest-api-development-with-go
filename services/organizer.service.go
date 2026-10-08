package services

import (
	"fmt"
	"sync"

	db "github.com/srikanta0427/rest_api_design/db/repos"
	"github.com/srikanta0427/rest_api_design/model"
	"github.com/srikanta0427/rest_api_design/util"
	security "github.com/srikanta0427/rest_api_design/util/security"
)

type OrganizerService interface {
	SendRegistrationOtp(org *model.Organizer) error
	VerifyOtp(email, otp string) error
	Login(email, password string) (model.Organizer, error)
}



type OrganizerServiceImpl struct {
	organizerRepository db.OrganizerRepo
	emailService        EmailService

	mu      sync.Mutex
	pending map[string]model.PendingOrganizer
}

func NewOrganizerService(
	repo db.OrganizerRepo,
	emailService EmailService,

) OrganizerService {
	return &OrganizerServiceImpl{
		organizerRepository: repo,
		emailService:        emailService,
		pending:             make(map[string]model.PendingOrganizer),
	}
}

func (s *OrganizerServiceImpl) SendRegistrationOtp(
	org *model.Organizer,
) error {
	otp, err := util.RandomNumber6()
	if err != nil {
		return err
	}

	organizer := *org

	// Hash before storing temporarily.
	hash, err := security.HashPassWord(organizer.Password)
	if err != nil {
		return err
	}
	organizer.Password = hash

	otpString := fmt.Sprintf("%d", otp)
	body := fmt.Sprintf(
		"<p>Your OTP is <strong>%s</strong></p>",
		otpString,
	)

	// Lock so verification cannot run before storage finishes.
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.emailService.SendOtp(
		organizer.Email,
		"OTP Verification",
		body,
	); err != nil {
		return err
	}

	s.pending[organizer.Email] = model.PendingOrganizer{
		Organizer: &organizer,
		OTP:       otpString,
	}

	return nil
}

func (s *OrganizerServiceImpl) VerifyOtp(
	email, otp string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pending, exists := s.pending[email]
	if !exists {
		return fmt.Errorf("no pending registration")
	}

	if pending.OTP != otp {
		return fmt.Errorf("incorrect OTP")
	}

	// Password is already hashed.
	_, err := s.organizerRepository.Create(pending.Organizer)
	if err != nil {
		return fmt.Errorf("could not create organizer")
	}

	delete(s.pending, email)
	return nil
}


func (s *OrganizerServiceImpl) Login(email, password string) (model.Organizer, error) {
	organizer , err := s.organizerRepository.GetByEmail(email)
	if err != nil {
		return model.Organizer{}, fmt.Errorf("could not get organizer")
	}

	// check password
	if password != organizer.Password {
		return model.Organizer{}, fmt.Errorf("incorrect password")
	}

	return organizer, nil
}
