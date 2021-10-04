# GoCryptoTrader Backtester: Holdings package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This holdings package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Holdings package overview

Holdings are used to calculate the holdings at any given time for a given exchange, asset, currency pair. If an order is placed, funds are removed from funding and placed under assets.
Every data event will update and calculate holdings value based on the new price. This will allow for statistics to be easily calculated at the end of a backtesting run



