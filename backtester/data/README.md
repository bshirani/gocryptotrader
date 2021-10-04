# GoCryptoTrader Backtester: Data package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This data package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Data package overview

The data package defines and implements a base version of the `Streamer` interface which is part of the `Handler` interface. These interfaces allow for the translation of data into individual intervals to be accessed and assessed as part of the `backtest` package.
This is a base implementation, the more proper implementation that is used throughout the backtester is under `./kline`

This can also be used to implement other means to load data for the backtester to process, however kline is currently the only supported method.





