# GoCryptoTrader Backtester: Risk package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This risk package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Risk package overview

The risk manager is responsible for ensuring that no order can be made if it is deemed too risky.
Risk is currently defined by ensuring that orders cannot have too much leverage for the individual order, overall with all orders in the portfolio as well as whether there are too many orders for an individual currency

See config package [readme](/backtester/config/README.md) to view the risk related fields to customise



