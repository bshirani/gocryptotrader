package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thrasher-corp/gocryptotrader/config"
	gctconfig "github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/engine"
	"github.com/thrasher-corp/gocryptotrader/factors"
	gctlog "github.com/thrasher-corp/gocryptotrader/log"
	"github.com/thrasher-corp/gocryptotrader/signaler"
)

func main() {
	var configPath, templatePath, reportOutput string
	var printLogo, generateReport, darkReport, live bool
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
		&live,
		"live",
		false,
		"run the system live")
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
		true,
		"print out the logo to the command line, projected profits likely won't be affected if disabled")
	flag.BoolVar(
		&darkReport,
		"darkreport",
		false,
		"sets the initial rerport to use a dark theme")
	flag.Parse()

	var bt *engine.BackTest
	var cfg *config.Config
	cfg, err = config.ReadConfigFromFile(configPath)
	if err != nil {
		fmt.Printf("Could not read config. Error: %v.\n", err)
		os.Exit(1)
	}

	// if printLogo {
	// 	fmt.Print(common.ASCIILogo)
	// }

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
		"ratelimiter":        true,
		"ordermanager":       true,
		"enablecommsrelayer": true,
	}
	bot, err = engine.NewFromSettings(&engine.Settings{
		ConfigFile:                    path,
		EnableDryRun:                  true,
		EnableAllPairs:                true,
		EnableExchangeHTTPRateLimiter: true,
		IsLive:                        live,
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
	bt, err = backtest.NewFromConfig(cfg, templatePath, reportOutput, bot)
	if err != nil {
		fmt.Printf("Could not setup backtester from config. Error: %v.\n", err)
		os.Exit(1)
	}

	bt.DataHistoryManager, err = engine.SetupDataHistoryManager(bot.ExchangeManager, bot.DatabaseManager, &bot.Config.DataHistoryManager)
	if err != nil {
		gctlog.Errorf(gctlog.Global, "database history manager unable to setup: %s", err)
	} else {
		err = bt.DataHistoryManager.Start()
		if err != nil {
			gctlog.Errorf(gctlog.Global, "database history manager unable to start: %s", err)
		}
	}

	if live {
		e, _ := factors.Setup()
		bt.FactorEngine = e
		// run catchup here
		// fmt.Println("live mode, running catchup")
		// bt.Catchup()
		// fmt.Println("catchup completed")
		go func() {
			err = bt.RunLive()
			if err != nil {
				fmt.Printf("Could not complete live run. Error: %v.\n", err)
				os.Exit(-1)
			}
		}()
		interrupt := signaler.WaitForInterrupt()
		gctlog.Infof(gctlog.Global, "Captured %v, shutdown requested.\n", interrupt)
		bt.Stop()
	} else {
		err = bt.Run()
		if err != nil {
			fmt.Printf("Could not complete run. Error: %v.\n", err)
			os.Exit(1)
		}
		bt.Stop()
	}

	// err = bt.Statistic.CalculateAllResults()
	// if err != nil {
	// 	gctlog.Error(gctlog.BackTester, err)
	// 	os.Exit(1)
	// }

	// BACKTEST ONLY
	if generateReport {
		bt.Reports.UseDarkMode(darkReport)
		err = bt.Reports.GenerateReport()
		if err != nil {
			gctlog.Error(gctlog.BackTester, err)
		}
	}
}
