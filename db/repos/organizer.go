package db

import (
	"database/sql"
	"fmt"

	"github.com/srikanta0427/rest_api_design/model"
)

type OrganizerRepo interface {
	Create(organizer *model.Organizer) (int64, error)
	GetByEmail(email string) (model.Organizer, error)
	GetById(id int64) (model.Organizer, error)
	UpdateLogoUrl(id int64, url string) error
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
	fmt.Println("organizer", org)
	if err != nil {
		return model.Organizer{}, fmt.Errorf("error getting organizer: %v", err)
	}
	return org, nil
}

func (o *OrganizerImpl) GetById(id int64) (model.Organizer, error) {
	const query = "SELECT id, name, email, password_hash, organization_name FROM organizers WHERE id = ?"

	var org model.Organizer

	err := o.db.QueryRow(query, id).Scan(
		&org.ID,
		&org.Name,
		&org.Email,
		&org.Password,
		&org.OrganizationName,
	)
	fmt.Println("organizer", org)
	if err != nil {
		return model.Organizer{}, fmt.Errorf("error getting organizer: %v", err)
	}
	return org, nil
}

func (o *OrganizerImpl) UpdateLogoUrl(id int64, url string) error {
	const query = "UPDATE organizers SET logo_url = ? WHERE id = ?"
	_, err := o.db.Exec(query, url, id)
	if err != nil {
		return fmt.Errorf("error updating logo url: %v", err)
	}
	return nil
}
