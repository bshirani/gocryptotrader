package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gocryptotrader/config"
	"gocryptotrader/currency/coinmarketcap"
	gctdatabase "gocryptotrader/database"
	"gocryptotrader/engine"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/kline"
	"gocryptotrader/log"
)

type Syncer struct {
	bot *engine.Engine
	cmc *coinmarketcap.Coinmarketcap
}

func main() {
	var configPath, templatePath, reportOutput, strategiesArg, pairsArg string
	var printLogo, generateReport, dryrun, darkReport bool
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}
	flag.BoolVar(&dryrun, "dryrun", true, "write orders/trades to db")
	flag.BoolVar(&generateReport, "generatereport", false, "whether to generate the report file")
	flag.StringVar(&configPath, "configpath", filepath.Join(wd, "config", "trend.strat"), "the config containing strategy params")
	flag.StringVar(&templatePath, "templatepath", filepath.Join(wd, "report", "tpl.gohtml"), "the report template to use")
	flag.StringVar(&reportOutput, "outputpath", filepath.Join(wd, "results"), "the path where to output results")
	flag.StringVar(&strategiesArg, "strategy", "", "strategies")
	flag.StringVar(&pairsArg, "pairs", "", "pairs")
	flag.BoolVar(&printLogo, "printlogo", false, "print out the logo to the command line, projected profits likely won't be affected if disabled")
	flag.BoolVar(&darkReport, "darkreport", false, "sets the initial rerport to use a dark theme")
	flag.Parse()

	var bot *engine.Engine
	flags := map[string]bool{
		"tickersync":         false,
		"orderbooksync":      false,
		"tradesync":          false,
		"ratelimiter":        false,
		"ordermanager":       false,
		"enablecommsrelayer": false,
	}

	path := config.DefaultFilePath()
	bot, err = engine.NewFromSettings(&engine.Settings{
		ConfigFile:                    path,
		EnableDryRun:                  dryrun,
		EnableAllPairs:                false,
		EnableExchangeHTTPRateLimiter: true,
		EnableLiveMode:                false,
	}, flags)
	if err != nil {
		fmt.Printf("Could not load backtester. Error: %v.\n", err)
		os.Exit(-1)
	}

	// err = bot.LoadExchange("kraken", nil)
	// if err != nil && !errors.Is(err, engine.ErrExchangeAlreadyLoaded) {
	// 	fmt.Println("error", err)
	// 	return
	// }

	// load all the exchanges
	bot.SetupExchanges()

	// load exchange settings
	err = bot.SetupExchangeSettings()
	if err != nil {
		fmt.Println("error setting up exchange settings", err)
	}

	bot.DatabaseManager, err = engine.SetupDatabaseConnectionManager(gctdatabase.DB.GetConfig())
	if err != nil {
		return
	} else {
		err = bot.DatabaseManager.Start(&bot.ServicesWG)
		if err != nil {
			log.Errorf(log.Global, "Database manager unable to start: %v", err)
		}
	}
	k, _ := bot.ExchangeManager.GetExchangeByName("kraken")
	ep, err := k.GetEnabledPairs(asset.Spot)

	t := time.Now()
	thisMinute := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location()).Add(time.Minute * -1)
	start := thisMinute.Add(time.Minute * -10)

	k.GetHistoricCandlesExtended(
		context.TODO(),
		ep[0],
		asset.Spot,
		start,
		thisMinute,
		kline.OneMin,
	)
}
