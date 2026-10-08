package model

import "time"

type Organizer struct {
	ID    int64  `json:"id"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	//Password           string    `json:"password" validate:"required,min=12,max=40,"`
	Password           string    `json:"password" validate:"required,min=12,max=40"`
	OrganizationName   string    `json:"organization_name" validate:"required,max=150"`
	Description        *string   `json:"description,omitempty"`
	LogoURL            *string   `json:"logo_url,omitempty"`
	ContactPhone       *string   `json:"contact_phone,omitempty"`
	WebsiteURL         *string   `json:"website_url,omitempty" validate:"omitempty,url"`
	VerificationStatus string    `json:"verification_status" validate:"oneof=pending verified rejected"`
	Status             string    `json:"status" validate:"oneof=active suspended"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (org *Organizer) NewOrganizer(
	id int64,
	name string,
	email string,
	password string,
	organizationName string,
	description *string,
	logoURL *string,
	contactPhone *string,
	websiteURL *string,
	verificationStatus string,
	status string,
	createdAt time.Time,
	updatedAt time.Time,
) Organizer {
	return Organizer{
		ID:                 id,
		Name:               name,
		Email:              email,
		Password:           password,
		OrganizationName:   organizationName,
		Description:        description,
		LogoURL:            logoURL,
		ContactPhone:       contactPhone,
		WebsiteURL:         websiteURL,
		VerificationStatus: verificationStatus,
		Status:             status,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}
}
