package main

import (
	"errors"
	"fmt"

	"gocryptotrader/config"
	"gocryptotrader/database"
	dbPSQL "gocryptotrader/database/drivers/postgres"
	"gocryptotrader/database/repository"

	"github.com/urfave/cli/v2"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var (
	dbConn *database.Instance
)

func load(c *cli.Context) error {
	var conf config.Config
	err := conf.LoadConfig(c.String("config"), true)
	if err != nil {
		return err
	}

	if !conf.Database.Enabled {
		return database.ErrDatabaseSupportDisabled
	}

	err = openDBConnection(c, &conf.Database)
	if err != nil {
		return err
	}

	drv := repository.GetSQLDialect()
	fmt.Printf("Connected to: %s\n", conf.Database.Host)

	return nil
}

func openDBConnection(c *cli.Context, cfg *database.Config) (err error) {
	if c.IsSet("verbose") {
		boil.DebugMode = true
	}
	dbConn, err = dbPSQL.Connect(cfg)
	if err != nil {
		return fmt.Errorf("database failed to connect: %v, some features that utilise a database will be unavailable", err)
	}
	return nil

	return errors.New("no connection established")
}
