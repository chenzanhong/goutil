package zlog

import (
	"fmt"
	"os"
	"sync"
)

var (
	globalHooks []LogHook
	hooksMutex  sync.RWMutex
)

type LogHook interface {
	OnLog(level Level, msg string, fields []Field) error
}

func RegisterLogHook(hook LogHook) {
	hooksMutex.Lock()
	defer hooksMutex.Unlock()
	globalHooks = append(globalHooks, hook)
}

// executeHooks is called within logWithFields
func executeHooks(zlogLevel Level, msg string, fields []Field) {
	hooksMutex.RLock()
	hooks := make([]LogHook, len(globalHooks))
	copy(hooks, globalHooks)
	hooksMutex.RUnlock()

	for _, hook := range hooks {
		if err := hook.OnLog(zlogLevel, msg, fields); err != nil {
			fmt.Fprintf(os.Stderr, "[zlog] LogHook error: %v\n", err)
		}
	}
}
