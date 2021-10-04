# GoCryptoTrader package Event manager

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">


## Current Features for Event manager
+ The event manager subsystem is used to push events to communication systems such as Slack
+ The only configurable aspects of the event manager are the delays between receiving an event and pushing it and enabling verbose:

### connectionMonitor

| Config | Description | Example |
| ------ | ----------- | ------- |
| eventmanagerdelay | Sets the event managers sleep delay between event checking by a Golang `time.Duration` |  `0` |
| verbose | Outputs debug messaging allowing for greater transparency for what the event manager is doing |  `false` |


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
