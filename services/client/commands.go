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
	cmder.initCommands()

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

	// first level page
	// page main options
	c.registerCommandPage(&CommandPage{PageID: Cmd_Page_Main, ParentPageID: -1, Cmds: make([]*Command, 0)})

	// seconde level page
	// page server connection options
	c.registerCommandPage(&CommandPage{PageID: Cmd_Page_Server, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// page role options
	c.registerCommandPage(&CommandPage{PageID: Cmd_Page_Role, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// page hero options
	c.registerCommandPage(&CommandPage{PageID: Cmd_Page_Hero, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// page item options
	c.registerCommandPage(&CommandPage{PageID: Cmd_Page_Item, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// page equip options
	c.registerCommandPage(&CommandPage{PageID: Cmd_Page_Equip, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// page token options
	c.registerCommandPage(&CommandPage{PageID: Cmd_Page_Token, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// page combat options
	c.registerCommandPage(&CommandPage{PageID: Cmd_Page_Combat, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// page fragment options
	c.registerCommandPage(&CommandPage{PageID: Cmd_Page_Fragment, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// page crystal options
	c.registerCommandPage(&CommandPage{PageID: Cmd_Page_Crystal, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})
}

func (c *Commander) initCommands() {
	// first level page
	// 0服务器连接管理
	c.registerCommand(&Command{Text: "服务器连接管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Server, Cb: nil})

	// 1角色管理
	c.registerCommand(&Command{Text: "角色管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Role, Cb: nil})

	// 2英雄管理
	c.registerCommand(&Command{Text: "英雄管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Hero, Cb: nil})

	// 3物品管理
	c.registerCommand(&Command{Text: "物品管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Item, Cb: nil})

	// 4装备管理
	c.registerCommand(&Command{Text: "装备管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Equip, Cb: nil})

	// 5代币管理
	c.registerCommand(&Command{Text: "代币管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Token, Cb: nil})

	// 7战斗管理
	c.registerCommand(&Command{Text: "战斗管理", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Combat, Cb: nil})

	// 8英雄碎片
	c.registerCommand(&Command{Text: "英雄碎片", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Fragment, Cb: nil})

	// 9晶石
	c.registerCommand(&Command{Text: "晶石", PageID: Cmd_Page_Main, GotoPageID: Cmd_Page_Crystal, Cb: nil})

	// 10退出
	c.registerCommand(&Command{Text: "退出", PageID: Cmd_Page_Main, GotoPageID: -1, Cb: c.CmdQuit})

	///////////////////////////////////////////////
	// 服务器连接管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Server, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1登录
	c.registerCommand(&Command{Text: "登录", PageID: Cmd_Page_Server, GotoPageID: -1, InputText: "请输入登录user ID和名字，以逗号分隔", DefaultInput: "1,dudu", Cb: c.CmdAccountLogon})

	// websocket连接登录
	c.registerCommand(&Command{Text: "websocket登录", PageID: Cmd_Page_Server, GotoPageID: -1, InputText: "请输入登录user ID和名字，以逗号分隔", DefaultInput: "1,dudu", Cb: c.CmdWebSocketAccountLogon})

	// 2发送心跳
	c.registerCommand(&Command{Text: "发送心跳", PageID: Cmd_Page_Server, GotoPageID: -1, Cb: c.CmdSendHeartBeat})

	// 3发送ClientMessage
	c.registerCommand(&Command{Text: "发送等待服务器返回消息", PageID: Cmd_Page_Server, GotoPageID: -1, Cb: c.CmdWaitResponseMessage})

	// 4客户端断开连接
	c.registerCommand(&Command{Text: "客户端断开连接", PageID: Cmd_Page_Server, GotoPageID: -1, Cb: c.CmdCliAccountDisconnect})

	// 5服务器断开连接
	c.registerCommand(&Command{Text: "服务器断开连接", PageID: Cmd_Page_Server, GotoPageID: -1, Cb: c.CmdServerAccountDisconnect})

	///////////////////////////////////////////////
	// 角色管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Role, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1查询账号下所有角色
	c.registerCommand(&Command{Text: "查询账号下所有角色", PageID: Cmd_Page_Role, GotoPageID: -1, Cb: c.CmdQueryPlayerInfo})

	// 2创建角色
	c.registerCommand(&Command{Text: "创建角色", PageID: Cmd_Page_Role, GotoPageID: -1, InputText: "请输入rpcid和角色名字:", DefaultInput: "加百列", Cb: c.CmdCreatePlayer})

	// 3gm命令
	c.registerCommand(&Command{Text: "gm命令", PageID: Cmd_Page_Role, GotoPageID: -1, InputText: "请输入gm命令", DefaultInput: "gm player exp 100", Cb: c.CmdGmCmd})

	///////////////////////////////////////////////
	// 英雄管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Hero, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1查询英雄信息
	c.registerCommand(&Command{Text: "查询英雄信息", PageID: Cmd_Page_Hero, GotoPageID: -1, Cb: c.CmdQueryHeros})

	// 3删除英雄
	c.registerCommand(&Command{Text: "删除英雄", PageID: Cmd_Page_Hero, GotoPageID: -1, InputText: "请输入要删除的英雄ID:", DefaultInput: "1", Cb: c.CmdDelHero})

	// 4查询英雄属性
	c.registerCommand(&Command{Text: "查询英雄属性", PageID: Cmd_Page_Hero, GotoPageID: -1, InputText: "请输入要查询的英雄ID:", DefaultInput: "1", Cb: c.CmdQueryHeroAtt})

	// 4增加经验
	//registerCommand(&Command{Text: "增加经验", PageID: 4, GotoPageID: -1, InputText: "请输入英雄id和经验，用逗号分隔:", DefaultInput: "1,110", Cb: CmdHeroAddExp})

	// 5增加等级
	//registerCommand(&Command{Text: "增加等级", PageID: 4, GotoPageID: -1, InputText: "请输入英雄id和等级，用逗号分隔:", DefaultInput: "1,3", Cb: CmdHeroAddLevel})

	///////////////////////////////////////////////
	// 物品管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Item, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1查询物品信息
	c.registerCommand(&Command{Text: "查询物品信息", PageID: Cmd_Page_Item, GotoPageID: -1, Cb: c.CmdQueryItems})

	// 3删除物品
	c.registerCommand(&Command{Text: "删除物品", PageID: Cmd_Page_Item, GotoPageID: -1, InputText: "请输入要删除的物品ID:", DefaultInput: "1", Cb: c.CmdDelItem})

	// 4使用物品
	c.registerCommand(&Command{Text: "使用物品", PageID: Cmd_Page_Item, GotoPageID: -1, InputText: "请输入要使用的物品ID:", Cb: c.CmdUseItem})

	///////////////////////////////////////////////
	// 装备管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Equip, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 2穿装备
	c.registerCommand(&Command{Text: "穿装备", PageID: Cmd_Page_Equip, GotoPageID: -1, InputText: "请输入英雄ID和物品ID:", DefaultInput: "1,1", Cb: c.CmdHeroPutonEquip})

	// 3脱装备
	c.registerCommand(&Command{Text: "脱装备", PageID: Cmd_Page_Equip, GotoPageID: -1, InputText: "请输入英雄ID和装备位置索引:", DefaultInput: "1,0", Cb: c.CmdHeroTakeoffEquip})

	// 4装备升级
	c.registerCommand(&Command{Text: "装备升级", PageID: Cmd_Page_Equip, GotoPageID: -1, InputText: "请输入装备ID:", Cb: c.CmdEquipLevelup})

	///////////////////////////////////////////////
	// 代币管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Token, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1查询代币信息
	c.registerCommand(&Command{Text: "查询代币信息", PageID: Cmd_Page_Token, GotoPageID: -1, Cb: c.CmdQueryTokens})

	///////////////////////////////////////////////
	// 战斗管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Combat, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1关卡战斗
	c.registerCommand(&Command{Text: "普通关卡战斗", PageID: Cmd_Page_Combat, GotoPageID: -1, Cb: c.CmdStartStageCombat})

	///////////////////////////////////////////////
	// 英雄碎片
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Fragment, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1请求碎片信息
	c.registerCommand(&Command{Text: "请求碎片信息", PageID: Cmd_Page_Fragment, GotoPageID: -1, Cb: c.CmdQueryFragments})

	// 2碎片合成
	c.registerCommand(&Command{Text: "碎片合成", PageID: Cmd_Page_Fragment, GotoPageID: -1, InputText: "请输入碎片ID:", DefaultInput: "1", Cb: c.CmdFragmentsCompose})

	///////////////////////////////////////////////
	// 晶石
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Crystal, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1晶石升级
	c.registerCommand(&Command{Text: "晶石升级", PageID: Cmd_Page_Crystal, GotoPageID: -1, InputText: "请输入晶石ID:", Cb: c.CmdCrystalLevelup})

	// 2装备晶石
	c.registerCommand(&Command{Text: "装备晶石", PageID: Cmd_Page_Crystal, GotoPageID: -1, InputText: "请输入英雄ID和晶石ID:", DefaultInput: "1,1", Cb: c.CmdPutonCrystal})

	// 3卸下晶石
	c.registerCommand(&Command{Text: "卸下晶石", PageID: Cmd_Page_Crystal, GotoPageID: -1, InputText: "请输入英雄ID和位置:", DefaultInput: "1,0", Cb: c.CmdTakeoffCrystal})
}
