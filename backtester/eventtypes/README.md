# GoCryptoTrader Backtester: Eventtypes package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This eventtypes package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Eventtypes overview

Event types are created after retrieving candle data. An individual candle is turned into a data event which is sent to the strategy for analysis. The event is then sent to the portfolio manager to determine whether there is appropriate funding, adequate risk and proper order sizing before raising an order event. The order event is taken to the exchange handler which will place the order and create a fill event. The fill event is used to update the portfolios individual holdings for analysis and decision making.
Below is an overview of how events are used
![workflow](https://i.imgur.com/Kup6IA9.png)



