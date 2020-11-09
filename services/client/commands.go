package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
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
		pages: make(map[int]*CommandPage, 0),
		c:     c,
	}

	cmder.initCommandPages()
	cmder.initCommands()

	return cmder
}

func reflectIntoMsg(msg proto.Message, result []string) error {
	// trans input into cmd.Message
	tp := reflect.TypeOf(msg).Elem()
	value := reflect.ValueOf(msg).Elem()

	// proto.Message struct has 3 invalid field
	if value.NumField()-3 != len(result) {
		return fmt.Errorf("输入数据无效")
	}

	// reflect into proto.Message
	for n := 0; n < len(result); n++ {
		ft := tp.Field(n).Type
		fv := value.Field(n)

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

func (cmd *Commander) CmdAccountLogon(ctx context.Context, result []string) (bool, string) {
	header := map[string]string{
		"Content-Type": "application/json",
	}

	var req struct {
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
	}

	req.UserID = result[0]
	req.UserName = result[1]

	body, err := json.Marshal(req)
	if err != nil {
		log.Warn().Err(err).Msg("json marshal failed when call CmdAccountLogon")
		return false, ""
	}

	resp, err := httpPost(cmd.c.transport.GetGateEndPoints(), header, body)
	if err != nil {
		log.Warn().Err(err).Msg("http post failed when call CmdAccountLogon")
		return false, ""
	}

	var gameInfo GameInfo
	if err := json.Unmarshal(resp, &gameInfo); err != nil {
		log.Warn().Err(err).Msg("json unmarshal failed when call CmdAccountLogon")
		return false, ""
	}

	log.Info().Interface("info", gameInfo).Msg("metadata unmarshaled result")

	if len(gameInfo.PublicTcpAddr) == 0 {
		log.Warn().Msg("invalid game public tcp address")
		return false, ""
	}

	cmd.c.transport.SetGameInfo(&gameInfo)
	cmd.c.transport.SetProtocol("tcp")
	if err := cmd.c.transport.StartConnect(ctx); err != nil {
		log.Warn().Err(err).Msg("tcp connect failed")
	}

	return true, "yokai_account.M2C_AccountLogon"
}

func (cmd *Commander) CmdWebSocketAccountLogon(ctx context.Context, result []string) (bool, string) {
	header := map[string]string{
		"Content-Type": "application/json",
	}

	var req struct {
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
	}

	req.UserID = result[0]
	req.UserName = result[1]

	body, err := json.Marshal(req)
	if err != nil {
		log.Warn().Err(err).Msg("json marshal failed when call CmdWebSocketAccountLogon")
		return false, ""
	}

	resp, err := httpPost(cmd.c.transport.GetGateEndPoints(), header, body)
	if err != nil {
		log.Warn().Err(err).Msg("http post failed when call CmdAccountLogon")
		return false, ""
	}

	var gameInfo GameInfo
	if err := json.Unmarshal(resp, &gameInfo); err != nil {
		log.Warn().Err(err).Msg("json unmarshal failed when call CmdAccountLogon")
		return false, ""
	}

	log.Info().Interface("info", gameInfo).Msg("metadata unmarshaled result")

	if len(gameInfo.PublicWsAddr) == 0 {
		log.Warn().Msg("invalid game public tcp address")
		return false, ""
	}

	cmd.c.transport.SetGameInfo(&gameInfo)
	cmd.c.transport.SetProtocol("ws")
	if err := cmd.c.transport.StartConnect(ctx); err != nil {
		log.Warn().Err(err).Msg("ws connect failed")
	}
	return true, "yokai_account.M2C_AccountLogon"
}

func (cmd *Commander) CmdCreatePlayer(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_CreatePlayer",
		Body: &pbGame.C2M_CreatePlayer{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdCreatePlayer command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_CreatePlayer"
}

func (cmd *Commander) CmdSendHeartBeat(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_account.C2M_HeartBeat",
		Body: &pbAccount.C2M_HeartBeat{},
	}

	cmd.c.transport.SendMessage(msg)

	return false, ""
}

func (cmd *Commander) CmdCliAccountDisconnect(ctx context.Context, result []string) (bool, string) {
	cmd.c.transport.StartDisconnect()
	return false, ""
}

func (cmd *Commander) CmdServerAccountDisconnect(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_account.C2M_AccountDisconnect",
		Body: &pbAccount.C2M_AccountDisconnect{},
	}

	cmd.c.transport.SendMessage(msg)

	return false, ""
}

func (cmd *Commander) CmdQueryPlayerInfo(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_QueryPlayerInfo",
		Body: &pbGame.C2M_QueryPlayerInfo{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_QueryPlayerInfo"
}

func (cmd *Commander) CmdChangeExp(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_ChangeExp",
		Body: &pbGame.C2M_ChangeExp{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdChangeExp command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_ExpUpdate"
}

func (cmd *Commander) CmdChangeLevel(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_ChangeLevel",
		Body: &pbGame.C2M_ChangeLevel{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdChangeLevel command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_ExpUpdate"
}

func (cmd *Commander) CmdSyncPlayerInfo(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_SyncPlayerInfo",
		Body: &pbGame.C2M_SyncPlayerInfo{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_SyncPlayerInfo"
}

func (cmd *Commander) CmdPublicSyncPlayerInfo(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_PublicSyncPlayerInfo",
		Body: &pbGame.C2M_PublicSyncPlayerInfo{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_PublicSyncPlayerInfo"
}

func (cmd *Commander) CmdQueryHeros(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_QueryHeros",
		Body: &pbGame.C2M_QueryHeros{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_HeroList"
}

func (cmd *Commander) CmdAddHero(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_AddHero",
		Body: &pbGame.C2M_AddHero{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdAddHero command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_HeroList"
}

func (cmd *Commander) CmdDelHero(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_DelHero",
		Body: &pbGame.C2M_DelHero{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdDelHero command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_HeroList"
}

func (cmd *Commander) CmdQueryItems(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_QueryItems",
		Body: &pbGame.C2M_QueryItems{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_ItemList"
}

func (cmd *Commander) CmdAddItem(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_AddItem",
		Body: &pbGame.C2M_AddItem{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdAddItem command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_ItemUpdate,yokai_game.M2C_ItemAdd"
}

func (cmd *Commander) CmdDelItem(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_DelItem",
		Body: &pbGame.C2M_DelItem{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdDelItem command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_DelItem"
}

func (cmd *Commander) CmdUseItem(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_UseItem",
		Body: &pbGame.C2M_UseItem{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdUseItem command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_DelItem,yokai_game.M2C_ItemUpdate"
}

func (cmd *Commander) CmdHeroPutonEquip(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_PutonEquip",
		Body: &pbGame.C2M_PutonEquip{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdHeroPutonEquip command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_HeroInfo"
}

func (cmd *Commander) CmdHeroTakeoffEquip(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_TakeoffEquip",
		Body: &pbGame.C2M_TakeoffEquip{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdHeroTakeoffEquip command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_HeroInfo"
}

func (cmd *Commander) CmdQueryTokens(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_QueryTokens",
		Body: &pbGame.C2M_QueryTokens{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdQueryTokens command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_TokenList"
}

func (cmd *Commander) CmdAddToken(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_AddToken",
		Body: &pbGame.C2M_AddToken{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdAddToken command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_TokenList"
}

func (cmd *Commander) CmdQueryTalents(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_QueryTalents",
		Body: &pbGame.C2M_QueryTalents{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdQueryTalents command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_TalentList"
}

func (cmd *Commander) CmdAddTalent(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_AddTalent",
		Body: &pbGame.C2M_AddTalent{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdAddTalent command failed:", err)
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_TalentList"
}

func (cmd *Commander) CmdStartStageCombat(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_StartStageCombat",
		Body: &pbGame.C2M_StartStageCombat{RpcId: 1},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "yokai_game.M2C_StartStageCombat"
}

func (c *Commander) registerCommand(cmd *Command) {
	cmdPage, ok := c.pages[cmd.PageID]
	if !ok {
		fmt.Println("register command failed:", cmd)
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

	// 9退出
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

	// 3客户端断开连接
	c.registerCommand(&Command{Text: "客户端断开连接", PageID: 2, GotoPageID: -1, Cb: c.CmdCliAccountDisconnect})

	// 4服务器断开连接
	c.registerCommand(&Command{Text: "服务器断开连接", PageID: 2, GotoPageID: -1, Cb: c.CmdServerAccountDisconnect})

	///////////////////////////////////////////////
	// 角色管理
	///////////////////////////////////////////////
	// 返回上页
	c.registerCommand(&Command{Text: "返回上页", PageID: 3, GotoPageID: 1, Cb: nil})

	// 1查询账号下所有角色
	c.registerCommand(&Command{Text: "查询账号下所有角色", PageID: 3, GotoPageID: -1, Cb: c.CmdQueryPlayerInfo})

	// 2创建角色
	c.registerCommand(&Command{Text: "创建角色", PageID: 3, GotoPageID: -1, InputText: "请输入rpcid和角色名字:", DefaultInput: "1,加百列", Cb: c.CmdCreatePlayer})

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

}
