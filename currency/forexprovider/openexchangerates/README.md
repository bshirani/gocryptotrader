# GoCryptoTrader package Openexchangerates

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">



This openexchangerates package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

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


