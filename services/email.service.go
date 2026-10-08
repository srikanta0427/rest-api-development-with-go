package services

import "github.com/srikanta0427/rest_api_design/util"

type EmailService interface {
	SendOtp(to, subject, body string) error
}

type EmailServiceImpl struct{}

func NewEmailService() EmailService {
	return &EmailServiceImpl{}
}

func (e *EmailServiceImpl) SendOtp(
	to, subject, body string,
) error {
	return util.SendMail(to, subject, body)
}
