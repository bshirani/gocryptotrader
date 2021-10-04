# GoCryptoTrader package Exchangerate.Host


This exchangerate.host package is part of the GoCryptoTrader codebase.

## Current Features for exchangerate.host

+ Fetches up to date curency data from [ExchangeRate.host API]("https://exchangerate.host")

### How to enable

+ [Enable via configuration](https://github.com/thrasher-corp/gocryptotrader/tree/master/config#enable-currency-via-config-example)

+ Individual package example below:
```go
import (
	"github.com/thrasher-corp/gocryptotrader/currency/forexprovider/base"
	"github.com/thrasher-corp/gocryptotrader/currency/forexprovider/exchangerate.host"
)

var c exchangeratehost.ExchangeRateHost

// Define configuration
newSettings := base.Settings{
	Name:             "ExchangeRateHost",
	// ...
}

c.Setup(newSettings)

rates, err := c.GetRates("USD", "EUR,AUD")
// Handle error
```


