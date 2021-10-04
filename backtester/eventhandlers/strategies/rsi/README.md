# GoCryptoTrader Backtester: Rsi package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This rsi package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Rsi package overview

The RSI strategy utilises [the gct-ta RSI package](https://github.com/thrasher-corp/gct-ta) to analyse market signals and output buy or sell signals based on the RSI output.
This strategy does support `SimultaneousSignalProcessing` aka [use-simultaneous-signal-processing](/backtester/config/README.md).
This strategy does support strategy customisation in the following ways:

| Field | Description |  Example |
| --- | ------- | --- |
|rsi-high| The upper bounds of RSI that when met, will trigger a Sell signal | 70 |
|rsi-low| The lower bounds of RSI that when met, will trigger a Buy signal | 30 |
|rsi-period| The consecutive candle periods used in order to generate a value. All values less than this number cannot output a buy or sell signal | 14 |


