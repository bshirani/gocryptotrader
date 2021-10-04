# GoCryptoTrader Backtester: Size package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This size package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Size package overview

The sizing package ensures that all potential orders raised are within both the CurrencySettings limits as well as the portfolio manager's limits.
- In the event that the order is to large, the sizing package will reduce the order until it fits that limit, inclusive of fees.
- When an order is sized under the limits, an order event cannot be raised an no order will be submitted by the exchange
- The portfolio manager's sizing rules override any CurrencySettings' rules if the sizing is outside the portfolio manager's



