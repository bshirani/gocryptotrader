# GoCryptoTrader Backtester: Eventholder package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This eventholder package is part of the GoCryptoTrader codebase.

## Eventholder package overview

The event holder is a simple interface implementation which allows the backtester to iterate over the event queue.
The event holder is based on the `EventHolder` interface and is implemented by `Holder`.
It is used by `backtest.Backtester` and it accepts appending any struct which implements the `common.EventHandler` interface, eg `order.Order`


