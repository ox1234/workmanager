package workmanager

import (
	"context"
	"fmt"
	"runtime"
)

// catchPanic catch panic
func catchPanic(ctx context.Context, format string, args ...any) {
	if e := recover(); e != nil {
		log.Errorf(format+": %v\n%v", append(args, e, catchStack())...)
	}
}

// wrapPanic wrap panic
func wrapPanic(format string, args ...any) {
	if e := recover(); e != nil {
		panic(fmt.Errorf(format+": %v\n%v", append(args, e, catchStack())...))
	}
}

// catchStack catch stack info
func catchStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return fmt.Sprintf("==> %s\n", string(buf[:n]))
}
