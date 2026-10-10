package router

import (
	"github.com/go-chi/chi"
	"github.com/srikanta0427/rest_api_design/controller"
	middlewere "github.com/srikanta0427/rest_api_design/middlewre"
)

type OrganizerRouter struct {
	organizer *controller.OrganizerController
}

func NewOrganizerRouter(organizer *controller.OrganizerController) *OrganizerRouter {
	return &OrganizerRouter{
		organizer: organizer,
	}
}

func (u *OrganizerRouter) Register(r chi.Router) {
	r.Post("/register-pending", u.organizer.RegisterOrganizer)
	r.Post("/verify-organizer", u.organizer.VerifyOrganizer)
	r.Post("/org-login", u.organizer.Login)
}

func (u *OrganizerRouter) Profile(r chi.Router) {
	r.With(middlewere.Auth).Post("/upload-logo", u.organizer.UploadLogo)
}
