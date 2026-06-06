package mailer


import "gopkg.in/gomail.v2"

/** 
This file contains the Mailer struct and the Send function to send emails using the gomail package.
 The Mailer struct holds a dialer for sending emails and the sender's information (from and fromName). 
 The MailerInfo function creates an instance of the Mailer struct with the provided dialer and sender information. 
 The Send method on the Mailer struct constructs an email message and sends it using the dialer. The standalone Send function loads the mailer configuration, creates a dialer,
 initializes a Mailer instance, and sends the email using the Mailer's Send method.


This file contains the Mailer struct and the Send function to send emails using the gomail package.
 The Mailer struct holds a dialer for sending emails and the sender's information (from and fromName). 
 The MailerInfo function creates an instance of the Mailer struct with the provided dialer and sender information. 
 The Send method on the Mailer struct constructs an email message and sends it using the dialer. The standalone Send function loads the mailer configuration, creates a dialer, initializes a Mailer instance, and sends the email using the Mailer's Send method.
*/


type Mailer struct {
	dialer *gomail.Dialer
	from  string
	fromName string
}

func MailerInfo(dialer *gomail.Dialer, from string, fromName string) *Mailer {
	return &Mailer{
		dialer: dialer,
		from: from,
		fromName: fromName,
	}
}

func (m *Mailer) Send(to string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return m.dialer.DialAndSend(msg)
}

func Send(to string, subject string, body string) error {
	mailerConfig := Environment()
	dialer := gomail.NewDialer(mailerConfig.Host, mailerConfig.Port, mailerConfig.Username, mailerConfig.Password)
	mailer := MailerInfo(dialer, mailerConfig.From, mailerConfig.FromName)
	return mailer.Send(to, subject, body)
}

