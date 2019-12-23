package utils

import (
	"fmt"
	"sync"

	"github.com/sony/sonyflake"
	"github.com/yokaiio/yokai_server/internal/define"
)

type Snowflakes struct {
	ids  []*sonyflake.Sonyflake
	once sync.Once
}

var sfs Snowflakes

func init() {
	sfs.ids = make([]*sonyflake.Sonyflake, 0, define.Plugin_End)
}

func InitMachineID(machineID uint16) {
	sfs.once.Do(func() {
		for n := 0; n < define.Plugin_End; n++ {
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
	if tp < 0 || tp >= define.Plugin_End {
		return -1, fmt.Errorf("wrong id generated, type:%d", tp)
	}

	id, err := sfs.ids[tp].NextID()
	return int64(id), err
}
