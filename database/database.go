package database

import (
	"context"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


/**
New establishes a new database connection using GORM with the provided configuration parameters.	
It returns a pointer to the GORM DB instance or an error if the connection fails.

Parameters:
- addr: The database connection string (e.g., "postgres://user:password@localhost:5432/dbname?sslmode=disable").
- maxOpenConns: The maximum number of open connections to the database.
- maxIdleConns: The maximum number of idle connections in the pool.
- maxIdleTime: The maximum amount of time a connection may be idle before being closed (e.g., "15m" for 15 minutes).

Returns:
- A pointer to the GORM DB instance if the connection is successful.
- An error if there is an issue connecting to the database or configuring the connection pool.
*/


func New(addr string, maxOpenConns, maxIdleConns int, maxIdletime string) (*gorm.DB, error) {
	duration, err := time.ParseDuration(maxIdletime)
	if err != nil {
		return nil, err
	}

	gdb, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return gdb, nil
}