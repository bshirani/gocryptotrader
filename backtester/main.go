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
	var configPath, templatePath, reportOutput string
	var printLogo, generateReport, darkReport bool
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}
	flag.StringVar(
		&configPath,
		"configpath",
		filepath.Join(
			wd,
			"..",
			"config",
			"examples",
			"trend.strat"),
		"the config containing strategy params")
	flag.StringVar(
		&templatePath,
		"templatepath",
		filepath.Join(
			wd,
			"report",
			"tpl.gohtml"),
		"the report template to use")
	flag.BoolVar(
		&generateReport,
		"generatereport",
		false,
		"whether to generate the report file")
	flag.StringVar(
		&reportOutput,
		"outputpath",
		filepath.Join(
			wd,
			"results"),
		"the path where to output results")
	flag.BoolVar(
		&printLogo,
		"printlogo",
		false,
		"print out the logo to the command line, projected profits likely won't be affected if disabled")
	flag.BoolVar(
		&darkReport,
		"darkreport",
		false,
		"sets the initial rerport to use a dark theme")
	flag.Parse()

	var bt *engine.TradeManager
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
	bt, err = engine.NewTradeManagerFromConfig(cfg, templatePath, reportOutput, bot, false)
	if err != nil {
		fmt.Printf("Could not setup backtester from config. Error: %v.\n", err)
		os.Exit(1)
	}

	err = bt.Run()
	if err != nil {
		fmt.Printf("Could not complete run. Error: %v.\n", err)
		os.Exit(1)
	}
	bt.Stop()

	// err = bt.Statistic.CalculateAllResults()
	// if err != nil {
	// 	gctlog.Error(gctlog.TradeManager, err)
	// 	os.Exit(1)
	// }

	// BACKTEST ONLY
	if generateReport {
		bt.Reports.UseDarkMode(darkReport)
		err = bt.Reports.GenerateReport()
		if err != nil {
			gctlog.Error(gctlog.TradeManager, err)
		}
	}
}
