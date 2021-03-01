package utils

import (
	"fmt"
	"sync"

	"bitbucket.org/funplus/server/define"
	"github.com/sony/sonyflake"
)

type Snowflakes struct {
	ids  []*sonyflake.Sonyflake
	once sync.Once
}

var sfs Snowflakes

func init() {
	sfs.ids = make([]*sonyflake.Sonyflake, 0, define.SnowFlake_End)
}

// snow flakes machine_id: 10 bits machineID + 6 bits plugin_type
func InitMachineID(machineID int16) {
	sfs.once.Do(func() {
		for n := 0; n < define.SnowFlake_End; n++ {
			var st sonyflake.Settings

			st.MachineID = func() (uint16, error) {
				newID := uint16(machineID<<6) + uint16(n)
				return newID, nil
			}

			sf := sonyflake.NewSonyflake(st)
			if sf == nil {
				panic("sonyflake not created")
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
