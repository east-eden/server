package att

import (
	"github.com/yokaiio/yokai_server/internal/define"
)

type AttManager struct {
	Owner  define.PluginObj
	Atts   [define.Att_End]int64   //一级属性
	AttExs [define.AttEx_End]int64 //二级属性
}

func NewAttManager(owner define.PluginObj, atts [define.Att_End]int64, att_exs [define.AttEx_End]int64) *AttManager {
	m := &AttManager{
		Owner: owner,
	}

	for k, v := range atts {
		m.Atts[k] = v
	}

	for k, v := range att_exs {
		m.AttExs[k] = v
	}

	return m
}

func (m *AttManager) GetAttValue(index int32) int64 {
	if index < 0 || index >= define.Att_End {
		return 0
	}

	return m.Atts[index]
}

func (m *AttManager) GetAttExValue(index int32) int64 {
	if index < 0 || index >= define.AttEx_End {
		return 0
	}

	return m.AttExs[index]
}

func (m *AttManager) CalcAtt() {
	// todo calc second att value by first att value
}
