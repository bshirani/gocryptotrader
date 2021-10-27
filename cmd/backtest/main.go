package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gocryptotrader/common"
	"gocryptotrader/engine"
	"gocryptotrader/log"
)

func main() {
	var startDate, endDate, tradeConfigPath, configPath, templatePath, reportOutput, pairsArg string
	var clearDB, printLogo, generateReport, dryrun, darkReport bool
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}
	flag.BoolVar(&dryrun, "dryrun", true, "write orders/trades to db")
	flag.BoolVar(&generateReport, "generatereport", false, "whether to generate the report file")
	flag.StringVar(&configPath, "config", "backtest", "the config containing strategy params")
	flag.StringVar(&tradeConfigPath, "trade", "", "the config containing strategy params")
	flag.BoolVar(&clearDB, "cleardb", true, "the config containing strategy params")
	flag.StringVar(&templatePath, "templatepath", filepath.Join(wd, "../portfolio/tradereport", "tpl.gohtml"), "the report template to use")
	flag.StringVar(&reportOutput, "outputpath", filepath.Join(wd, "results"), "the path where to output results")
	flag.StringVar(&startDate, "start", "", "start date")
	flag.StringVar(&endDate, "end", "", "enddate")
	flag.StringVar(&pairsArg, "pairs", "", "pairs")
	flag.BoolVar(&printLogo, "printlogo", false, "print out the logo to the command line, projected profits likely won't be affected if disabled")
	flag.BoolVar(&darkReport, "darkreport", false, "sets the initial rerport to use a dark theme")
	flag.Parse()

	configPath = filepath.Join(wd, "../confs/dev", fmt.Sprintf("%s.json", configPath))

	if printLogo {
		fmt.Print(common.ASCIILogo)
	}

	if tradeConfigPath == "" {
		tradeConfigPath = filepath.Join(wd, "../confs/prod.strat")
	} else {
		tradeConfigPath = filepath.Join(wd, "../confs/dev/strategy", fmt.Sprintf("%s.strat", tradeConfigPath))
	}
	// fmt.Println("Loading TradeManager Config", strategy)

	// cfg, err = config.ReadStrategyConfigFromFile(strategy)
	// if err != nil {
	// 	fmt.Printf("Could not read config. Error: %v. Path: %s\n", err, strategy)
	// 	os.Exit(1)
	// }
	// if endDate != "" {
	// 	cfg.DataSettings.DatabaseData.EndDate = endDate
	// }
	// path := config.DefaultFilePath()
	// if cfg.GoCryptoTraderConfigPath != "" {
	// 	path = cfg.GoCryptoTraderConfigPath
	// }

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

	// print range of backtest
	// print all symbols in backtest
	// print all strategies in backtest
	log.Infof(
		log.Global,
		"%d trades, %d strategies",
		len(tm.Portfolio.GetAllClosedTrades()),
		len(tm.Strategies),
	)
	tm.Stop()

	// err = tm.Statistic.CalculateAllResults()
	// if err != nil {
	// 	log.Error(log.Global, err)
	// 	os.Exit(1)
	// }

	// if generateReport {
	// tm.Reports.UseDarkMode(darkReport)
	// err = tm.Reports.GenerateReport()
	// if err != nil {
	// 	log.Error(log.Global, err)
	// }

	// tm.TradeReports.AddTrades(tm.Portfolio.GetAllClosedTrades())
	// tm.TradeReports.UseDarkMode(darkReport)
	// err = tm.TradeReports.GenerateReport()
	// if err != nil {
	// 	log.Error(log.Global, err)
	// }
	// }
}
