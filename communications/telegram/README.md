# GoCryptoTrader package Telegram

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">


[![Build Status](https://github.com/thrasher-corp/gocryptotrader/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/thrasher-corp/gocryptotrader/actions/workflows/tests.yml)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/thrasher-corp/gocryptotrader/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/thrasher-corp/gocryptotrader?status.svg)](https://godoc.org/github.com/thrasher-corp/gocryptotrader/communications/telegram)
[![Coverage Status](http://codecov.io/github/thrasher-corp/gocryptotrader/coverage.svg?branch=master)](http://codecov.io/github/thrasher-corp/gocryptotrader?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thrasher-corp/gocryptotrader)](https://goreportcard.com/report/github.com/thrasher-corp/gocryptotrader)


This telegram package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Telegram Communications package

### What is telegram?

+ Telegram is a cloud-based instant messaging and voice over IP service
developed by Telegram Messenger LLP
+ Please visit: [Telegram](https://telegram.org/) for more information

### Current Features

+ Creation of bot that can retrieve
	- Bot status

	### How to enable

	+ [Enable via configuration](https://github.com/thrasher-corp/gocryptotrader/tree/master/config#enable-communications-via-config-example)

	+ Individual package example below:
	```go
	import (
	"github.com/thrasher-corp/gocryptotrader/communications/telegram"
	"github.com/thrasher-corp/gocryptotrader/config"
	)

	t := new(telegram.Telegram)

	// Define Telegram configuration
	commsConfig := config.CommunicationsConfig{TelegramConfig: config.TelegramConfig{
	Name: "Telegram",
		Enabled: true,
		Verbose: false,
	VerificationToken: "token",
	}}

	t.Setup(commsConfig)
	err := t.Connect
	// Handle error
	```

+ Once the bot has started you can interact with the bot using these commands
via Telegram:

```
/start			- Will authenticate your ID
/status			- Displays the status of the bot
/help			- Displays current command list
/settings		- Displays current bot settings
```

### Please click GoDocs chevron above to view current GoDoc information for this package


