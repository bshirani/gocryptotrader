package run

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gocryptotrader/config"
	gctdatabase "gocryptotrader/database"
	"gocryptotrader/database/repository/datahistoryjob"
	"gocryptotrader/engine"
	"gocryptotrader/log"

	"github.com/urfave/cli/v2"
)

var configPath, tradeConfigPath string
var settings engine.Settings

var RunCommand = &cli.Command{
	Name:   "run",
	Usage:  "run script",
	Action: runConfig,
}

func runConfig(c *cli.Context) error {
	fmt.Println(c.Args())
	flags := map[string]bool{}
	var bot *engine.Engine
	bot, err := engine.NewFromSettings(&settings, flags)
	if err != nil {
		fmt.Printf("Could not run live. Error: %v.\n", err)
		os.Exit(-1)
	}
	engine.Bot = bot
	config.Cfg = *engine.Bot.Config

	err = engine.Bot.LoadExchange("gateio", nil)
	if err != nil && !errors.Is(err, engine.ErrExchangeAlreadyLoaded) {
		fmt.Println("error", err)
		return err
	}

	err = engine.Bot.SetupExchangeSettings()
	if err != nil {
		fmt.Println("error setting up exchange settings", err)
	}

	engine.Bot.DatabaseManager, err = engine.SetupDatabaseConnectionManager(gctdatabase.DB.GetConfig())
	if err != nil {
		return err
	} else {
		err = engine.Bot.DatabaseManager.Start(&engine.Bot.ServicesWG)
		if err != nil {
			log.Errorf(log.Global, "Database manager unable to start: %v", err)
		}
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

	// fmt.Println(0)
	// var localWG sync.WaitGroup
	// localWG.Add(1)
	db := engine.Bot.DatabaseManager.GetInstance()
	dhj, err := datahistoryjob.Setup(db)
	if err != nil {
		fmt.Println("error", err)
	}
	dhj.ClearJobs()

	log.Infoln(log.TradeMgr, "Catching up days...", engine.Bot.Config.DataHistory.DaysBack)
	daysBack := make([]int, engine.Bot.Config.DataHistory.DaysBack)

	for i := range daysBack {
		i += 1
		engine.Bot.DataHistoryManager.CatchupDays(int64(i))

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
	log.Infoln(log.TradeMgr, "Done with catchup")

	return nil
}
