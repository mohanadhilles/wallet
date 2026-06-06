package mailer

import (
	"os"
)

/** * This file contains the Initializer struct and the Environment function to load mailer configuration from environment variables.
 * The Initializer struct holds the necessary fields for configuring the mailer, such as Host, Port, Username, Password, From, and FromName.
 * The Environment function reads the corresponding environment variables and returns an instance of the Initializer struct populated with those values.
 */

type Initializer struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

/*Environment loads the mailer configuration from environment variables and returns an Initializer struct. It reads the SMTP_HOST,
SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM, and SMTP_FROM_NAME environment variables to populate the fields of the Initializer struct.
 The function returns an instance of Initializer with the loaded configuration values.
*/

func Environment() Initializer {
	return Initializer{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     587,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		FromName: os.Getenv("SMTP_FROM_NAME"),
	}
}
