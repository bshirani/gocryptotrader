package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gocryptotrader/common"
	"gocryptotrader/config"
	gctconfig "gocryptotrader/config"
	"gocryptotrader/engine"
	gctlog "gocryptotrader/log"
)

func main() {
	var configPath, templatePath, reportOutput, strategiesArg, pairsArg string
	var printLogo, generateReport, darkReport bool
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}
	flag.StringVar(&configPath, "configpath", filepath.Join(wd, "config", "trend.strat"), "the config containing strategy params")
	flag.StringVar(&templatePath, "templatepath", filepath.Join(wd, "report", "tpl.gohtml"), "the report template to use")
	flag.BoolVar(&generateReport, "generatereport", false, "whether to generate the report file")
	flag.StringVar(&reportOutput, "outputpath", filepath.Join(wd, "results"), "the path where to output results")
	flag.StringVar(&strategiesArg, "strategy", "", "strategies")
	flag.StringVar(&pairsArg, "pairs", "", "pairs")
	flag.BoolVar(&printLogo, "printlogo", false, "print out the logo to the command line, projected profits likely won't be affected if disabled")
	flag.BoolVar(&darkReport, "darkreport", false, "sets the initial rerport to use a dark theme")
	flag.Parse()

	var tm *engine.TradeManager
	var cfg *config.Config
	cfg, err = config.ReadConfigFromFile(configPath)
	if err != nil {
		fmt.Printf("Could not read config. Error: %v. Path: %s\n", err, configPath)
		os.Exit(1)
	}

	if printLogo {
		fmt.Print(common.ASCIILogo)
	}

	path := gctconfig.DefaultFilePath()
	if cfg.GoCryptoTraderConfigPath != "" {
		fmt.Println("using custom config", path)
		path = cfg.GoCryptoTraderConfigPath
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
		ConfigFile:                    path,
		EnableDryRun:                  true,
		EnableAllPairs:                false,
		EnableExchangeHTTPRateLimiter: false,
	}, flags)
	if err != nil {
		fmt.Printf("Could not load backtester. Error: %v.\n", err)
		os.Exit(-1)
	}

	err = cfg.Validate()
	if err != nil {
		fmt.Printf("Could not read config. Error: %v.\n", err)
		os.Exit(1)
	}
	tm, err = engine.NewTradeManagerFromConfig(cfg, templatePath, reportOutput, bot, false)
	if err != nil {
		fmt.Printf("Could not setup trade manager from config. Error: %v.\n", err)
		os.Exit(1)
	}

	err = tm.Run()
	if err != nil {
		fmt.Printf("Could not complete run. Error: %v.\n", err)
		os.Exit(1)
	}

	// for _, t := range tm.Portfolio.GetAllClosedTrades() {
	// 	if t != nil {
	// 		fmt.Println("trade:", t)
	// 	}
	// }

	tm.Stop()

	// err = tm.Statistic.CalculateAllResults()
	// if err != nil {
	// 	gctlog.Error(gctlog.TradeManager, err)
	// 	os.Exit(1)
	// }

	if generateReport {
		tm.Reports.UseDarkMode(darkReport)
		err = tm.Reports.GenerateReport()
		if err != nil {
			gctlog.Error(gctlog.TradeManager, err)
		}
	}
}
