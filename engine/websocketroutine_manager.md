# GoCryptoTrader package Websocketroutine manager

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">


## Current Features for Websocketroutine manager
+ The websocket routine manager subsystem is used process websocket data in a unified manner across enabled exchanges with websocket support
+ It can help process orders to the order manager subsystem when it receives new data
+ Logs output of ticker and orderbook updates
+ The websocket routine manager subsystem can be enabled or disabled via runtime command `-websocketroutine=false` defaulting to true
+ Logs can be customised to display values the config value `fiatDisplayCurrency` under `currencyConfig`


