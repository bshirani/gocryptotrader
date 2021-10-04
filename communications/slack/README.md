# GoCryptoTrader package Slack

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">



This slack package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Slack Communications package

### What is Slack?

+ Slack is a code-centric collaboration hub that allows users to connect via an
app and share different types of data
+ Please visit: [Slack](https://slack.com/) for more information and account setup

### Current Features

+ Basic communication to your slack channel information includes:
	- Working status of bot

### How to enable

+ [Enable via configuration](https://github.com/thrasher-corp/gocryptotrader/tree/master/config#enable-communications-via-config-example)

+ Individual package example below:
```go
import (
"github.com/thrasher-corp/gocryptotrader/communications/slack"
"github.com/thrasher-corp/gocryptotrader/config"
)

s := new(slack.Slack)

// Define slack configuration
commsConfig := config.CommunicationsConfig{SlackConfig: config.SlackConfig{
	Name:              "Slack",
	Enabled:           true,
	Verbose:           false,
	TargetChannel:     "targetChan",
	VerificationToken: "slackGeneratedToken",
}}

s.Setup(commsConfig)
err := s.Connect
// Handle error
```

Once the bot has started you can interact with the bot using these commands
via Slack:

```
!status 		- Displays current working status of bot
!help 			- Displays help text
!settings		- Displays current settings
```


