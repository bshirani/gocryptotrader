package currency

import (
	"time"

	"gocryptotrader/currency/coinmarketcap"
)

// MainConfiguration is the main configuration from the config.json file
type MainConfiguration struct {
	ForexProviders         []FXSettings
	CryptocurrencyProvider coinmarketcap.Settings
	Cryptocurrencies       Currencies
	CurrencyPairFormat     interface{}
	FiatDisplayCurrency    Code
	CurrencyDelay          time.Duration
	FxRateDelay            time.Duration
}

// BotOverrides defines a bot overriding factor for quick running currency
// subsystems
type BotOverrides struct {
	Coinmarketcap       bool
	FxCurrencyConverter bool
	FxCurrencyLayer     bool
	FxFixer             bool
	FxOpenExchangeRates bool
	FxExchangeRateHost  bool
}

// CoinmarketcapSettings refers to settings
type CoinmarketcapSettings coinmarketcap.Settings

// SystemsSettings defines incoming system settings
type SystemsSettings struct {
	Coinmarketcap     coinmarketcap.Settings
	Currencyconverter FXSettings
	Currencylayer     FXSettings
	Fixer             FXSettings
	Openexchangerates FXSettings
}

// FXSettings defines foreign exchange requester settings
type FXSettings struct {
	Name             string        `json:"name"`
	Enabled          bool          `json:"enabled"`
	Verbose          bool          `json:"verbose"`
	RESTPollingDelay time.Duration `json:"restPollingDelay"`
	APIKey           string        `json:"apiKey"`
	APIKeyLvl        int           `json:"apiKeyLvl"`
	PrimaryProvider  bool          `json:"primaryProvider"`
}

// File defines a full currency file generated by the currency storage
// analysis system
type File struct {
	LastMainUpdate interface{} `json:"lastMainUpdate"`
	Cryptocurrency []Item      `json:"cryptocurrencies"`
	FiatCurrency   []Item      `json:"fiatCurrencies"`
	UnsetCurrency  []Item      `json:"unsetCurrencies"`
	Contracts      []Item      `json:"contracts"`
	Token          []Item      `json:"tokens"`
}

// Const here are packaged defined delimiters
const (
	UnderscoreDelimiter   = "_"
	DashDelimiter         = "-"
	ForwardSlashDelimiter = "/"
	ColonDelimiter        = ":"
)

// delimiters is a delimiter list
var delimiters = []string{UnderscoreDelimiter,
	DashDelimiter,
	ForwardSlashDelimiter,
	ColonDelimiter}
