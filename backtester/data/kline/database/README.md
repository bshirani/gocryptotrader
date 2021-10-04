# GoCryptoTrader Backtester: Database package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This database package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Database package overview

This package is responsible for the loading of kline data via a user's existing GoCryptoTrader database. It can load existing data from the `candles` and `trades` tables.
For more information on the GoCryptoTrader database, read [this readme](/database/README.md).
Ensure that your database has data and has been seeded with exchanges. For more information on this, please see [this readme](/cmd/dbseed/README.md).

### Database credentials
#### Defaults
The default database will be loaded from your GoCryptoTrader config. See [this](/database) for database configuration and implementation.

#### Overriding the GoCryptoTrader config
Database configuration details can be overridden in the `.strat` config file to allow other sources to be used and not rely on existing GoCryptoTrader configuration. See [this readme](/backtester/config/README.md) for details on config customisation


