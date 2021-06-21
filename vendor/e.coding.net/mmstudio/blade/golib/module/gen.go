package module

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	defaultTimerMaxBucket = 600
	defaultExecQueueLen   = 64
)

//go:generate optionGen --option_with_struct_name=true --v=true
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"InvokeTimeout":        time.Duration(time.Duration(5) * time.Second),
		"InvokeQueueLen":       int(defaultExecQueueLen),
		"TimerDispatcherLen":   int(defaultExecQueueLen),
		"StrictOrderInvoke":    true, // action need be executed in strict sequence
		"DetectLoopInvoke":     true, // detect loop invoke like : a->b->c->a
		"StepDetectLoopInvoke": 5,
		"LoopInvokeTrigger":    func(cs []string) { log.Error().Strs("chain", cs).Msg("skeleton got loop invoke") },
		"EnableGls":            true,
		"TerminateCleanQueue":  true,
		"LiveProbe":            false,
		"LiveNotifyInternal":   time.Duration(time.Duration(5) * time.Second),
		"LiveTTL":              time.Duration(time.Duration(20) * time.Second),
		"ShouldRestart":        func(s Skeleton, errMsg string) bool { return true },

		"PanicWithStack": false,                   // 主要为了测试使用
		"CatchPanic":     true,                    // 主要为了测试使用
		"OnStopped":      (func(s Skeleton))(nil), // 主要为了测试使用
		"OnTerminated":   (func(s Skeleton))(nil), // 主要为了测试使用
	}
}

func init() {
	InstallOptionsWatchDog(func(cc *Options) {
		if cc.InvokeQueueLen == 0 {
			cc.InvokeQueueLen = defaultExecQueueLen
			log.Info().Msg(fmt.Sprintf("skeleton InvokeQueueLen is zero, changed to:%d", defaultExecQueueLen))
		}
	})
}
