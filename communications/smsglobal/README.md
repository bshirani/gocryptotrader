# GoCryptoTrader package Smsglobal

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">



This smsglobal package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## SMSGlobal Communications package

### What is SMSGlobal?

+ SMSGlobal allows bulk sending of messages via their API
+ Please visit: [SMSGlobal](https://www.smsglobal.com/) for more information and account setup

### Current Features

+ Sending of events to a list of recipients

### How to enable

+ [Enable via configuration](https://github.com/thrasher-corp/gocryptotrader/tree/master/config#enable-communications-via-config-example)

+ Individual package example below:
```go
import (
"github.com/thrasher-corp/gocryptotrader/communications/smsglobal"
"github.com/thrasher-corp/gocryptotrader/config"
)

s := new(smsglobal.SMSGlobal)

// Define SMSGlobal configuration
commsConfig := config.CommunicationsConfig{SMSGlobalConfig: config.SMSGlobalConfig{
	Name:     "SMSGlobal",
	Enabled:  true,
	Verbose:  false,
	Username: "username",
	Password: "password",
	Contacts: []config.SMSContact{}
}}

s.Setup(commsConfig)
err := s.Connect
// Handle error
```


