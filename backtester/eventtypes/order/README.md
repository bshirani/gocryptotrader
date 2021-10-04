# GoCryptoTrader Backtester: Order package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This order package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Order package overview

The Order Event Type is an event type raised after the portfolio manager has passed all its checks and wishes to make an order
It is sent to the Exchange to process and if successful, will raise a Fill Event.

The Order Event Type is based on `common.EventHandler` and `common.Directioner` while also having the following custom functions
```
	SetAmount(float64)
	GetAmount() float64
	IsOrder() bool
	GetStatus() order.Status
	SetID(id string)
	GetID() string
	GetLimit() float64
	IsLeveraged() bool
```


