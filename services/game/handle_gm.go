package game

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode"

	"bitbucket.org/funplus/server/define"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	pbMail "bitbucket.org/funplus/server/proto/server/mail"
	pbPubSub "bitbucket.org/funplus/server/proto/server/pubsub"
	"bitbucket.org/funplus/server/services/game/player"
	"bitbucket.org/funplus/server/transport"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	ErrInvalidGmCmd       = errors.New("invalid gm cmd")
	ErrPrivilegeNotEnough = errors.New("privilege not enough")
	registerCmds          = map[string]func(*player.Account, *MsgRegister, []string) error{
		"player": handleGmPlayer,
		"hero":   handleGmHero,
		"item":   handleGmItem,
		"token":  handleGmToken,
		"stage":  handleGmStage,
		"pub":    handleGmPub,
		"mail":   handleGmMail,
	}
)

func (r *MsgRegister) handleGmCmd(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_GmCmd)
	if !ok {
		return errors.New("handleGmCmd failed: recv message body error")
	}

	_, err := r.am.GetPlayerByAccount(acct)
	if err != nil {
		return err
	}

	return gmCmd(acct, r, msg.Cmd)
}

// 玩家相关gm命令
func handleGmPlayer(acct *player.Account, r *MsgRegister, cmds []string) error {
	switch strings.ToLower(cmds[0]) {
	case "level":
		change, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmPlayer failed", cmds, acct.ID) {
			return err
		}

		acct.GetPlayer().GmChangeLevel(int32(change))

	case "exp":
		change, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmPlayer failed", cmds, acct.ID) {
			return err
		}

		acct.GetPlayer().ChangeExp(int32(change))

	case "vip":
		change, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmPlayer failed", cmds, acct.ID) {
			return err
		}

		acct.GetPlayer().GmChangeVipLevel(int32(change))
	}

	return nil
}

// 英雄相关gm命令
func handleGmHero(acct *player.Account, r *MsgRegister, cmds []string) error {
	switch strings.ToLower(cmds[0]) {

	// 添加
	case "add":
		typeId, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmHero failed", cmds, acct.ID) {
			return err
		}

		acct.GetPlayer().HeroManager().AddHeroByTypeId(int32(typeId))

	// 经验改变
	case "exp":
		typeId, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmHero failed", cmds, acct.ID) {
			return err
		}

		exp, err := strconv.Atoi(cmds[2])
		if !utils.ErrCheck(err, "handleGmHero failed", cmds, acct.ID) {
			return err
		}

		h := acct.GetPlayer().HeroManager().GetHeroByTypeId(int32(typeId))
		if h == nil {
			return player.ErrHeroNotFound
		}

		return acct.GetPlayer().HeroManager().GmExpChange(h.Id, int32(exp))

	// 等级改变
	case "level":
		typeId, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmHero failed", cmds, acct.ID) {
			return err
		}

		level, err := strconv.Atoi(cmds[2])
		if !utils.ErrCheck(err, "handleGmHero failed", cmds, acct.ID) {
			return err
		}

		h := acct.GetPlayer().HeroManager().GetHeroByTypeId(int32(typeId))
		if h == nil {
			return player.ErrHeroNotFound
		}

		return acct.GetPlayer().HeroManager().GmLevelChange(h.Id, int32(level))

	// 突破
	case "promote":
		typeId, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmHero failed", cmds, acct.ID) {
			return err
		}

		promote, err := strconv.Atoi(cmds[2])
		if !utils.ErrCheck(err, "handleGmHero failed", cmds, acct.ID) {
			return err
		}

		h := acct.GetPlayer().HeroManager().GetHeroByTypeId(int32(typeId))
		if h == nil {
			return player.ErrHeroNotFound
		}

		return acct.GetPlayer().HeroManager().GmPromoteChange(h.Id, int32(promote))
	}

	return nil
}

// 物品相关gm命令
func handleGmItem(acct *player.Account, r *MsgRegister, cmds []string) error {
	switch strings.ToLower(cmds[0]) {

	// 添加
	case "add":
		typeId, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmItem failed", cmds, acct.ID) {
			return err
		}

		num := 1
		if len(cmds) >= 3 {
			num, err = strconv.Atoi(cmds[2])
			if !utils.ErrCheck(err, "handleGmItem failed", cmds, acct.ID) {
				return err
			}
		}

		return acct.GetPlayer().ItemManager().GainLoot(int32(typeId), int32(num))

	// 删除
	case "delete":
		fallthrough
	case "del":
		typeId, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmItem failed", cmds, acct.ID) {
			return err
		}

		num := 1
		if len(cmds) >= 3 {
			num, err = strconv.Atoi(cmds[2])
			if !utils.ErrCheck(err, "handleGmItem failed", cmds, acct.ID) {
				return err
			}
		}

		return acct.GetPlayer().ItemManager().DoCost(int32(typeId), int32(num))
	}
	return nil
}

