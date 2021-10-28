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
	"gocryptotrader/engine"
	"gocryptotrader/log"
	"gocryptotrader/portfolio/analyze"

	"github.com/urfave/cli/v2"
)

var configPath, tradeConfigPath string
var settings engine.Settings

var (
	app = &cli.App{
		Name:                 "gct",
		Version:              core.Version(false),
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       "live",
				Usage:       "config file",
				Destination: &settings.ConfigFile,
			},
			&cli.StringFlag{
				Name:        "trade",
				Value:       "all",
				Usage:       "trade config file",
				Destination: &settings.TradeConfigFile,
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
	workingDir string
	verbose    bool
	bot        *engine.Engine
)

func main() {
	app.Run(os.Args)
	err := bot.Config.SaveConfigToFile(bot.Settings.ConfigFile)
	if err != nil {
		log.Errorln(log.Global, "Unable to save config.")
	} else {
		log.Debugln(log.Global, "Config file saved successfully.")
	}
}

func startBot() error {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}

	configPath = filepath.Join(wd, "confs/dev", fmt.Sprintf("%s.json", settings.ConfigFile))
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
			settings.ConfigFile = filepath.Join(configDir, "dev/live.json")
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
	return err
}

func updateWeights(c *cli.Context) error {
	pf, err := getPF()
	prodWeighted := filepath.Join(workingDir, "confs/prod.strat")
	fmt.Println("saving", len(pf.Weights.Strategies), "pf weights to", prodWeighted)
	pf.Weights.Save(prodWeighted)
	return err
}

func analyzePF(c *cli.Context) error {
	pf, err := getPF()
	filename := fmt.Sprintf(
		"portfolio_analysis_%v.json",
		time.Now().Format("2006-01-02-15-04-05"))
	filename = filepath.Join(workingDir, "results/pf", filename)
	fmt.Println("saved portfolio analysis to")
	fmt.Println(filename)
	pf.Save(filename)
	return err
}

func generateAll(c *cli.Context) error {
	pf, err := getPF()
	allPath := filepath.Join(workingDir, "confs/dev/strategy/all.strat")
	fmt.Println("saving all.strat to", allPath)
	pf.SaveAllStrategiesConfigFile(allPath)
	return err
}

func getPF() (*analyze.PortfolioAnalysis, error) {
	startBot()
	pf := &analyze.PortfolioAnalysis{
		Config: bot.Config,
	}
	err := pf.Analyze("")
	if err != nil {
		fmt.Println("error analyzeTrades", err)
	}
	return pf, err
}
