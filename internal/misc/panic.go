package misc

import (
	"flag"
	"log/slog"
	"time"
)

func Panic(err error) {
	// if not unit testing give some time to flush logs etc.
	if flag.Lookup("test.v") == nil {
		slog.Error("PANIC ... waiting 15s to exit ...")
		<-time.After(15 * time.Second)
	}
	panic(err)
}
