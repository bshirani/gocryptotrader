# GoCryptoTrader package cache

<img src="https://github.com/thrasher-corp/gocryptotrader/blob/master/web/src/assets/page-logo.png?raw=true" width="350px" height="350px" hspace="70">



This cache package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Current Features for cache package

+ Basic LRU cache system with both goroutine safe (via mutex locking) and non-goroutine safe options

## How to use

##### Basic Usage:

```go
package main

import ("github.com/thrasher-corp/gocryptotrader/common/cache")

func main() {
	lruCache := cache.New(5)
	lruCache.Add("hello", "world")
	c := lruCache.Contains("hello")
	if !c {
		fmt.Println("expected cache to contain \"hello\" key")
	}

	v := lruCache.Get("hello")
	if v == nil {
		fmt.Println("expected cache to contain \"hello\" key")
	}
	fmt.Println(v)
}
```

