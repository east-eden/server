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
	c.registerCommandPage(&CommandPage{PageID: 1, ParentPageID: -1, Cmds: make([]*Command, 0)})

	// seconde level page
	// page server connection options
	c.registerCommandPage(&CommandPage{PageID: 2, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page role options
	c.registerCommandPage(&CommandPage{PageID: 3, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page hero options
	c.registerCommandPage(&CommandPage{PageID: 4, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page item options
	c.registerCommandPage(&CommandPage{PageID: 5, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page equip options
	c.registerCommandPage(&CommandPage{PageID: 6, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page token options
	c.registerCommandPage(&CommandPage{PageID: 7, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page blade options
	c.registerCommandPage(&CommandPage{PageID: 8, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page combat options
	c.registerCommandPage(&CommandPage{PageID: 9, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page fragment options
	c.registerCommandPage(&CommandPage{PageID: 10, ParentPageID: 1, Cmds: make([]*Command, 0)})
}

func (c *Commander) initCommands() {
	// first level page
	// 0服务器连接管理
	c.registerCommand(&Command{Text: "服务器连接管理", PageID: 1, GotoPageID: 2, Cb: nil})

	// 1角色管理
	c.registerCommand(&Command{Text: "角色管理", PageID: 1, GotoPageID: 3, Cb: nil})

	// 2英雄管理
	c.registerCommand(&Command{Text: "英雄管理", PageID: 1, GotoPageID: 4, Cb: nil})

	// 3物品管理
	c.registerCommand(&Command{Text: "物品管理", PageID: 1, GotoPageID: 5, Cb: nil})

	// 4装备管理
	c.registerCommand(&Command{Text: "装备管理", PageID: 1, GotoPageID: 6, Cb: nil})

	// 5代币管理
	c.registerCommand(&Command{Text: "代币管理", PageID: 1, GotoPageID: 7, Cb: nil})

	// 6异刃管理
	c.registerCommand(&Command{Text: "异刃管理", PageID: 1, GotoPageID: 8, Cb: nil})

	// 7战斗管理
	c.registerCommand(&Command{Text: "战斗管理", PageID: 1, GotoPageID: 9, Cb: nil})

	// 9英雄碎片
	c.registerCommand(&Command{Text: "英雄碎片", PageID: 1, GotoPageID: 10, Cb: nil})

	// 10退出
	c.registerCommand(&Command{Text: "退出", PageID: 1, GotoPageID: -1, Cb: c.CmdQuit})

	///////////////////////////////////////////////
	// 服务器连接管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: 2, GotoPageID: 1, Cb: nil})

	// 1登录
	c.registerCommand(&Command{Text: "登录", PageID: 2, GotoPageID: -1, InputText: "请输入登录user ID和名字，以逗号分隔", DefaultInput: "1,dudu", Cb: c.CmdAccountLogon})

	// websocket连接登录
	c.registerCommand(&Command{Text: "websocket登录", PageID: 2, GotoPageID: -1, InputText: "请输入登录user ID和名字，以逗号分隔", DefaultInput: "1,dudu", Cb: c.CmdWebSocketAccountLogon})

	// 2发送心跳
	c.registerCommand(&Command{Text: "发送心跳", PageID: 2, GotoPageID: -1, Cb: c.CmdSendHeartBeat})

	// 3发送ClientMessage
	c.registerCommand(&Command{Text: "发送等待服务器返回消息", PageID: 2, GotoPageID: -1, Cb: c.CmdWaitResponseMessage})

	// 4客户端断开连接
	c.registerCommand(&Command{Text: "客户端断开连接", PageID: 2, GotoPageID: -1, Cb: c.CmdCliAccountDisconnect})

	// 5服务器断开连接
	c.registerCommand(&Command{Text: "服务器断开连接", PageID: 2, GotoPageID: -1, Cb: c.CmdServerAccountDisconnect})

	///////////////////////////////////////////////
	// 角色管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: 3, GotoPageID: 1, Cb: nil})

	// 1查询账号下所有角色
	c.registerCommand(&Command{Text: "查询账号下所有角色", PageID: 3, GotoPageID: -1, Cb: c.CmdQueryPlayerInfo})

	// 2创建角色
	c.registerCommand(&Command{Text: "创建角色", PageID: 3, GotoPageID: -1, InputText: "请输入rpcid和角色名字:", DefaultInput: "加百列", Cb: c.CmdCreatePlayer})

	// 3改变经验
	c.registerCommand(&Command{Text: "改变经验", PageID: 3, GotoPageID: -1, InputText: "请输入要改变的经验值:", DefaultInput: "120", Cb: c.CmdChangeExp})

	// 4改变等级
	c.registerCommand(&Command{Text: "改变等级", PageID: 3, GotoPageID: -1, InputText: "请输入要改变的等级:", DefaultInput: "10", Cb: c.CmdChangeLevel})

	// 5同步玩家信息到gate
	c.registerCommand(&Command{Text: "同步gate", PageID: 3, GotoPageID: -1, Cb: c.CmdSyncPlayerInfo})

	// 6publish玩家信息
	c.registerCommand(&Command{Text: "publish玩家信息", PageID: 3, GotoPageID: -1, Cb: c.CmdPublicSyncPlayerInfo})

	///////////////////////////////////////////////
	// 英雄管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: 4, GotoPageID: 1, Cb: nil})

	// 1查询英雄信息
	c.registerCommand(&Command{Text: "查询英雄信息", PageID: 4, GotoPageID: -1, Cb: c.CmdQueryHeros})

	// 2添加英雄
	c.registerCommand(&Command{Text: "添加英雄", PageID: 4, GotoPageID: -1, InputText: "请输入要添加的英雄TypeID:", DefaultInput: "1", Cb: c.CmdAddHero})

	// 3删除英雄
	c.registerCommand(&Command{Text: "删除英雄", PageID: 4, GotoPageID: -1, InputText: "请输入要删除的英雄ID:", DefaultInput: "1", Cb: c.CmdDelHero})

	// 4增加经验
	//registerCommand(&Command{Text: "增加经验", PageID: 4, GotoPageID: -1, InputText: "请输入英雄id和经验，用逗号分隔:", DefaultInput: "1,110", Cb: CmdHeroAddExp})

	// 5增加等级
	//registerCommand(&Command{Text: "增加等级", PageID: 4, GotoPageID: -1, InputText: "请输入英雄id和等级，用逗号分隔:", DefaultInput: "1,3", Cb: CmdHeroAddLevel})

	///////////////////////////////////////////////
	// 物品管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: 5, GotoPageID: 1, Cb: nil})

	// 1查询物品信息
	c.registerCommand(&Command{Text: "查询物品信息", PageID: 5, GotoPageID: -1, Cb: c.CmdQueryItems})

	// 2添加物品
	c.registerCommand(&Command{Text: "添加物品", PageID: 5, GotoPageID: -1, InputText: "请输入要添加的物品TypeID:", DefaultInput: "1", Cb: c.CmdAddItem})

	// 3删除物品
	c.registerCommand(&Command{Text: "删除物品", PageID: 5, GotoPageID: -1, InputText: "请输入要删除的物品ID:", DefaultInput: "1", Cb: c.CmdDelItem})

	// 4使用物品
	c.registerCommand(&Command{Text: "使用物品", PageID: 5, GotoPageID: -1, InputText: "请输入要使用的物品ID:", Cb: c.CmdUseItem})

	///////////////////////////////////////////////
	// 装备管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: 6, GotoPageID: 1, Cb: nil})

	// 2穿装备
	c.registerCommand(&Command{Text: "穿装备", PageID: 6, GotoPageID: -1, InputText: "请输入英雄ID和物品ID:", DefaultInput: "1,1", Cb: c.CmdHeroPutonEquip})

	// 3脱装备
	c.registerCommand(&Command{Text: "脱装备", PageID: 6, GotoPageID: -1, InputText: "请输入英雄ID和装备位置索引:", DefaultInput: "1,0", Cb: c.CmdHeroTakeoffEquip})

	// 4装备升级
	c.registerCommand(&Command{Text: "装备升级", PageID: 6, GotoPageID: -1, InputText: "请输入装备ID:", Cb: c.CmdEquipLevelup})

	///////////////////////////////////////////////
	// 代币管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: 7, GotoPageID: 1, Cb: nil})

	// 1查询代币信息
	c.registerCommand(&Command{Text: "查询代币信息", PageID: 7, GotoPageID: -1, Cb: c.CmdQueryTokens})

	// 2变更代币数量
	c.registerCommand(&Command{Text: "变更代币数量", PageID: 7, GotoPageID: -1, InputText: "请输入要变更的代币类型和数量，用逗号分隔:", DefaultInput: "0,1000", Cb: c.CmdAddToken})

	///////////////////////////////////////////////
	// 异刃管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: 8, GotoPageID: 1, Cb: nil})

	// 1查询天赋信息
	c.registerCommand(&Command{Text: "查询天赋信息", PageID: 8, GotoPageID: -1, InputText: "请输入异刃ID:", DefaultInput: "1", Cb: c.CmdQueryTalents})

	// 2增加天赋
	c.registerCommand(&Command{Text: "增加天赋", PageID: 8, GotoPageID: -1, InputText: "请输入异刃ID和天赋ID:", DefaultInput: "1,1", Cb: c.CmdAddTalent})

	///////////////////////////////////////////////
	// 战斗管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: 9, GotoPageID: 1, Cb: nil})

	// 1关卡战斗
	c.registerCommand(&Command{Text: "普通关卡战斗", PageID: 9, GotoPageID: -1, Cb: c.CmdStartStageCombat})

	///////////////////////////////////////////////
	// 英雄碎片
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: 10, GotoPageID: 1, Cb: nil})

	// 1请求碎片信息
	c.registerCommand(&Command{Text: "请求碎片信息", PageID: 10, GotoPageID: -1, Cb: c.CmdQueryFragments})

	// 2碎片合成
	c.registerCommand(&Command{Text: "碎片合成", PageID: 10, GotoPageID: -1, InputText: "请输入碎片ID:", DefaultInput: "1", Cb: c.CmdFragmentsCompose})
}
