package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gocryptotrader/config"
	"gocryptotrader/currency/coinmarketcap"
	"gocryptotrader/database"
	gctdatabase "gocryptotrader/database"
	modelPSQL "gocryptotrader/database/models/postgres"
	"gocryptotrader/database/repository/instrument"
	"gocryptotrader/engine"
	"gocryptotrader/log"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Syncer struct {
	bot *engine.Engine
}

func main() {
	var configPath, templatePath, reportOutput, strategiesArg, pairsArg string
	var printLogo, generateReport, dryrun, darkReport bool
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}
	flag.BoolVar(&dryrun, "dryrun", true, "write orders/trades to db")
	flag.BoolVar(&generateReport, "generatereport", false, "whether to generate the report file")
	flag.StringVar(&configPath, "configpath", filepath.Join(wd, "config", "trend.strat"), "the config containing strategy params")
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

	path := config.DefaultFilePath()
	bot, err = engine.NewFromSettings(&engine.Settings{
		ConfigFile:                    path,
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

	// insertGateIOPairs()

	syncer := Syncer{
		bot: bot,
	}
	// syncer.insertGateIOPairs()
	res, _ := syncer.downloadCMCMap()

	for _, coin := range res {
		// p, err := currency.NewPairFromString(coin.Symbol)
		// if err != nil {
		// 	fmt.Println("error creating pair from string", err, coin.Symbol)
		// }
		details := instrument.Details{
			CMCID: coin.ID,
			// Base:                p.Base,
			// Quote:               p.Quote,
			Name:                coin.Name,
			Symbol:              coin.Symbol,
			Slug:                coin.Slug,
			FirstHistoricalData: coin.FirstHistoricalData,
			LastHistoricalData:  coin.LastHistoricalData,
			Active:              coin.IsActive == 1,
		}
		err = instrument.Insert(details)
		if err != nil {
			// fmt.Println(err)
			os.Exit(123)
		}
		// instruments = append(instruments, details)
	}

}

func (s *Syncer) listAllInstruments() {
	whereQM := qm.Where("1=1")
	ins, _ := modelPSQL.Instruments(whereQM).All(context.Background(), database.DB.SQL)
	for _, i := range ins {
		fmt.Println(i)
		// pair, _ := currency.NewPairFromStrings(i.Base, i.Quote)
		// fmt.Println(pair)
	}
}

// func (s *Syncer) insertGateIOPairs() {
// 	pairs, _ := s.bot.Config.GetAvailablePairs("gateio", asset.Spot)
// 	// pair := pair.NewPairFromString("BTC_USD")
// 	for _, p := range pairs {
// 		details := instrument.Details{
// 			Base:  p.Base,
// 			Quote: p.Quote,
// 		}
// 		err := instrument.Insert(details)
// 		if err != nil {
// 			fmt.Println("error", err)
// 			os.Exit(123)
// 		}
// 		// upsert
// 	}
// }

func (s *Syncer) downloadCMCMap() ([]coinmarketcap.CryptoCurrencyMap, error) {
	settings := coinmarketcap.Settings{
		APIkey:      s.bot.Config.Currency.CryptocurrencyProvider.APIkey,
		AccountPlan: s.bot.Config.Currency.CryptocurrencyProvider.AccountPlan,
		Verbose:     false,
		Enabled:     true,
	}

	cmc := new(coinmarketcap.Coinmarketcap)
	cmc.SetDefaults()
	cmc.Setup(settings)
	return cmc.GetCryptocurrencyIDMap()
	// f, _ := os.Create("cmcresponse.json")
	// for _, l := range res {
	// 	f.WriteString(l)
	// }
	// f.Close()
}
