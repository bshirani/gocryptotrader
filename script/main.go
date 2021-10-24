package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/currency/coinmarketcap"
	"gocryptotrader/database"
	gctdatabase "gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/database/repository/instrument"
	"gocryptotrader/engine"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/log"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Syncer struct {
	bot *engine.Engine
	cmc *coinmarketcap.Coinmarketcap
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
		TradeConfigFile:               configPath,
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

	syncer := SetupSyncer(bot)
	// syncer.insertGateIOPairs()
	// res, _ := syncer.downloadCMCMap()
	syncer.saveCMCLatestListings()
}

func SetupSyncer(bot *engine.Engine) Syncer {
	settings := coinmarketcap.Settings{
		APIkey:      bot.Config.Currency.CryptocurrencyProvider.APIkey,
		AccountPlan: bot.Config.Currency.CryptocurrencyProvider.AccountPlan,
		Verbose:     true,
		Enabled:     true,
	}

	cmc := new(coinmarketcap.Coinmarketcap)
	cmc.SetDefaults()
	cmc.Setup(settings)

	return Syncer{
		bot: bot,
		cmc: cmc,
	}
}

func (s *Syncer) listAllInstruments() {
	whereQM := qm.Where("1=1")
	ins, _ := postgres.Instruments(whereQM).All(context.Background(), database.DB.SQL)
	for _, i := range ins {
		fmt.Println(i)
		// pair, _ := currency.NewPairFromStrings(i.Base, i.Quote)
		// fmt.Println(pair)
	}
}

type GateIOCoin struct {
	Base  currency.Code
	Quote currency.Code
}

func (s *Syncer) insertGateIOPairs() {
	pairs, _ := s.bot.Config.GetAvailablePairs("gateio", asset.Spot)
	// pair := pair.NewPairFromString("BTC_USD")
	for _, p := range pairs {
		err := s.insertPair(p)
		if err != nil {
			fmt.Println("error", err)
			os.Exit(123)
		}
		// upsert
	}
}

func (s *Syncer) downloadCMCMap() ([]coinmarketcap.CryptoCurrencyMap, error) {

	res, err := s.cmc.GetCryptocurrencyIDMap()
	if err != nil {
		fmt.Println("error getting map", err)
	}
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
	// f, _ := os.Create("cmcresponse.json")
	// for _, l := range res {
	// 	f.WriteString(l)
	// }
	// f.Close()
	return res, nil
}

func (s *Syncer) insertPair(p currency.Pair) error {
	if database.DB.SQL == nil {
		return database.ErrDatabaseSupportDisabled
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var tempInsert = postgres.Gateiocoin{
		Quote: p.Quote.String(),
		Base:  p.Base.String(),
	}

	err = tempInsert.Upsert(ctx, tx, true, []string{"quote", "base"}, boil.Infer(), boil.Infer())
	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			log.Errorln(log.DatabaseMgr, errRB)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Syncer) saveCMCLatestListings() error {
	listings, err := s.cmc.GetCryptocurrencyLatestListing(0, 1000)
	if err != nil {
		fmt.Println(err)
	}
	for _, l := range listings {
		s.insertCMCLatestListing(l)
	}
	return nil
}

func (s *Syncer) insertCMCLatestListing(l coinmarketcap.CryptocurrencyLatestListings) error {
	if database.DB.SQL == nil {
		return database.ErrDatabaseSupportDisabled
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	USD := l.Quote.USD

	var tempInsert = postgres.CMCLatestListing{
		CMCRank:                l.CmcRank,
		CirculatingSupply:      l.CirculatingSupply,
		DateAdded:              l.DateAdded,
		FullyDilutedMarketCap:  USD.FullyDilutedMarketCap,
		ID:                     l.ID,
		LastUpdated:            l.LastUpdated,
		MarketCap:              USD.MarketCap,
		MarketCapByTotalSupply: l.MarketCapByTotalSupply,
		MarketCapDominance:     USD.MarketCapDominance,
		MaxSupply:              l.MaxSupply,
		Name:                   l.Name,
		NumMarketPairs:         l.NumMarketPairs,
		PercentChange1H:        USD.PercentChange1H,
		PercentChange24H:       USD.PercentChange24H,
		PercentChange30D:       USD.PercentChange30D,
		PercentChange60D:       USD.PercentChange60D,
		PercentChange7D:        USD.PercentChange7D,
		PercentChange90D:       USD.PercentChange90D,
		PercentChangeVolume24H: USD.PercentChangeVolume24H,
		PercentChangeVolume30D: USD.PercentChangeVolume30D,
		PercentChangeVolume7D:  USD.PercentChangeVolume7D,
		Price:                  USD.Price,
		Slug:                   l.Slug,
		Symbol:                 l.Symbol,
		TotalMarketCap:         USD.TotalMarketCap,
		TotalSupply:            l.TotalSupply,
		Volume24H:              USD.Volume24H,
		VolumeChange24H:        USD.VolumeChange24H,
		// Volume7D:               USD.Volume7D,
		// Volume30D:              USD.Volume30D,
		// Volume24HReported:      USD.Volume24HReported,
		// Volume30DReported:      USD.Volume30DReported,
		// Volume7DReported:       USD.Volume7DReported,
	}

	err = tempInsert.Upsert(ctx, tx, true, []string{"id"}, boil.Infer(), boil.Infer())
	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			log.Errorln(log.DatabaseMgr, errRB)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
