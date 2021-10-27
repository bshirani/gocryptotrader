package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gocryptotrader/config"
	"gocryptotrader/core"
	"gocryptotrader/portfolio/analyze"

	"github.com/urfave/cli/v2"
)

var (
	app = &cli.App{
		Name:                 "analyze",
		Version:              core.Version(false),
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "command",
				Value:       "",
				Usage:       "command to run",
				Destination: &command,
			},
		},
	}
	workingDir string
	configFile string
	verbose    bool
	command    string
)

func main() {
	// flag.StringVar(&command, "command", "", "command to run status|up|up-by-one|up-to|down|create")
	// flag.Parse()

	err := app.Run(os.Args)

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}
	configPath := filepath.Join(workingDir, "../confs/dev/backtest.json")
	cfg, err := config.ReadConfigFromFile(configPath)
	if err != nil {
		fmt.Printf("Could not read config. Error: %v. Path: %s\n", err, configPath)
		os.Exit(1)
	}

	pf := &analyze.PortfolioAnalysis{
		Config: cfg,
	}

	if command == "" {
		fmt.Println("no command given")
		return
	} else if command == "pf" {
		analyzePF(pf)
		return
	} else if command == "weights" {
		pf.Analyze("")
		prodWeighted := filepath.Join(workingDir, "../confs/prod.strat")
		fmt.Println("saving", len(pf.Weights.Strategies), "pf weights to", prodWeighted)
		pf.Weights.Save(prodWeighted)
	} else if command == "store_all" {
		allPath := filepath.Join(workingDir, "../confs/dev/strategy/all.strat")
		fmt.Println("saving all.strat to", allPath)
		pf.SaveAllStrategiesConfigFile(allPath)
	}
}

func analyzePF(pf *analyze.PortfolioAnalysis) (err error) {
	err = pf.Analyze("")
	if err != nil {
		fmt.Println("error analyzeTrades", err)
	}
	filename := fmt.Sprintf(
		"portfolio_analysis_%v.json",
		time.Now().Format("2006-01-02-15-04-05"))
	filename = filepath.Join(workingDir, "results", filename)
	fmt.Println("saved portfolio analysis to")
	fmt.Println(filename)
	return pf.Save(filename)
}
