package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gocryptotrader/config"
	"gocryptotrader/core"
	"gocryptotrader/database"
	"gocryptotrader/database/repository"
)

var (
	configFile   string
	outputFolder string
)

var sqlboilerConfig map[string]driverConfig

type driverConfig struct {
	DBName    string   `json:"dbname,omitempty"`
	Host      string   `json:"host,omitempty"`
	Port      uint16   `json:"port,omitempty"`
	User      string   `json:"user,omitempty"`
	Pass      string   `json:"pass,omitempty"`
	Schema    string   `json:"schema,omitempty"`
	SSLMode   string   `json:"sslmode,omitempty"`
	Blacklist []string `json:"blacklist,omitempty"`
}

func main() {
	fmt.Println("GoCryptoTrader SQLBoiler config generation tool")
	fmt.Println(core.Copyright)
	fmt.Println()

	flag.StringVar(&configFile, "config", config.DefaultFilePath(), "config file to load")
	flag.StringVar(&outputFolder, "outdir", "", "overwrite default output folder")
	flag.Parse()

	var cfg config.Config
	err := cfg.LoadConfig(configFile, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	convertGCTtoSQLBoilerConfig(&cfg.Database)

	jsonOutput, err := json.MarshalIndent(sqlboilerConfig, "", " ")
	if err != nil {
		fmt.Printf("Marshal failed: %v", err)
		os.Exit(1)
	}

	path := filepath.Join(outputFolder, "sqlboiler.json")
	err = ioutil.WriteFile(path, jsonOutput, 0770)
	if err != nil {
		fmt.Printf("Write failed: %v", err)
		os.Exit(1)
	}
	fmt.Println("sqlboiler.json file created")
}

func convertGCTtoSQLBoilerConfig(c *database.Config) {
	tempConfig := driverConfig{
		Blacklist: []string{"goose_db_version"},
	}

	sqlboilerConfig = make(map[string]driverConfig)

	dbType := repository.GetSQLDialect()

	if dbType == database.DBPostgreSQL {
		dbType = "psql"
	}
	tempConfig.User = c.Username
	tempConfig.Pass = c.Password
	tempConfig.Port = c.Port
	tempConfig.Host = c.Host
	tempConfig.DBName = c.Database
	tempConfig.SSLMode = c.SSLMode

	sqlboilerConfig[dbType] = tempConfig
}
