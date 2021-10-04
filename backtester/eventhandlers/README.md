# GoCryptoTrader Backtester: Eventhandlers package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This eventhandlers package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Eventhandlers overview

Event handlers are responsible for taking in an event, analysing its contents and outputting another event to be handled. An individual candle is turned into a data event which handled via the strategy event handler. The strategy handler outputs a signal event, which the portfolio eventhandler will size and risk analyse before raising an order event. The event is then sent to the portfolio manager to determine whether there is appropriate funding, adequate risk and proper order sizing before raising an order event. The order event is taken to the exchange handler which will place the order and create a fill event.
Below is an overview of how event handlers are used
![workflow](https://i.imgur.com/Kup6IA9.png)



