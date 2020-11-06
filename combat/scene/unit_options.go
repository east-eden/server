package scene

import (
	"log"
	"strconv"
	"strings"

	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/internal/att"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type UnitOption func(*UnitOptions)
type UnitOptions struct {
	TypeId     int32
	AttValue   []int64
	Position   [3]float32
	Entry      *define.UnitEntry
	AttManager *att.AttManager
}

func DefaultUnitOptions() *UnitOptions {
	return &UnitOptions{
		TypeId:   -1,
		Position: [3]float32{0, 0, 0},
		Entry:    nil,
	}
}

func WithUnitTypeId(typeId int32) UnitOption {
	return func(o *UnitOptions) {
		o.TypeId = typeId
	}
}

func WithUnitEntry(entry *define.UnitEntry) UnitOption {
	return func(o *UnitOptions) {
		o.Entry = entry
	}
}

func WithUnitAttValue(value []int64) UnitOption {
	return func(o *UnitOptions) {
		o.AttValue = value
	}
}

func WithUnitAttList(attList []*pbGame.Att) UnitOption {
	return func(o *UnitOptions) {
		o.AttValue = make([]int64, define.Att_End)

		for _, v := range attList {
			o.AttValue[v.AttType] = v.AttValue
		}
	}
}

func WithUnitPosition(pos []float32) UnitOption {
	return func(o *UnitOptions) {
		if len(pos) != 3 {
			return
		}

		for k, v := range pos {
			o.Position[k] = v
		}
	}
}

func WithUnitPositionString(pos string) UnitOption {
	return func(o *UnitOptions) {
		arr := strings.Split(pos, ",")
		if len(arr) != 3 {
			return
		}

		for k, v := range arr {
			p, err := strconv.ParseFloat(v, 32)
			if err != nil {
				log.Println("parsing position string error:", err)
				return
			}

			o.Position[k] = float32(p)
		}
	}
}

func WithAttManager(attId int32) UnitOption {
	return func(o *UnitOptions) {
		o.AttManager = att.NewAttManager(attId)
	}
}
