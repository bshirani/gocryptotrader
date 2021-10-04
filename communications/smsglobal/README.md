# GoCryptoTrader package Smsglobal

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">



This smsglobal package is part of the GoCryptoTrader codebase.

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


