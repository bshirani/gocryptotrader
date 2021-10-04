# GoCryptoTrader Backtester: Backtest package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This backtest package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Backtest package overview

The backtest package is the most important package of the GoCryptoTrader backtester. It is the engine which combines all elements.
It is responsible for the following functionality
- Loading settings from a provided config file
- Retrieving data
- Loading the data into assessable chunks
- Analysing the data via the `handleEvent` function
- Looping through all data
- Outputting results into a report


A flow of the application is as follows:
![workflow](https://i.imgur.com/Kup6IA9.png)



