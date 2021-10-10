package settings

import (
	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/bt_config"
	"github.com/thrasher-corp/gocryptotrader/compliance"
	"github.com/thrasher-corp/gocryptotrader/holdings"
)

// Settings holds all important information for the portfolio manager
// to assess purchasing decisions
type Settings struct {
	Fee               decimal.Decimal
	BuySideSizing     config.MinMax
	SellSideSizing    config.MinMax
	Leverage          config.Leverage
	HoldingsSnapshots []holdings.Holding
	ComplianceManager compliance.Manager
}
