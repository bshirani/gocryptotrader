package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gocryptotrader/currency/coinmarketcap"
	gctdatabase "gocryptotrader/database"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/engine"
	"gocryptotrader/exchange/order"
	"gocryptotrader/log"

	"github.com/shopspring/decimal"
)

type Syncer struct {
	bot *engine.Engine
	cmc *coinmarketcap.Coinmarketcap
}

func main() {
	var csvPath, configPath, templatePath, reportOutput, strategiesArg, pairsArg string
	var printLogo, generateReport, dryrun, darkReport bool
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}
	flag.StringVar(&configPath, "configpath", filepath.Join(wd, "../confs/dev/backtest.json"), "config")
	flag.StringVar(&csvPath, "csvpath", filepath.Join(wd, "../confs/dev/backtest.json"), "csv")

	flag.BoolVar(&dryrun, "dryrun", true, "write orders/trades to db")
	flag.BoolVar(&generateReport, "generatereport", false, "whether to generate the report file")
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

	// path := config.DefaultFilePath()
	bot, err = engine.NewFromSettings(&engine.Settings{
		ConfigFile:                    configPath,
		TradeConfigFile:               configPath,
		EnableDryRun:                  dryrun,
		EnableAllPairs:                false,
		EnableExchangeHTTPRateLimiter: true,
		EnableLiveMode:                false,
	}, flags)
	if err != nil {
		fmt.Printf("Could not load backtester. Error: %v.\n", err)
		os.Exit(-1)
	}

	err = bot.LoadExchange("gateio", nil)
	if err != nil && !errors.Is(err, engine.ErrExchangeAlreadyLoaded) {
		fmt.Println("error", err)
		return
	}

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

	// syncer := SetupSyncer(bot)
	// syncer.insertGateIOPairs()
	// res, _ := syncer.downloadCMCMap()
	// syncer.saveCMCLatestListings()

	// load the trades csv
	// generate the statistics
	// output a report and the weight configuration

	err = analyzeTrades("")
	if err != nil {
		fmt.Println("error analyzeTrades", err)
	}
}

func analyzeTrades(filepath string) error {
	fmt.Println("analyze trades", filepath)
	// load all the trades from the csv into trade details
	lf := livetrade.LastResult()
	trades, err := livetrade.LoadCSV(lf)
	fmt.Println("loaded", len(trades), "trades from", lf)
	for _, t := range trades {
		fmt.Printf("enter=%v exit=%v enter=%v exit=%v profit=%v minutes=%d\n",
			t.EntryTime,
			t.ExitTime,
			t.EntryPrice,
			t.ExitPrice,
			getProfit(t),
			getDurationMin(t),
		)
	}
	return err
}

func getProfit(trade livetrade.Details) decimal.Decimal {
	if trade.Side == order.Buy {
		return trade.ExitPrice.Sub(trade.EntryPrice)
	} else if trade.Side == order.Sell {
		return trade.EntryPrice.Sub(trade.ExitPrice)
	}
	return decimal.Decimal{}
}

func getDurationMin(trade livetrade.Details) int {
	return int(trade.ExitTime.Sub(trade.EntryTime).Minutes())
}
