# GoCryptoTrader Backtester: Configbuilder package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This configbuilder package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Configbuilder package overview

### What does the config builder do?
The config builder runs you through the process of creating a strategy config (`.strat`) file. Configs can also be generated via test code under `config_test.go`.
Once the config is created, when running the backtester, you can reference it via `go run . -configpath=(path-to-strat-file)`

### How do I run it?
`go run .`

### Anything else?
The config builder will ask you all the necessary questions required to create a config file. If there is anything confusing, feel free to ask a question in our Slack group or open an issue!



