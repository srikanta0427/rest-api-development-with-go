package model

type PendingOrganizer struct {
	Organizer *Organizer
	OTP       string
	Verified  bool
}
