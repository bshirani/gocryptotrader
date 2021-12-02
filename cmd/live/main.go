package live

import (
	"fmt"
	"gocryptotrader/common"
	"gocryptotrader/config"
	"gocryptotrader/core"
	"gocryptotrader/engine"
	"gocryptotrader/gctscript"
	"gocryptotrader/log"
	"gocryptotrader/signaler"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var startDate, endDate, tradeConfigPath, configPath, templatePath, reportOutput, pairsArg string
var settings engine.Settings
var clearDB, printLogo, generateReport, dryrun, darkReport, versionFlag bool

var LiveCommand = &cli.Command{
	Name:   "live",
	Usage:  "run live",
	Action: runLive,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Value:       "live",
			Usage:       "config path",
			Destination: &settings.ConfigFile,
		},
		&cli.StringFlag{
			Name:        "trade",
			Value:       "",
			Usage:       "trade config path",
			Destination: &settings.TradeConfigFile,
		},
		&cli.BoolFlag{
			Name:        "version",
			Usage:       "get version",
			Value:       false,
			Destination: &versionFlag,
		},
		&cli.BoolFlag{
			Name:        "production",
			Usage:       "enable produciton mode (real $$$)",
			Value:       false,
			Destination: &settings.EnableProductionMode,
		},
		&cli.BoolFlag{
			Name:        "datahistory",
			Usage:       "enable data history manager",
			Value:       false,
			Destination: &settings.EnableDataHistoryManager,
		},
		&cli.BoolFlag{
			Name:        "cleardb",
			Usage:       "clear database",
			Value:       false,
			Destination: &settings.EnableClearDB,
		},
		&cli.BoolFlag{
			Name:        "db",
			Usage:       "enable database manager",
			Value:       true,
			Destination: &settings.EnableDatabaseManager,
		},
		&cli.BoolFlag{
			Name:        "dryrun",
			Usage:       "dont write to database",
			Value:       false,
			Destination: &settings.EnableDryRun,
		},
		&cli.BoolFlag{
			Name:        "sync",
			Usage:       "enable sync manager",
			Value:       false,
			Destination: &settings.EnableExchangeSyncManager,
		},
		&cli.BoolFlag{
			Name:        "grpc",
			Usage:       "enable grpc",
			Value:       false,
			Destination: &settings.EnableGRPC,
		},
		&cli.BoolFlag{
			Name:        "live",
			Usage:       "enable live mode",
			Value:       false,
			Destination: &settings.EnableLiveMode,
		},
		&cli.BoolFlag{
			Name:        "ntpclient",
			Usage:       "enable ntpclient",
			Value:       true,
			Destination: &settings.EnableNTPClient,
		},
		&cli.BoolFlag{
			Name:        "account",
			Usage:       "enable portfolio manager",
			Value:       false,
			Destination: &settings.EnablePortfolioManager,
		},
		&cli.BoolFlag{
			Name:        "dataimporter",
			Usage:       "enable data importer",
			Value:       false,
			Destination: &settings.EnableDataImporter,
		},
	},
}

func runLive(c *cli.Context) error {
	if versionFlag {
		fmt.Print(core.Version(true))
		os.Exit(0)
	}
	if printLogo {
		fmt.Print(common.ASCIILogo)
	}
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}

	configPath = filepath.Join(wd, "confs/dev", fmt.Sprintf("%s.json", configPath))

	configDir := filepath.Join(wd, "confs")
	if settings.EnableProductionMode {

		settings.ConfigFile = filepath.Join(configDir, "prod.json")
		settings.TradeConfigFile = filepath.Join(configDir, "prod.strat")

	} else {
		if settings.TradeConfigFile == "" {
			settings.TradeConfigFile = filepath.Join(configDir, "prod.strat")
		} else {
			settings.TradeConfigFile = filepath.Join(configDir, fmt.Sprintf("dev/strategy/%s.strat", settings.TradeConfigFile))
		}

		if settings.ConfigFile == "" {
			settings.ConfigFile = filepath.Join(configDir, "dev/live.json")
		} else {
			settings.ConfigFile = filepath.Join(configDir, fmt.Sprintf("dev/%s.json", settings.ConfigFile))
		}
	}

	if tradeConfigPath == "" {
		tradeConfigPath = filepath.Join(wd, "confs/prod.strat")
	} else {
		tradeConfigPath = filepath.Join(wd, "confs/dev/strategy", fmt.Sprintf("%s.strat", tradeConfigPath))
	}

	flags := map[string]bool{}
	engine.Bot, err = engine.NewFromSettings(&settings, flags)
	if err != nil {
		fmt.Printf("Could not run live. Error: %v.\n", err)
		os.Exit(-1)
	}
	config.Cfg = *engine.Bot.Config

	gctscript.Setup()

	engine.PrintSettings(&engine.Bot.Settings)
	if err = engine.Bot.Start(); err != nil {
		log.Errorf(log.Global, "Unable to start bot engine. Error: %s\n", err)
		os.Exit(1)
	}

	interrupt := signaler.WaitForInterrupt()
	log.Infof(log.Global, "Captured %v, shutdown requested.\n", interrupt)
	engine.Bot.Stop()
	log.Infoln(log.Global, "Exiting.")
	return nil
}
