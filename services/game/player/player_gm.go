package player

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"bitbucket.org/funplus/server/utils"
)

var (
	ErrInvalidGmCmd       = errors.New("invalid gm cmd")
	ErrPrivilegeNotEnough = errors.New("privilege not enough")
	registerCmds          = map[string]func(*Player, []string) error{
		"player": handleGmPlayer,
		"hero":   handleGmHero,
		"item":   handleGmItem,
		"token":  handleGmToken,
	}
)

// 玩家相关gm命令
func handleGmPlayer(p *Player, cmds []string) error {
	switch strings.ToLower(cmds[0]) {
	case "level":
		change, err := strconv.Atoi(cmds[1])
		if pass := utils.ErrCheck(err, "handleGmPlayer failed", cmds, p.ID); !pass {
			return err
		}

		p.ChangeLevel(int32(change))

	case "exp":
		change, err := strconv.Atoi(cmds[1])
		if pass := utils.ErrCheck(err, "handleGmPlayer failed", cmds, p.ID); !pass {
			return err
		}

		p.ChangeExp(int32(change))
	}

	return nil
}

// 英雄相关gm命令
func handleGmHero(p *Player, cmds []string) error {
	switch strings.ToLower(cmds[0]) {
	case "add":
		typeId, err := strconv.Atoi(cmds[1])
		if pass := utils.ErrCheck(err, "handleGmHero failed", cmds, p.ID); !pass {
			return err
		}

		p.HeroManager().AddHeroByTypeId(int32(typeId))

	case "exp":
		typeId, err := strconv.Atoi(cmds[1])
		if pass := utils.ErrCheck(err, "handleGmHero failed", cmds, p.ID); !pass {
			return err
		}

		exp, err := strconv.Atoi(cmds[2])
		if pass := utils.ErrCheck(err, "handleGmHero failed", cmds, p.ID); !pass {
			return err
		}

		h := p.HeroManager().GetHeroByTypeId(int32(typeId))
		if h == nil {
			return ErrHeroNotFound
		}

		return p.HeroManager().GmExpChange(h.Id, int32(exp))

	case "level":
		typeId, err := strconv.Atoi(cmds[1])
		if pass := utils.ErrCheck(err, "handleGmHero failed", cmds, p.ID); !pass {
			return err
		}

		level, err := strconv.Atoi(cmds[2])
		if pass := utils.ErrCheck(err, "handleGmHero failed", cmds, p.ID); !pass {
			return err
		}

		h := p.HeroManager().GetHeroByTypeId(int32(typeId))
		if h == nil {
			return ErrHeroNotFound
		}

		return p.HeroManager().GmLevelChange(h.Id, int32(level))
	}

	return nil
}

// 物品相关gm命令
func handleGmItem(p *Player, cmds []string) error {
	switch strings.ToLower(cmds[0]) {
	case "add":
		typeId, err := strconv.Atoi(cmds[1])
		if pass := utils.ErrCheck(err, "handleGmItem failed", cmds, p.ID); !pass {
			return err
		}

		num := 1
		if len(cmds) >= 3 {
			num, err = strconv.Atoi(cmds[2])
			if pass := utils.ErrCheck(err, "handleGmItem failed", cmds, p.ID); !pass {
				return err
			}
		}

		return p.ItemManager().GainLoot(int32(typeId), int32(num))

	case "del":
		typeId, err := strconv.Atoi(cmds[1])
		if pass := utils.ErrCheck(err, "handleGmItem failed", cmds, p.ID); !pass {
			return err
		}

		num := 1
		if len(cmds) >= 3 {
			num, err = strconv.Atoi(cmds[2])
			if pass := utils.ErrCheck(err, "handleGmItem failed", cmds, p.ID); !pass {
				return err
			}
		}

		return p.ItemManager().DoCost(int32(typeId), int32(num))
	}
	return nil
}

// 代币相关gm命令
func handleGmToken(p *Player, cmds []string) error {
	switch strings.ToLower(cmds[0]) {
	case "add":
		tp, err := strconv.Atoi(cmds[1])
		if pass := utils.ErrCheck(err, "handleGmToken failed", cmds, p.ID); !pass {
			return err
		}

		add := 1000
		if len(cmds) >= 3 {
			add, err = strconv.Atoi(cmds[2])
			if pass := utils.ErrCheck(err, "handleGmToken failed", cmds, p.ID); !pass {
				return err
			}
		}

		return p.TokenManager().GainLoot(int32(tp), int32(add))
	}

	return nil
}

func GmCmd(p *Player, cmd string) error {
	// 权限不够
	if p.acct.Privilege <= 0 {
		return ErrPrivilegeNotEnough
	}

	if len(cmd) == 0 {
		return ErrInvalidGmCmd
	}
	cmds := strings.FieldsFunc(cmd, unicode.IsSpace)
	return handleGmCmds(p, cmds[1:])
}

func handleGmCmds(p *Player, cmds []string) error {
	fn, ok := registerCmds[strings.ToLower(cmds[0])]
	if !ok {
		return ErrInvalidGmCmd
	}

	return fn(p, cmds[1:])
}
