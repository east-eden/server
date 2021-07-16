package time2

import (
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/golib/suturev4"
)

var supervisorOnce sync.Once
var defaultSupervisorSpec = suturev4.Spec{
	FailureDecay:     30,                              // 30 seconds
	FailureThreshold: 5,                               //5 failures
	FailureBackoff:   time.Duration(15) * time.Second, //15 seconds
	Timeout:          time.Duration(10) * time.Second, // 10 seconds
}

var time2Supervisor *suturev4.Supervisor

func SetSkeletonSupervisor(supervisor *suturev4.Supervisor) { time2Supervisor = supervisor }
