# GoCryptoTrader package Exchangerate.Host

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">



This exchangerate.host package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

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


