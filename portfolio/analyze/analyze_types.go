package analyze

import (
	"gocryptotrader/config"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/portfolio/compliance"
	"gocryptotrader/portfolio/holdings"
	"gocryptotrader/portfolio/strategies"
	"time"

	"github.com/shopspring/decimal"
)

type TradeCSVData struct {
	Trades []*livetrade.Details
	Path   string
}

type DetailedTrade struct {
	livetrade.Details
	DurationMinutes int
}

type PortfolioAnalysis struct {
	Strategies map[int]strategies.Handler
	Report     *PortfolioReport
	Weights    *PortfolioWeights
}

type PortfolioWeights struct {
	Strategies []*config.StrategySetting `json:"strategies"`
}

type PortfolioReport struct {
	StrategiesAnalyses map[int]*StrategyAnalysis
	NumTrades          int64
	NumStrategies      int64
	AverageDurationMin float64
}

// CurrencyStatistic Holds all events and statistics relevant to an exchange, asset type and currency pair
type StrategyAnalysis struct {
	// Trades                       []*livetrade.Details
	NumTrades                    int                   `json:"num-trades"`
	MaxDrawdown                  Swing                 `json:"max-drawdown,omitempty"`
	StartingClosePrice           decimal.Decimal       `json:"starting-close-price"`
	EndingClosePrice             decimal.Decimal       `json:"ending-close-price"`
	LowestClosePrice             decimal.Decimal       `json:"lowest-close-price"`
	HighestClosePrice            decimal.Decimal       `json:"highest-close-price"`
	MarketMovement               decimal.Decimal       `json:"market-movement"`
	StrategyMovement             decimal.Decimal       `json:"strategy-movement"`
	HighestCommittedFunds        HighestCommittedFunds `json:"highest-committed-funds"`
	RiskFreeRate                 decimal.Decimal       `json:"risk-free-rate"`
	BuyOrders                    int64                 `json:"buy-orders"`
	GeometricRatios              Ratios                `json:"geometric-ratios"`
	ArithmeticRatios             Ratios                `json:"arithmetic-ratios"`
	CompoundAnnualGrowthRate     decimal.Decimal       `json:"compound-annual-growth-rate"`
	SellOrders                   int64                 `json:"sell-orders"`
	TotalOrders                  int64                 `json:"total-orders"`
	InitialHoldings              holdings.Holding      `json:"initial-holdings-holdings"`
	FinalHoldings                holdings.Holding      `json:"final-holdings"`
	FinalOrders                  compliance.Snapshot   `json:"final-orders"`
	ShowMissingDataWarning       bool                  `json:"-"`
	IsStrategyProfitable         bool                  `json:"is-strategy-profitable"`
	DoesPerformanceBeatTheMarket bool                  `json:"does-performance-beat-the-market"`
}

// Ratios stores all the ratios used for statistics
type Ratios struct {
	SharpeRatio      decimal.Decimal `json:"sharpe-ratio"`
	SortinoRatio     decimal.Decimal `json:"sortino-ratio"`
	InformationRatio decimal.Decimal `json:"information-ratio"`
	CalmarRatio      decimal.Decimal `json:"calmar-ratio"`
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
