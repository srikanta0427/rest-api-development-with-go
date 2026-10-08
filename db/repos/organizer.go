package db

import (
	"database/sql"
	"fmt"

	"github.com/srikanta0427/rest_api_design/model"
)

type OrganizerRepo interface {
	Create(organizer *model.Organizer) (int64, error)
	GetByEmail(email string) (model.Organizer, error)
}

type OrganizerImpl struct {
	db *sql.DB
}

func NewOrganizer(db *sql.DB) OrganizerRepo {
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	return &OrganizerImpl{
		db: db,
	}
}

func (o *OrganizerImpl) Create(org *model.Organizer) (int64, error) {
	const query = ` INSERT INTO organizers (
        name,
        email,
        password_hash,
        organization_name,
        description,
        logo_url,
        contact_phone,
        website_url,
        verification_status,
        status
    )
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := o.db.Exec(query,
		org.Name,
		org.Email,
		org.Password,
		org.OrganizationName,
		org.Description,
		org.LogoURL,
		org.ContactPhone,
		org.WebsiteURL,
		org.VerificationStatus,
		org.Status)

	if err != nil {
		return 0, fmt.Errorf("error creating organizer: %v", err)
	}
	id, err := result.LastInsertId()
	return id, nil

}

func (o *OrganizerImpl) GetByEmail(email string) (model.Organizer, error) {
	const query = "SELECT id, name, email, password_hash, organization_name FROM organizers WHERE email = ?"

	var org model.Organizer

	err := o.db.QueryRow(query, email).Scan(
		&org.ID,
		&org.Name,
		&org.Email,
		&org.Password,
		&org.OrganizationName,
	)
	if err != nil {
		return model.Organizer{}, fmt.Errorf("error getting organizer: %v", err)
	}
	return org, nil
}
