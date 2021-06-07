package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"e.coding.net/mmstudio/blade/server/define"
	"github.com/sony/sonyflake"
)

type Snowflakes struct {
	ids  []*sonyflake.Sonyflake
	once sync.Once
	cb   func()
}

var sfs Snowflakes

func init() {
	sfs.ids = make([]*sonyflake.Sonyflake, 0, define.SnowFlake_End)
}

// snow flakes machine_id: 10 bits machineID + 6 bits plugin_type
func InitMachineID(machineID int16, startTime int64, cb func()) {
	sfs.cb = cb
	sfs.once.Do(func() {
		for n := 0; n < define.SnowFlake_End; n++ {
			var st sonyflake.Settings

			st.MachineID = func() (uint16, error) {
				newID := uint16(machineID<<6) + uint16(n)
				return newID, nil
			}

			st.CheckMachineID = func(id uint16) bool {
				return id <= (1<<16 - 1)
			}

			st.StartTime = time.Unix(startTime, 0)

			sf := sonyflake.NewSonyflake(st)
			if sf == nil {
				log.Panic().Str("start_time", st.StartTime.String()).Msg("sonyflake not created")
			}

			sfs.ids = append(sfs.ids, sf)
		}
	})
}

func NextID(tp int) (int64, error) {
	if tp < 0 || tp >= define.SnowFlake_End {
		return -1, fmt.Errorf("wrong id generated, type:%d", tp)
	}

	id, err := sfs.ids[tp].NextID()
	if err == nil {
		sfs.cb()
	}

	return int64(id), err
}

func MachineID(id int64) int16 {
	m := sonyflake.Decompose(uint64(id))
	return int16(m["machine-id"])
}

func MachineIDHigh(id int64) int16 {
	return MachineID(id) >> 6
}

func MachineIDLow(id int64) int16 {
	return MachineID(id) & 15
}
