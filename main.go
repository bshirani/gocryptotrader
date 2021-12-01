package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gocryptotrader/cmd/backtest"
	"gocryptotrader/cmd/live"
	"gocryptotrader/cmd/run"
	"gocryptotrader/core"
	"gocryptotrader/database"
	"gocryptotrader/database/repository/datahistoryjob"
	"gocryptotrader/engine"
	"gocryptotrader/log"
	"gocryptotrader/portfolio/analyze"

	"github.com/urfave/cli/v2"
)

var (
	app = &cli.App{
		Name:                 "gct",
		Version:              core.Version(false),
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			// fmt.Printf("Hello %q", c.Args().Get(0))
			runCommand(c)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       "live",
				Usage:       "config file",
				Aliases:     []string{"c"},
				Destination: &settings.ConfigFile,
			},
			&cli.StringFlag{
				Name:        "trade",
				Aliases:     []string{"t"},
				Value:       "all",
				Usage:       "trade config file",
				Destination: &settings.TradeConfigFile,
			},
			&cli.BoolFlag{
				Name:        "exchange_verbose",
				Aliases:     []string{"ev"},
				Value:       false,
				Usage:       "exchange verbose",
				Destination: &settings.EnableExchangeVerbose,
			},
			&cli.StringFlag{
				Name:        "command",
				Value:       "",
				Usage:       "command to run",
				Destination: &command,
			},
		},
		Commands: []*cli.Command{
			backtest.BacktestCommand,
			live.LiveCommand,
			run.RunCommand,
			{
				Name:   "analyze_pf",
				Usage:  "analyze pf",
				Action: analyzePF,
			},
			{
				Name:   "backtest_model",
				Usage:  "backtest model",
				Action: backtestModel,
			},
			{
				Name:   "update_weights",
				Usage:  "calculate and update pf weights",
				Action: updateWeights,
			},
			{
				Name:   "generate_all_strategies",
				Usage:  "generate all.strat",
				Action: generateAll,
			},
		},
	}
	tradeConfigPath string
	settings        engine.Settings
	workingDir      string
	verbose         bool
	bot             *engine.Engine
	command         string
)

func main() {
	app.Run(os.Args)
	if bot != nil {
		err := bot.Config.SaveConfigToFile(bot.Settings.ConfigFile)
		if err != nil {
			log.Errorln(log.Global, "Unable to save config.")
		} else {
			log.Debugln(log.Global, "Config file saved successfully.")
		}
	}
}

func runCommand(c *cli.Context) {
	cmd := c.Args().Get(0)
	settings.ConfigFile = cmd
	fmt.Println("set config file to", cmd)
	if cmd == "analyze_pf" {
		analyzePF(c)
	} else if cmd == "catchup" {
		settings.ConfigFile = "catchup"
		setupBot()
		catchup(c)
	}
}

func script(c *cli.Context) error {
	startOfflineServices()
	return nil
}

func catchup(c *cli.Context) error {
	startOfflineServices()
	db := bot.DatabaseManager.GetInstance()
	dhj, err := datahistoryjob.Setup(db)
	if err != nil {
		fmt.Println("error", err)
	}
	dhj.ClearJobs()

	log.Infoln(log.TradeMgr, "Catching up days...", bot.Config.DataHistory.DaysBack)
	daysBack := make([]int, bot.Config.DataHistory.DaysBack)
	for _, cs := range bot.CurrencySettings {
		fmt.Println(cs.CurrencyPair)
	}
	// return nil

	for i := range daysBack {
		i += 1
		bot.DataHistoryManager.CatchupDays(int64(i))

		for {
			active, err := dhj.CountActive()
			if err != nil {
				fmt.Println("error", err)
			}
			if active == 0 {
				fmt.Println("starting days back", i)
				break
			}
			time.Sleep(time.Second)
		}
	}
	return nil
}

