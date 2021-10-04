# GoCryptoTrader package cache

<img src="https://github.com/thrasher-corp/gocryptotrader/blob/master/web/src/assets/page-logo.png?raw=true" width="350px" height="350px" hspace="70">



This cache package is part of the GoCryptoTrader codebase.

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

