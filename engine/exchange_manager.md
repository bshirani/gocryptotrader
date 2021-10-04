# GoCryptoTrader package Exchange manager

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">


## Current Features for Exchange manager
+ The exchange manager subsystem is used load and store exchanges so that the engine Bot can use them to track orderbooks, submit orders etc etc
+ The exchange manager itself is not customisable, it is always enabled.
+ The exchange manager by default will load all exchanges that are enabled in your config, however, it will also load exchanges by request via GRPC commands