// 代币相关gm命令
func handleGmToken(acct *player.Account, r *MsgRegister, cmds []string) error {
	switch strings.ToLower(cmds[0]) {
	case "add":
		tp, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmToken failed", cmds, acct.ID) {
			return err
		}

		add := 1000
		if len(cmds) >= 3 {
			add, err = strconv.Atoi(cmds[2])
			if !utils.ErrCheck(err, "handleGmToken failed", cmds, acct.ID) {
				return err
			}
		}

		return acct.GetPlayer().TokenManager().GainLoot(int32(tp), int32(add))
	}

	return nil
}

// 关卡相关gm命令
func handleGmStage(acct *player.Account, r *MsgRegister, cmds []string) error {
	switch strings.ToLower(cmds[0]) {
	case "pass":
		stageId, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmStage failed", cmds, acct.ID) {
			return err
		}

		return acct.GetPlayer().ChapterStageManager.StagePass(int32(stageId), []bool{true, true, true})
	}

	return nil
}

func handleGmPub(acct *player.Account, r *MsgRegister, cmds []string) error {
	switch cmds[0] {
	case "multi_publish_test":
		id, err := strconv.Atoi(cmds[1])
		if !utils.ErrCheck(err, "handleGmPub failed", cmds, acct.ID) {
			return err
		}

		name := cmds[2]

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = r.pubSub.PubMultiPublishTest(ctx, &pbPubSub.MultiPublishTest{
			Id:   int32(id),
			Name: name,
		})
		utils.ErrPrint(err, "PubMultiPublishTest failed when handleGmPub")

	case "game.StartGate":
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := r.pubSub.PubStartGate(ctx, &pbGlobal.AccountInfo{
			Id:    999,
			Name:  "StartGate Name",
			Level: 99,
		})
		utils.ErrPrint(err, "PubStartGate failed when handleGmPub")

	case "game.SyncPlayerInfo":
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := r.pubSub.PubSyncPlayerInfo(ctx, &player.PlayerInfo{
			ID:        99,
			AccountID: 999,
			Name:      "SyncPlayerInfo Name",
			Exp:       9999,
			Level:     9,
		})
		utils.ErrPrint(err, "PubSyncPlayerInfo failed when handleGmPub")
	}
	return nil
}

func handleGmMail(acct *player.Account, r *MsgRegister, cmds []string) error {
	switch cmds[0] {
	case "create":
		var title string
		var content string
		if len(cmds) >= 2 {
			title = cmds[1]
		}

		if len(cmds) >= 3 {
			content = cmds[2]
		}

		req := &pbMail.CreateMailRq{
			ReceiverId:  acct.GetPlayer().ID,
			SenderId:    1237475,
			Type:        pbGlobal.MailType_System,
			SenderName:  "来自深渊",
			Title:       title,
			Content:     content,
			Attachments: make([]*pbGlobal.LootData, 0),
		}

		req.Attachments = append(
			req.Attachments,
			&pbGlobal.LootData{
				Type: pbGlobal.LootType(define.CostLoot_Item),
				Misc: 1,
				Num:  2,
			},
			&pbGlobal.LootData{
				Type: pbGlobal.LootType(define.CostLoot_Token),
				Misc: 1,
				Num:  99,
			},
		)

		rsp, err := r.rpcHandler.CallCreateMail(req)
		if !utils.ErrCheck(err, "rpc call CreateSystemMail failed", req) {
			return err
		} else {
			log.Info().Interface("response", rsp).Msg("rpc call CreateSystemMail success")
		}
	case "read":
		_ = acct.GetPlayer().MailManager().ReadAllMail()
	case "gain":
		_ = acct.GetPlayer().MailManager().GainAllMailsAttachments()
	case "del":
		_ = acct.GetPlayer().MailManager().DelAllMails()
	}
	return nil
}

func gmCmd(acct *player.Account, r *MsgRegister, cmd string) error {
	// 权限不够
	if acct.Privilege <= 0 {
		return ErrPrivilegeNotEnough
	}

	if len(cmd) == 0 {
		return ErrInvalidGmCmd
	}
	cmds := strings.FieldsFunc(cmd, unicode.IsSpace)

	fn, ok := registerCmds[strings.ToLower(cmds[1])]
	if !ok {
		return ErrInvalidGmCmd
	}

	return fn(acct, r, cmds[2:])
}
