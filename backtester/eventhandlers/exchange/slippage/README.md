# GoCryptoTrader Backtester: Slippage package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This slippage package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Slippage package overview

Slippage refers to the difference between the expected price of a trade and the price at which the trade is executed. Slippage is used here to simulate what would occur if trading was live as no perfect conditions exist for placing orders.
Slippage is calculated in two ways in the GoCryptoTrader Backtester

### If `RealOrders` is `true`
- The orderbook is frequently requested during live cycle candle retrieval
- When the order is being calculated in the `ExecuteOrder` eventhandler, it will use the orderbook to simulate placing the order and adjust the order price

### If `RealOrders` is `false`
- The `min-slippage-percent` and `max-slippage-percent` values for the specific exchange, asset and currency pair will be used as bounds to simulate an orderbook using a random number
  - If it is a buy order, it will raise the price by a random percentage between the two values
  - If the order is a sell order, it will reduce the price by a random percentage between the two values


