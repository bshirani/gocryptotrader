# GoCryptoTrader package Fixer.Io


This fixer.io package is part of the GoCryptoTrader codebase.

## Current Features for fixer.io

+ Fetches up to date curency data from [Fixer.io](https://fixer.io/)

### How to enable

+ [Enable via configuration](https://github.com/thrasher-corp/gocryptotrader/tree/master/config#enable-currency-via-config-example)

+ Individual package example below:
```go
import (
	"github.com/thrasher-corp/gocryptotrader/currency/forexprovider/base"
	"github.com/thrasher-corp/gocryptotrader/currency/forexprovider/fixer.io"
)

c := fixer.Fixer{}

// Define configuration
newSettings := base.Settings{
	Name:             "Fixer",
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


