package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type Mailer struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewMailer(name string, fromEmailAddress string, fromEmailPassword string) *Mailer {
	return &Mailer{
		name,
		fromEmailAddress,
		fromEmailPassword,
	}
}
func (mail *Mailer) SendGmail(code string, to string) error {
	e := email.NewEmail()
	content := `
	<h1>Your OTP Verification code is: </h1>
	<p> ` + code + `</p>
	
	`
	e.From = fmt.Sprintf(" %s <%s>", mail.name, mail.fromEmailAddress)
	//	To := []string{to}
	e.Subject = "Verify Your Email"
	e.HTML = []byte(content)
	e.To = []string{to}
	smtpauth := smtp.PlainAuth("", mail.fromEmailAddress, mail.fromEmailPassword, smtpAuthAddress)
	if err := e.Send(smtpServerAddress, smtpauth); err != nil {
		return err
	}
	return nil
}