func setupBot() error {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}

	configDir := filepath.Join(wd, "confs")

	if settings.EnableProductionMode {
		settings.ConfigFile = filepath.Join(configDir, fmt.Sprintf("prod.json", settings.ConfigFile))
		settings.TradeConfigFile = filepath.Join(configDir, fmt.Sprintf("all.strat", settings.TradeConfigFile))

	} else {
		if settings.TradeConfigFile == "" {
			settings.TradeConfigFile = filepath.Join(configDir, "all.strat")
		} else {
			settings.TradeConfigFile = filepath.Join(configDir, fmt.Sprintf("dev/strategy/%s.strat", settings.TradeConfigFile))
		}

		if settings.ConfigFile == "" {
			panic("config file is empty")
			// fmt.Println("config file", settings.ConfigFile)
			// fmt.Println("emoty config file")
			if settings.EnableLiveMode {
				settings.ConfigFile = filepath.Join(configDir, "dev/live.json")
			} else {
				settings.ConfigFile = filepath.Join(configDir, "dev/backtest.json")
			}
		} else {
			settings.ConfigFile = filepath.Join(configDir, fmt.Sprintf("dev/%s.json", settings.ConfigFile))
		}
	}

	flags := map[string]bool{}
	bot, err = engine.NewFromSettings(&settings, flags)
	if err != nil {
		fmt.Printf("Could not run engine. Error: %v.\n", err)
		os.Exit(-1)
	}

	bot.DatabaseManager, err = engine.SetupDatabaseConnectionManager(database.DB.GetConfig())
	if err != nil {
		return err
	} else {
		err = bot.DatabaseManager.Start(&bot.ServicesWG)
		if err != nil {
			log.Errorf(log.Global, "Database manager unable to start: %v", err)
		}
	}

	return err
}

func updateWeights(c *cli.Context) error {
	setupBot()
	pf, err := getPF()
	pf.Analyze()
	prodWeighted := filepath.Join(workingDir, "confs/prod.strat")
	fmt.Println("saving", len(pf.Weights.Strategies), "pf weights to", prodWeighted)
	pf.Weights.Save(prodWeighted)
	return err
}

func backtestModel(c *cli.Context) error {
	err := setupBot()
	pf, err := getPF()
	pf.BacktestModel()
	pf.Analyze()
	// fmt.Println("backtest models", pf)
	// prodWeighted := filepath.Join(workingDir, "confs/prod.strat")
	// fmt.Println("saving", len(pf.Weights.Strategies), "pf weights to", prodWeighted)
	// pf.Weights.Save(prodWeighted)
	return err
}

func analyzePF(c *cli.Context) error {
	setupBot()
	pf, err := getPF()
	pf.Analyze()
	return err
}

func generateAll(c *cli.Context) error {
	ss := analyze.GenerateAllStrategies()
	allPath := filepath.Join(workingDir, "confs/dev/strategy/all.strat")
	fmt.Println("saving all.strat to", allPath)
	analyze.SaveStrategiesConfigFile(allPath, ss)
	return nil
}

func getPF() (*analyze.PortfolioAnalysis, error) {
	return analyze.SetupPortfolio(bot.Config, "")
}

func startOfflineServices() (err error) {
	// fmt.Println("start offline services")
	if bot.Config.LiveMode {
		panic("cannot run offline services in live mode")
	}

	// err = bot.LoadExchange("FTX", nil)
	// if err != nil && !errors.Is(err, engine.ErrExchangeAlreadyLoaded) {
	// 	fmt.Println("error", err)
	// 	return err
	// }

	err = bot.SetupExchanges()
	if err != nil {
		return err
	}

	// fmt.Println("setting up exchange settings")
	err = bot.SetupExchangeSettings()
	if err != nil {
		panic(err)
	}

	if bot.Config.DataHistory.Enabled {
		if bot.DataHistoryManager == nil {
			bot.DataHistoryManager, err = engine.SetupDataHistoryManager(bot, bot.ExchangeManager, bot.DatabaseManager, &bot.Config.DataHistory)
			if err != nil {
				log.Errorf(log.Global, "database history manager unable to setup: %s", err)
			} else {
				err = bot.DataHistoryManager.Start()
				if err != nil {
					log.Errorf(log.Global, "database history manager unable to start: %s", err)
				}
			}
		}
	}

	bot.DatabaseManager, err = engine.SetupDatabaseConnectionManager(database.DB.GetConfig())
	if err != nil {
		return err
	} else {
		err = bot.DatabaseManager.Start(&bot.ServicesWG)
		if err != nil {
			log.Errorf(log.Global, "Database manager unable to start: %v", err)
		}
	}
	return err
}
