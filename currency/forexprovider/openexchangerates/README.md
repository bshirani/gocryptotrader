# GoCryptoTrader package Openexchangerates

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">



This openexchangerates package is part of the GoCryptoTrader codebase.

## Current Features for openexchangerates

+ Fetches up to date curency data from [Open Exchange Rates](https://openexchangerates.org/)

### How to enable

+ [Enable via configuration](https://github.com/thrasher-corp/gocryptotrader/tree/master/config#enable-currency-via-config-example)

+ Individual package example below:
```go
import (
	"github.com/thrasher-corp/gocryptotrader/currency/forexprovider/base"
	"github.com/thrasher-corp/gocryptotrader/currency/forexprovider/openexchangerates"
)

c := openexchangerates.OXR{}

// Define configuration
newSettings := base.Settings{
	Name:             "openexchangerates",
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


