package main

import (
	"wallet/internal/store"
	"wallet/provider/mailer"
)

/**
This file contains the bootstrap code for the application. It loads the configuration from environment variables,
establishes a connection to the database, and initializes the application struct with the necessary dependencies.

The mailer configuration is loaded using the Environment function from the mailer package,
which reads the necessary environment variables to populate the Initializer struct.
 This allows the application to be configured for sending emails without hardcoding sensitive information in the codebase.
*/
type ConnectionDatabase struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}


type Settings struct {
	addr   string
	db     ConnectionDatabase
	mailer mailer.Initializer
}

type Application struct {
	setting Settings
	store  store.Initializer
}

