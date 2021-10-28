package backtest

import (
	"fmt"
	"gocryptotrader/common"
	"gocryptotrader/engine"
	"gocryptotrader/log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

var startDate, endDate, tradeConfigPath, configPath, templatePath, reportOutput, pairsArg string
var clearDB, printLogo, generateReport, dryrun, darkReport bool

var BacktestCommand = &cli.Command{
	Name:   "backtest",
	Usage:  "backtest pf",
	Action: backtest,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "dryrun",
			Usage:       "dry run",
			Value:       true,
			Destination: &dryrun,
		},
		&cli.BoolFlag{
			Name:        "cleardb",
			Usage:       "clear database",
			Value:       false,
			Destination: &clearDB,
		},
		&cli.StringFlag{
			Name:        "config",
			Value:       "backtest",
			Usage:       "config path",
			Destination: &configPath,
		},
		&cli.StringFlag{
			Name:        "trade",
			Value:       "",
			Usage:       "trade config path",
			Destination: &tradeConfigPath,
		},
		&cli.StringFlag{
			Name:        "start",
			Value:       "",
			Usage:       "start date",
			Destination: &startDate,
		},
		&cli.StringFlag{
			Name:        "end",
			Value:       "",
			Usage:       "end date",
			Destination: &endDate,
		},
	},
}

func backtest(c *cli.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}

	configPath = filepath.Join(wd, "confs/dev", fmt.Sprintf("%s.json", configPath))

	if printLogo {
		fmt.Print(common.ASCIILogo)
	}

	if tradeConfigPath == "" {
		tradeConfigPath = filepath.Join(wd, "confs/prod.strat")
	} else {
		tradeConfigPath = filepath.Join(wd, "confs/dev/strategy", fmt.Sprintf("%s.strat", tradeConfigPath))
	}

	var bot *engine.Engine
	flags := map[string]bool{
		"tickersync":         false,
		"orderbooksync":      false,
		"tradesync":          false,
		"ratelimiter":        false,
		"ordermanager":       false,
		"enablecommsrelayer": false,
	}
	bot, err = engine.NewFromSettings(&engine.Settings{
		ConfigFile:                    configPath,
		TradeConfigFile:               tradeConfigPath,
		EnableDryRun:                  dryrun,
		EnableAllPairs:                false,
		EnableClearDB:                 clearDB,
		EnableExchangeHTTPRateLimiter: true,
		EnableLiveMode:                false,
	}, flags)
	if err != nil {
		fmt.Printf("Could not load backtester. Error: %v.\n", err)
		os.Exit(-1)
	}

	// should do this somewhere else
	tformat := "2006-01-02"
	if startDate != "" {
		start, err := time.Parse(tformat, startDate)
		if err != nil {
			fmt.Println("error date", err)
		}
		bot.Config.DataSettings.DatabaseData.StartDate = start
	}
	if endDate != "" {
		end, err := time.Parse(tformat, endDate)
		if err != nil {
			fmt.Println("error date", err)
		}
		bot.Config.DataSettings.DatabaseData.EndDate = end
	}
	// set default date somewhere else.

	var tm *engine.TradeManager
	tm, err = engine.NewTradeManager(bot)
	if err != nil {
		fmt.Printf("Could not setup trade manager from config. Error: %v.\n", err)
		os.Exit(1)
	}

	err = tm.Run()
	if err != nil {
		fmt.Printf("Could not complete run. Error: %v.\n", err)
		os.Exit(1)
	}

	log.Infof(
		log.Global,
		"%d trades, %d strategies",
		len(tm.Portfolio.GetAllClosedTrades()),
		len(tm.Strategies),
	)
	tm.Stop()
	return err
}
