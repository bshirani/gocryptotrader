# GoCryptoTrader Backtester: Signal package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This signal package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Signal package overview

The signal event is created as a result of a data event being analysed via a strategy. Typically, there are three types of signal that should be expected `buy`, `sell` and `donothing`. An example of this is demonstrated in the RSI strategy. However, other signals can be raised such as `MissingData`.
The signal event will contain data such as price, the direction as well as the reasoning for the signal decision with the `GetWhy()` function


