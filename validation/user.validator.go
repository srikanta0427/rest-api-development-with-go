package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/srikanta0427/rest_api_design/model"
)

//var validate *validator.Validate

func ValidateUserStruct(user model.User) error {
	validate := validator.New()
	return validate.Struct(user)
}

func OrganizerValidation(organizer model.Organizer) error {
	validate := validator.New()
	return validate.Struct(organizer)
}
