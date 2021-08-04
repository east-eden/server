package iface

import (
	pbCombat "github.com/east-eden/server/proto/server/combat"
	pbMail "github.com/east-eden/server/proto/server/mail"
	pbRank "github.com/east-eden/server/proto/server/rank"
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

	// 排行相关
	CallQueryRankByObjId(*pbRank.QueryRankByObjIdRq) (*pbRank.QueryRankByObjIdRs, error)
	CallSetRankScore(*pbRank.SetRankScoreRq) (*pbRank.SetRankScoreRs, error)
}
