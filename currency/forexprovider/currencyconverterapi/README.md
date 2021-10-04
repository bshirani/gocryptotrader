# GoCryptoTrader package Currencyconverterapi

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">



This currencyconverterapi package is part of the GoCryptoTrader codebase.

## Current Features for currencyconverterapi

+ Fetches up to date curency data from [Currency Coverter API](https://free.currencyconverterapi.com/)

### How to enable

+ [Enable via configuration](https://github.com/thrasher-corp/gocryptotrader/tree/master/config#enable-currency-via-config-example)

+ Individual package example below:
```go
import (
	"github.com/thrasher-corp/gocryptotrader/currency/forexprovider/base"
	"github.com/thrasher-corp/gocryptotrader/currency/forexprovider/currencyconverter"
)

c := currencyconverter.CurrencyConverter{}

// Define configuration
newSettings := base.Settings{
	Name:             "CurrencyConverter",
	Enabled:          true,
	Verbose:          false,
	RESTPollingDelay: time.Duration,
	APIKey:           "key",
	APIKeyLvl:        "keylvl",
	PrimaryProvider:  true,
}

c.Setup(newSettings)

mapstringfloat, err := c.GetRates("USD", "EUR,CHY")
// Handle error
```


