package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	db "github.com/srikanta0427/rest_api_design/db/repos"
	"github.com/srikanta0427/rest_api_design/model"
	"github.com/srikanta0427/rest_api_design/util"
	security "github.com/srikanta0427/rest_api_design/util/security"
)

type OrganizerService interface {
	SendRegistrationOtp(org *model.Organizer) error
	VerifyOtp(email, otp string, ctx context.Context) error
	SaveOrganizer(email string, accept bool) error
	Login(email, password string) (model.Organizer, error)
	UploadLogo(id int64) (model.Organizer, error)
	UpdateLogoUrl(id int64, url string) error
}

type OrganizerServiceImpl struct {
	organizerRepository db.OrganizerRepo
	emailService        EmailService
	rdb                 *redis.Client
}

func NewOrganizerService(
	repo db.OrganizerRepo,
	emailService EmailService,
	rdb *redis.Client,
) OrganizerService {
	return &OrganizerServiceImpl{
		organizerRepository: repo,
		emailService:        emailService,
		rdb:                 rdb,
	}
}

// STEP 1: Send OTP and store organizer in Redis.
func (s *OrganizerServiceImpl) SendRegistrationOtp(
	org *model.Organizer,
) error {
	if org == nil {
		return errors.New("organizer cannot be nil")
	}

	otp, err := util.RandomNumber6()
	if err != nil {
		return err
	}

	organizer := *org

	// Hash password before storing in Redis.
	hash, err := security.HashPassWord(organizer.Password)
	if err != nil {
		return err
	}

	organizer.Password = hash

	otpString := fmt.Sprintf("%06d", otp)

	modelOrg := model.PendingOrganizer{
		Organizer: &organizer,
		OTP:       otpString,
		Verified:  false,
	}

	data, err := json.Marshal(modelOrg)
	if err != nil {
		return err
	}

	ctx := context.Background()
	key := "organizer:pending:" + organizer.Email

	// Store data for 5 minutes.
	err = s.rdb.Set(
		ctx,
		key,
		data,
		5*time.Minute,
	).Err()

	if err != nil {
		return fmt.Errorf("failed to store in Redis: %w", err)
	}

	body := fmt.Sprintf(
		"<p>Your OTP is <strong>%s</strong></p>",
		otpString,
	)

	if err := s.emailService.SendOtp(
		organizer.Email,
		"OTP Verification",
		body,
	); err != nil {
		return err
	}

	return nil
}

// STEP 2: Verify OTP and update Redis.
func (s *OrganizerServiceImpl) VerifyOtp(
	email, otp string,
	ctx context.Context,
) error {
	key := "organizer:pending:" + email

	// Retrieve stored registration.
	data, err := s.rdb.Get(ctx, key).Bytes()

	if errors.Is(err, redis.Nil) {
		return errors.New("OTP expired or registration not found")
	}
	if err != nil {
		return err
	}

	var modelOrg model.PendingOrganizer

	if err := json.Unmarshal(data, &modelOrg); err != nil {
		return err
	}

	if modelOrg.OTP != otp {
		return errors.New("incorrect OTP")
	}

	// Mark organizer as verified.
	modelOrg.Verified = true

	jsonData, err := json.Marshal(modelOrg)
	if err != nil {
		return err
	}

	// Update Redis without resetting expiration.
	return s.rdb.SetArgs(
		ctx,
		key,
		jsonData,
		redis.SetArgs{
			KeepTTL: true,
		},
	).Err()
}

// STEP 3: Accept terms and save organizer to MySQL.
func (s *OrganizerServiceImpl) SaveOrganizer(
	email string,
	accept bool,
) error {
	if !accept {
		return errors.New("terms and conditions not accepted")
	}

	ctx := context.Background()
	key := "organizer:pending:" + email

	// Retrieve organizer from Redis.
	data, err := s.rdb.Get(ctx, key).Bytes()

	if errors.Is(err, redis.Nil) {
		return errors.New("pending registration expired or not found")
	}
	if err != nil {
		return err
	}

	var pending model.PendingOrganizer

	if err := json.Unmarshal(data, &pending); err != nil {
		return err
	}

	// Check whether OTP was verified.
	if !pending.Verified {
		return errors.New("organizer OTP not verified")
	}

	if pending.Organizer == nil {
		return errors.New("organizer data not found")
	}

	// Insert organizer into MySQL.
	_, err = s.organizerRepository.Create(pending.Organizer)

	if err != nil {
		return fmt.Errorf("could not create organizer: %w", err)
	}

	// Delete temporary registration after successful insert.
	if err := s.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf(
			"organizer saved but Redis cleanup failed: %w",
			err,
		)
	}

	return nil
}

// Login
func (s *OrganizerServiceImpl) Login(
	email, password string,
) (model.Organizer, error) {
	organizer, err := s.organizerRepository.GetByEmail(email)

	if err != nil {
		return model.Organizer{}, errors.New(
			"invalid email or password",
		)
	}

	if !security.ComparePassWord(
		password,
		organizer.Password,
	) {
		return model.Organizer{}, errors.New(
			"invalid email or password",
		)
	}

	return organizer, nil
}

// Upload logo
func (s *OrganizerServiceImpl) UploadLogo(
	id int64,
) (model.Organizer, error) {
	return s.organizerRepository.GetById(id)
}

// Update logo URL
func (s *OrganizerServiceImpl) UpdateLogoUrl(
	id int64,
	url string,
) error {
	return s.organizerRepository.UpdateLogoUrl(id, url)
}
