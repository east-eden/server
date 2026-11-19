package game

import (
	"context"
	"errors"
	"time"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/alliance"
	"github.com/east-eden/server/services/game/player"
)

// sendAllianceInfo 下发公会信息
func sendAllianceInfo(pl *player.Player) {
	g, err := pl.AllianceManager().Info()
	if err != nil || g == nil {
		return
	}

	info := &pbGlobal.AllianceInfo{
		Id:       g.Id,
		Name:     g.Name,
		LeaderId: g.LeaderId,
		Level:    g.Level,
		Exp:      g.Exp,
		Notice:   g.Notice,
	}

	for _, id := range pl.AllianceManager().MemberList(g) {
		info.Members = append(info.Members, &pbGlobal.AllianceMember{Id: id, Position: g.Members[id]})
	}

	pl.SendProtoMessage(&pbGlobal.S2C_AllianceInfo{Info: info})
}

func (m *MsgRegister) handleAllianceCreate(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_AllianceCreate)
	if !ok {
		return errors.New("handleAllianceCreate failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	if _, err := pl.AllianceManager().Create(msg.GetName()); err != nil {
		return err
	}

	sendAllianceInfo(pl)
	return nil
}

func (m *MsgRegister) handleAllianceJoin(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_AllianceJoin)
	if !ok {
		return errors.New("handleAllianceJoin failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	if err := pl.AllianceManager().Join(msg.GetId()); err != nil {
		return err
	}

	sendAllianceInfo(pl)
	return nil
}

func (m *MsgRegister) handleAllianceLeave(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	if err := pl.AllianceManager().Leave(); err != nil {
		return err
	}

	pl.SendProtoMessage(&pbGlobal.S2C_AllianceInfo{})
	return nil
}

func (m *MsgRegister) handleAllianceKick(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_AllianceKick)
	if !ok {
		return errors.New("handleAllianceKick failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	if err := pl.AllianceManager().Kick(msg.GetMemberId()); err != nil {
		return err
	}

	sendAllianceInfo(pl)
	return nil
}

func (m *MsgRegister) handleAllianceBuild(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	if err := pl.AllianceManager().Build(timeNow()); err != nil {
		return err
	}

	sendAllianceInfo(pl)
	return nil
}

var timeNow = func() time.Time { return time.Now() }

var _ = alliance.PositionLeader
