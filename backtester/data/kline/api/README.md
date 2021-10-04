# GoCryptoTrader Backtester: Api package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This api package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Api package overview

This package is responsible for the loading of kline data via the API. It can retrieve candle data or trade data which is converted into candle data.
This package uses existing GoCryptoTrader exchange implementations.

See individual exchange implementations [here](/exchanges) and the interface used [here](/exchanges/interfaces.go)


