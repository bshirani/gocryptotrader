package analyze

import (
	"gocryptotrader/currency"
	"gocryptotrader/exchange/order"
	"time"
)

type StrategyStats struct {
	NumTrades       int     `json:"numTrades"`
	NetProfit       float64 `json:"netProfit"`
	WinPercentage   float64 `json:"winPercentage"`
	AveragePL       float64 `json:"averagePl"`
	AverageWin      float64 `json:"averageWin"`
	AverageLoss     float64 `json:"averageLoss"`
	AvgWinByAvgLoss float64 `json:"avgWinByAvgLoss"`
	// MaxDrawdown                  Swing                 `json:"maxDrawdown,omitempty"`
	// StartingClosePrice           decimal.Decimal       `json:"startingClosePrice"`
	// EndingClosePrice             decimal.Decimal       `json:"endingClosePrice"`
	// LowestClosePrice             decimal.Decimal       `json:"lowestClosePrice"`
	// HighestClosePrice            decimal.Decimal       `json:"highestClosePrice"`
	// MarketMovement               decimal.Decimal       `json:"marketMovement"`
	// StrategyMovement             decimal.Decimal       `json:"strategyMovement"`
	// HighestCommittedFunds        HighestCommittedFunds `json:"highestCommittedFunds"`
	// RiskFreeRate                 decimal.Decimal       `json:"riskFreeRate"`
	// BuyOrders                    int64                 `json:"buyOrders"`
	// GeometricRatios              Ratios                `json:"geometricRatios"`
	// ArithmeticRatios             Ratios                `json:"arithmeticRatios"`
	// CompoundAnnualGrowthRate     decimal.Decimal       `json:"compoundAnnualGrowthRate"`
	// SellOrders                   int64                 `json:"sellOrders"`
	// TotalOrders                  int64                 `json:"totalOrders"`
	// InitialHoldings              holdings.Holding      `json:"initialHoldingsHoldings"`
	// FinalHoldings                holdings.Holding      `json:"finalHoldings"`
	// FinalOrders                  compliance.Snapshot   `json:"finalOrders"`
	// ShowMissingDataWarning       bool                  `json:""`
	// IsStrategyProfitable         bool                  `json:"isStrategyProfitable"`
	// DoesPerformanceBeatTheMarket bool                  `json:"doesPerformanceBeatTheMarket"`
}

type StrategyAnalysis struct {
	// Trades                       []*livetrade.Details
	Pair      currency.Pair `json:"pair"`
	Direction order.Side    `json:"direction"`
	Name      string        `json:"name"`
	Label     string        `json:"label"`

	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`

	Base       *StrategyStats
	Prediction *StrategyStats
}
