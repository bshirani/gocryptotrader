package repository

import (
	"gocryptotrader/database"
)

// GetSQLDialect returns current SQL Dialect based on enabled driver
func GetSQLDialect() string {
	cfg := database.DB.GetConfig()
	switch cfg.Driver {
	case "psql", "postgres", "postgresql":
		return database.DBPostgreSQL
	}
	return "invalid driver"
}
