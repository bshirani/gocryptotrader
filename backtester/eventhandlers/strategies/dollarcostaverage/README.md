# GoCryptoTrader Backtester: Dollarcostaverage package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This dollarcostaverage package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Dollarcostaverage package overview

The dollar cost average is a strategy which is designed to purchase on _every_ data candle. Unless data is missing, all output signals will be to buy.
This strategy supports simultaneous signal processing, aka `config.StrategySettings.SimultaneousSignalProcessing` set to true will use the function `OnSignals(d []data.Handler, p portfolio.Handler) ([]signal.Event, error)`. This function, like the basic `OnSignal` function, will signal to buy on every iteration.
This strategy does not support customisation



