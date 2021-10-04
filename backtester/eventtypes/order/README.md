# GoCryptoTrader Backtester: Order package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This order package is part of the GoCryptoTrader codebase.

## Order package overview

The Order Event Type is an event type raised after the portfolio manager has passed all its checks and wishes to make an order
It is sent to the Exchange to process and if successful, will raise a Fill Event.

The Order Event Type is based on `common.EventHandler` and `common.Directioner` while also having the following custom functions
```
	SetAmount(float64)
	GetAmount() float64
	IsOrder() bool
	GetStatus() order.Status
	SetID(id string)
	GetID() string
	GetLimit() float64
	IsLeveraged() bool
```


