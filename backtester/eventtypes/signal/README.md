# GoCryptoTrader Backtester: Signal package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This signal package is part of the GoCryptoTrader codebase.

## Signal package overview

The signal event is created as a result of a data event being analysed via a strategy. Typically, there are three types of signal that should be expected `buy`, `sell` and `donothing`. An example of this is demonstrated in the RSI strategy. However, other signals can be raised such as `MissingData`.
The signal event will contain data such as price, the direction as well as the reasoning for the signal decision with the `GetWhy()` function


