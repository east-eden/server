package player

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/willf/bitset"
)

var (
	Guide_Max_Num uint = 256 // 引导数量上限

	ErrGuideInvalidIndex = errors.New("invalid guide index")
)

type GuideManager struct {
	owner     *Player        `bson:"-" json:"-"`
	guideBits *bitset.BitSet `bson:"-" json:"-"`
	GuideData []uint64       `bson:"guide_data" json:"guide_data"`
}

func NewGuideManager(owner *Player) *GuideManager {
	m := &GuideManager{
		owner:     owner,
		guideBits: bitset.New(Guide_Max_Num),
		GuideData: make([]uint64, 0, 4),
	}

	return m
}

func (m *GuideManager) AfterLoad() {
	loadBits := bitset.From(m.GuideData)
	m.guideBits = m.guideBits.Union(loadBits)
}

func (m *GuideManager) GuidePass(idx int32) error {
	if !utils.Between(idx, 0, int32(Guide_Max_Num)) {
		return ErrGuideInvalidIndex
	}

	m.guideBits.Set(uint(idx))
	m.GuideData = m.guideBits.Bytes()

	fields := map[string]any{
		"guide_data": m.GuideData,
	}

	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when GuideManager.GuidePass", m.owner.ID, fields)
	return err
}

func (m *GuideManager) GenGuideInfoPB() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, m.guideBits.Bytes())
	return buf.Bytes()
}
