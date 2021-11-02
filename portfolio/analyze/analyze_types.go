package analyze

import (
	"gocryptotrader/config"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/portfolio/strategies"
	"time"

	"github.com/shopspring/decimal"
)

type PortfolioAnalysis struct {
	AllSettings        []*config.StrategySetting `json:"strategies"`
	GroupedSettings    []*config.StrategySetting `json:"strategies"`
	Strategies         []strategies.Handler
	Report             *Report
	Weights            *PortfolioWeights
	groupedTrades      map[string][]*livetrade.Details
	trades             []*livetrade.Details
	Config             *config.Config
	Filepath           string
	StrategiesAnalyses []*StrategyAnalysis
}

type TradeCSVData struct {
	Trades []*livetrade.Details
	Path   string
}

type DetailedTrade struct {
	livetrade.Details
	DurationMinutes int
}

type PortfolioWeights struct {
	Strategies []*config.StrategySetting `json:"strategies"`
}

type Report struct {
	Portfolio  *PortfolioReport    `json:"portfolio"`
	Strategies []*StrategyAnalysis `json:"strategies"`
}

type PortfolioReport struct {
	NumTrades          int64
	NumStrategies      int64
	AverageDurationMin float64
}

// Swing holds a drawdown
type Swing struct {
	Highest          Iteration       `json:"highest"`
	Lowest           Iteration       `json:"lowest"`
	DrawdownPercent  decimal.Decimal `json:"drawdown"`
	IntervalDuration int64
}

// Iteration is an individual iteration of price at a time
type Iteration struct {
	Time  time.Time       `json:"time"`
	Price decimal.Decimal `json:"price"`
}

// HighestCommittedFunds is an individual iteration of price at a time
type HighestCommittedFunds struct {
	Time  time.Time       `json:"time"`
	Value decimal.Decimal `json:"value"`
}

// type PortfolioStats interface {
// 	TotalEquityReturn() (decimal.Decimal, error)
// 	MaxDrawdown() Swing
// 	LongestDrawdown() Swing
// 	SharpeRatio(decimal.Decimal) decimal.Decimal
// 	SortinoRatio(decimal.Decimal) decimal.Decimal
// }

// // Handler contains all functions required to generate statistical reporting for backtesting results
// type Handler interface {
// 	GenerateReport() error
// 	AddKlineItem(*kline.Item)
// 	UpdateItem(*kline.Item)
// 	UseDarkMode(bool)
// }

// DetailedKline enhances kline details for the purpose of rich reporting results
// type DetailedKline struct {
// 	IsOverLimit bool
// 	Watermark   string
// 	Exchange    string
// 	Asset       asset.Item
// 	Pair        currency.Pair
// 	Interval    kline.Interval
// 	Candles     []DetailedCandle
// }
//
// // DetailedCandle contains extra details to enable rich reporting results
// type DetailedCandle struct {
// 	Time           int64
// 	Open           decimal.Decimal
// 	High           decimal.Decimal
// 	Low            decimal.Decimal
// 	Close          decimal.Decimal
// 	Volume         decimal.Decimal
// 	VolumeColour   string
// 	MadeOrder      bool
// 	OrderDirection order.Side
// 	OrderAmount    decimal.Decimal
// 	Shape          string
// 	Text           string
// 	Position       string
// 	Colour         string
// 	PurchasePrice  decimal.Decimal
// }

// type TradeData struct {
// 	OriginalCandles []*kline.Item
// 	EnhancedCandles []DetailedKline
// 	Statistics      *statistics.Statistic
// 	Config          *config.Config
// 	TemplatePath    string
// 	OutputPath      string
// 	Warnings        []Warning
// 	UseDarkTheme    bool
// }
