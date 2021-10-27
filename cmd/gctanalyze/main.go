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
		Name:                 "gctanalyze",
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
		Commands: []*cli.Command{
			backtestCommand,
			{
				Name:   "analyze_pf",
				Usage:  "analyze pf",
				Action: analyzePF,
			},
			{
				Name:   "calculate_weights",
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
	configFile string
	verbose    bool
	command    string
)

func main() {
	app.Run(os.Args)
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
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}
	configPath := filepath.Join(workingDir, "confs/dev/backtest.json")
	cfg, err := config.ReadConfigFromFile(configPath)
	if err != nil {
		fmt.Printf("Could not read config. Error: %v. Path: %s\n", err, configPath)
		os.Exit(1)
	}

	pf := &analyze.PortfolioAnalysis{
		Config: cfg,
	}
	err = pf.Analyze("")
	if err != nil {
		fmt.Println("error analyzeTrades", err)
	}
	return pf, err
}
