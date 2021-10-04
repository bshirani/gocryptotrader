# GoCryptoTrader package Websocketroutine manager

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">


## Current Features for Websocketroutine manager
+ The websocket routine manager subsystem is used process websocket data in a unified manner across enabled exchanges with websocket support
+ It can help process orders to the order manager subsystem when it receives new data
+ Logs output of ticker and orderbook updates
+ The websocket routine manager subsystem can be enabled or disabled via runtime command `-websocketroutine=false` defaulting to true
+ Logs can be customised to display values the config value `fiatDisplayCurrency` under `currencyConfig`


### Please click GoDocs chevron above to view current GoDoc information for this package

## Contribution

Please feel free to submit any pull requests or suggest any desired features to be added.

When submitting a PR, please abide by our coding guidelines:

+ Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
+ Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
+ Code must adhere to our [coding style](https://github.com/thrasher-corp/gocryptotrader/blob/master/doc/coding_style.md).
+ Pull requests need to be based on and opened against the `master` branch.

## Donations

<img src="https://github.com/thrasher-corp/gocryptotrader/blob/master/web/src/assets/donate.png?raw=true" hspace="70">

If this framework helped you in any way, or you would like to support the developers working on it, please donate Bitcoin to:

***bc1qk0jareu4jytc0cfrhr5wgshsq8282awpavfahc***
