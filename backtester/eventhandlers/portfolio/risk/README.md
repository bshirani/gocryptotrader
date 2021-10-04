# GoCryptoTrader Backtester: Risk package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This risk package is part of the GoCryptoTrader codebase.

## Risk package overview

The risk manager is responsible for ensuring that no order can be made if it is deemed too risky.
Risk is currently defined by ensuring that orders cannot have too much leverage for the individual order, overall with all orders in the portfolio as well as whether there are too many orders for an individual currency

See config package [readme](/backtester/config/README.md) to view the risk related fields to customise



