package client

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

const (
	Cmd_Page_Main     = 1  // 主选项
	Cmd_Page_Server   = 2  // 服务器选项
	Cmd_Page_Role     = 3  // 角色选项
	Cmd_Page_Hero     = 4  // 英雄选项
	Cmd_Page_Item     = 5  // 物品选项
	Cmd_Page_Equip    = 6  // 装备选项
	Cmd_Page_Token    = 7  // 代币选项
	Cmd_Page_Combat   = 8  // 战斗选项
	Cmd_Page_Fragment = 9  // 英雄碎片选项
	Cmd_Page_Crystal  = 10 // 晶石选项
)

type Command struct {
	Number       int
	Text         string
	PageID       int
	GotoPageID   int
	Cb           func(context.Context, []string) (bool, string)
	InputText    string
	DefaultInput string
}

type CommandPage struct {
	PageID       int
	ParentPageID int
	Cmds         []*Command
}

type Commander struct {
	pages map[int]*CommandPage
	c     *Client
}

func NewCommander(c *Client) *Commander {
	cmder := &Commander{
		pages: make(map[int]*CommandPage),
		c:     c,
	}

	cmder.initCommandPages()

	return cmder
}

func reflectIntoMsg(msg proto.Message, result []string) error {
	var fieldOffset = 3

	// trans input into cmd.Message
	tp := reflect.TypeOf(msg).Elem()
	value := reflect.ValueOf(msg).Elem()

	// proto.Message struct has 3 invalid field
	if value.NumField()-3 != len(result) {
		return fmt.Errorf("输入数据无效")
	}

	// reflect into proto.Message
	for n := 0; n < len(result); n++ {
		ft := tp.Field(n + fieldOffset).Type
		fv := value.Field(n + fieldOffset)

		if ft.Kind() >= reflect.Int && ft.Kind() <= reflect.Uint64 {
			inputValue, err := strconv.ParseInt(result[n], 10, ft.Bits())
			if err != nil {
				return fmt.Errorf("input value<%s> cannot assert to type<%s>\r\n", result[n], ft.Name())
			}

			input := reflect.ValueOf(inputValue).Convert(ft)
			fv.Set(input)
		}

		if ft.Kind() == reflect.String {
			fv.Set(reflect.ValueOf(result[n]))
		}
	}

	return nil
}

func (cmd *Commander) CmdQuit(ctx context.Context, result []string) (bool, string) {
	os.Exit(0)
	return false, ""
}

func (c *Commander) registerCommand(cmd *Command) {
	cmdPage, ok := c.pages[cmd.PageID]
	if !ok {
		log.Info().Interface("cmd", cmd).Msg("register command failed")
		return
	}

	cmdPage.Cmds = append(cmdPage.Cmds, cmd)
	cmd.Number = len(cmdPage.Cmds)
}

func (c *Commander) registerCommandPage(p *CommandPage) {
	c.pages[p.PageID] = p
}

func (c *Commander) initCommandPages() {
	c.initMainCommands()
	c.initServerCommands()
	c.initRoleCommands()
	c.initHeroCommands()
	c.initItemCommands()
	c.initEquipCommands()
	c.initTokenCommands()
	c.initCombatCommands()
	c.initFragmentCommands()
	c.initCrystalCommands()
}

func (cmd *Commander) initMainCommands() {
	// page main options
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Main, ParentPageID: -1, Cmds: make([]*Command, 0)})

	// 0服务器连接管理
	cmd.registerCommand(&Command{Text: "服务器连接管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Server, Cb: nil})

	// 1角色管理
	cmd.registerCommand(&Command{Text: "角色管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Role, Cb: nil})

	// 2英雄管理
	cmd.registerCommand(&Command{Text: "英雄管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Hero, Cb: nil})

	// 3物品管理
	cmd.registerCommand(&Command{Text: "物品管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Item, Cb: nil})

	// 4装备管理
	cmd.registerCommand(&Command{Text: "装备管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Equip, Cb: nil})

	// 5代币管理
	cmd.registerCommand(&Command{Text: "代币管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Token, Cb: nil})

	// 7战斗管理
	cmd.registerCommand(&Command{Text: "战斗管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Combat, Cb: nil})

	// 8英雄碎片
	cmd.registerCommand(&Command{Text: "英雄碎片", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Fragment, Cb: nil})

	// 9晶石
	cmd.registerCommand(&Command{Text: "晶石", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Crystal, Cb: nil})

	// 10退出
	cmd.registerCommand(&Command{Text: "退出", PageID: Cmd_Page_Main, GotoPageID: -1, Cb: cmd.CmdQuit})
}
