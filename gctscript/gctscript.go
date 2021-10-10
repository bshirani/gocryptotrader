package gctscript

import (
	"gocryptotrader/gctscript/modules"
	"gocryptotrader/gctscript/wrappers/gct"
)

// Setup configures the wrapper interface to use
func Setup() {
	modules.SetModuleWrapper(gct.Setup())
}
