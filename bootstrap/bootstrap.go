package main

import (
	"log"
	"os"
	"wallet/database"
	"wallet/provider/mailer"
	"wallet/internal/store"
)

/**
This file contains the bootstrap code for the application. It loads the configuration from environment variables,
establishes a connection to the database, and initializes the application struct with the necessary dependencies.

/*The mailer configuration is loaded using the Environment function from the mailer package,
which reads the necessary environment variables to populate the Initializer struct.
 This allows the application to be configured for sending emails without hardcoding sensitive information in the codebase.
*/

func Build() Settings {
	return Settings{
		addr: os.Getenv("ADDR"),
		db: ConnectionDatabase{
			addr:         os.Getenv("DB_ADDR"),
			maxOpenConns: 25,
			maxIdleConns: 25,
			maxIdleTime:  os.Getenv("DB_MAX_IDLE_TIME"),
		},
		mailer: mailer.Environment()}
}



/**
bootstrap initializes the application by loading the configuration, connecting to the database, and returning an instance of the application struct.
 It logs any errors encountered during the database connection process and exits the application if a connection cannot be established.
*/

func bootstrap() *Application {
	setting := Build()

	dbConn, err := database.New(setting.db.addr, setting.db.maxOpenConns, setting.db.maxIdleConns, setting.db.maxIdleTime)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Successfully connected to database", dbConn)

	return &Application{
		setting: setting,
		store: store.LoadInitializer(dbConn),
	}
}

	