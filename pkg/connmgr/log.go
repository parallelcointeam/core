package connmgr

import (
	cl "git.parallelcoin.io/pod/pkg/clog"
)

// Log is the logger for the connmgr package
var Log = cl.NewSubSystem("pkg/connmgr    ", "info")
var log = Log.Ch

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(
	logger *cl.SubSystem) {
	Log = logger
	log = Log.Ch
}
