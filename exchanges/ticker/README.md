# GoCryptoTrader package Ticker

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">



This ticker package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Current Features for ticker

+ This services the exchanges package by ticker functions.

+ This package facilitates ticker generation.
+ Attaches methods to an ticker
	- Returns a string of a value

+ Gets a loaded ticker by exchange, asset type and currency pair.

+ This package is primarily used in conjunction with but not limited to the
exchange interface system set by exchange wrapper orderbook functions in
"exchange"_wrapper.go.

Examples below:

```go
tick, err := yobitExchange.FetchTicker()
if err != nil {
	// Handle error
}

// Converts ticker value to string
tickerValString := tick.PriceToString(...)
```

+ or if you have a routine setting an exchange orderbook you can access it via
the package itself.

```go
tick, err := ticker.GetTicker(...)
if err != nil {
	// Handle error
}
```


