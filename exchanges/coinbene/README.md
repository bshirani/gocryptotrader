# GoCryptoTrader package Coinbene

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">



This coinbene package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Coinbene Exchange

### Current Features

+ REST Support
+ Websocket Support

### How to enable

+ [Enable via configuration](https://github.com/thrasher-corp/gocryptotrader/tree/master/config#enable-exchange-via-config-example)

+ Individual package example below:

```go
	// Exchanges will be abstracted out in further updates and examples will be
	// supplied then
```

### How to do REST public/private calls

+ If enabled via "configuration".json file the exchange will be added to the
IBotExchange array in the ```go var bot Bot``` and you will only be able to use
the wrapper interface functions for accessing exchange data. View routines.go
for an example of integration usage with GoCryptoTrader. Rudimentary example
below:

main.go
```go
var c exchange.IBotExchange

for i := range Bot.Exchanges {
	if Bot.Exchanges[i].GetName() == "Coinbene" {
		c = Bot.Exchanges[i]
	}
}

// Public calls - wrapper functions

// Fetches current ticker information
tick, err := c.FetchTicker()
if err != nil {
	// Handle error
}

// Fetches current orderbook information
ob, err := c.FetchOrderbook()
if err != nil {
	// Handle error
}

// Private calls - wrapper functions - make sure your APIKEY and APISECRET are
// set and AuthenticatedAPISupport is set to true

// Fetches current account information
accountInfo, err := c.GetAccountInfo()
if err != nil {
	// Handle error
}
```

+ If enabled via individually importing package, rudimentary example below:

```go
// Public calls

// Fetches current ticker information
ticker, err := c.GetTicker()
if err != nil {
	// Handle error
}

// Fetches current orderbook information
ob, err := c.GetOrderbook()
if err != nil {
	// Handle error
}

// Private calls - make sure your APIKEY and APISECRET are set and
// AuthenticatedAPISupport is set to true

// GetUserInfo returns account info
accountInfo, err := c.GetUserInfo(...)
if err != nil {
	// Handle error
}

// Submits an order and the exchange and returns its tradeID
resp, err := c.SubmitOrder(...)
if err != nil {
	// Handle error
}
```

### How to do Websocket public/private calls

```go
	// Exchanges will be abstracted out in further updates and examples will be
	// supplied then
```


