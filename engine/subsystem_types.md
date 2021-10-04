# GoCryptoTrader package Subsystem types

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">


## Current Features for Subsystem types
+ Subsystem contains subsystems that are used at run time by an `engine.Engine`, however they can be setup and run individually.
+ Subsystems are designed to be self contained
+ All subsystems have a public `Setup(...) (..., error)` function to return a valid subsystem ready for use
  + Subsystems which are designed to be switched off also have `Start(...) error`, `IsRunning() bool` and `Stop(...) error` functions to allow the main `engine.Engine` instance to manage them
+ Common subsystem types such as errors can be found within the `subsystem.go` file

