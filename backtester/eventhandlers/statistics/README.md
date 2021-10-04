# GoCryptoTrader Backtester: Statistics package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This statistics package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Statistics package overview

The statistics package is used for storing all relevant data over the course of a GoCryptoTrader Backtesting run. All types of events are tracked by exchange, asset and currency pair.
When multiple currencies are included in your strategy, the statistics package will be able to calculate which exchange asset currency pair has performed the best, along with the biggest drop downs in the market.




