package iface

import (
	pbCombat "e.coding.net/mmstudio/blade/server/proto/server/combat"
	pbMail "e.coding.net/mmstudio/blade/server/proto/server/mail"
)

type RpcCaller interface {
	// 邮箱相关
	CallCreateMail(*pbMail.CreateMailRq) (*pbMail.CreateMailRs, error)
	CallQueryPlayerMails(*pbMail.QueryPlayerMailsRq) (*pbMail.QueryPlayerMailsRs, error)
	CallReadMail(*pbMail.ReadMailRq) (*pbMail.ReadMailRs, error)
	CallGainAttachments(*pbMail.GainAttachmentsRq) (*pbMail.GainAttachmentsRs, error)
	CallDelMail(*pbMail.DelMailRq) (*pbMail.DelMailRs, error)

	// 战斗相关
	CallStageCombat(*pbCombat.StageCombatRq) (*pbCombat.StageCombatRs, error)
}
