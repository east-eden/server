package client

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type Command struct {
	Number       int
	Text         string
	PageID       int
	GotoPageID   int
	Cb           func(*TcpClient, []string) bool
	InputText    string
	DefaultInput string
}

type CommandPage struct {
	PageID       int
	ParentPageID int
	Cmds         []*Command
}

var (
	CmdPages = make(map[int]*CommandPage, 0)
)

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

func CmdQuit(c *TcpClient, result []string) bool {
	os.Exit(0)
	return false
	//syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

func CmdClientLogon(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_ClientLogon",
		Body: &pbClient.MC_ClientLogon{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdClientLogon command failed:", err)
		return false
	}

	logon, ok := msg.Body.(*pbClient.MC_ClientLogon)
	if !ok {
		logger.Info("cannot assert to yokai_client.MC_ClientLogon")
		return false
	}

	c.Connect(logon.ClientId, logon.ClientName)
	return true
}

func CmdCreatePlayer(c *TcpClient, result []string) bool {
	if !c.connected {
		logger.Warn("未连接到服务器")
		return false
	}

	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_CreatePlayer",
		Body: &pbGame.MC_CreatePlayer{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdCreatePlayer command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdSelectPlayer(c *TcpClient, result []string) bool {
	if !c.connected {
		logger.Warn("未连接到服务器")
		return false
	}

	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_SelectPlayer",
		Body: &pbGame.MC_SelectPlayer{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdSelectPlayer command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdSendHeartBeat(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_HeartBeat",
		Body: &pbClient.MC_HeartBeat{},
	}

	c.SendMessage(msg)

	return false
}

func CmdClientDisconnect(c *TcpClient, result []string) bool {
	c.Disconnect()
	return false
}

func CmdQueryPlayerInfos(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_QueryPlayerInfos",
		Body: &pbGame.MC_QueryPlayerInfos{},
	}

	c.SendMessage(msg)
	return true
}

func CmdChangeExp(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_ChangeExp",
		Body: &pbGame.MC_ChangeExp{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdChangeExp command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdChangeLevel(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_ChangeLevel",
		Body: &pbGame.MC_ChangeLevel{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdChangeLevel command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdQueryHeros(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_QueryHeros",
		Body: &pbGame.MC_QueryHeros{},
	}

	c.SendMessage(msg)
	return true
}

func CmdAddHero(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_AddHero",
		Body: &pbGame.MC_AddHero{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdAddHero command failed:", err)
		return false
	}

	c.SendMessage(msg)

	return true
}

func CmdDelHero(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_DelHero",
		Body: &pbGame.MC_DelHero{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdDelHero command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdHeroAddExp(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_HeroAddExp",
		Body: &pbGame.MC_HeroAddExp{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdHeroAddExp command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdHeroAddLevel(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_HeroAddLevel",
		Body: &pbGame.MC_HeroAddLevel{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdHeroAddLevel command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdQueryItems(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_QueryItems",
		Body: &pbGame.MC_QueryItems{},
	}

	c.SendMessage(msg)
	return true
}

func CmdAddItem(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_AddItem",
		Body: &pbGame.MC_AddItem{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdAddItem command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdDelItem(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_DelItem",
		Body: &pbGame.MC_DelItem{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdDelItem command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdQueryHeroEquips(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_QueryHeroEquips",
		Body: &pbGame.MC_QueryHeroEquips{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdQueryHeroEquips command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdHeroPutonEquip(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_PutonEquip",
		Body: &pbGame.MC_PutonEquip{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdHeroPutonEquip command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdHeroTakeoffEquip(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_TakeoffEquip",
		Body: &pbGame.MC_TakeoffEquip{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdHeroTakeoffEquip command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdQueryTokens(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_QueryTokens",
		Body: &pbGame.MC_QueryTokens{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdQueryTokens command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdAddToken(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_AddToken",
		Body: &pbGame.MC_AddToken{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdAddToken command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdQueryTalents(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_QueryTalents",
		Body: &pbGame.MC_QueryTalents{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdQueryTalents command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdAddTalent(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_AddTalent",
		Body: &pbGame.MC_AddTalent{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdAddTalent command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func registerCommand(c *Command) {
	cmdPage, ok := CmdPages[c.PageID]
	if !ok {
		fmt.Println("register command failed:", c)
		return
	}

	cmdPage.Cmds = append(cmdPage.Cmds, c)
	c.Number = len(cmdPage.Cmds)
}

func registerCommandPage(p *CommandPage) {
	CmdPages[p.PageID] = p
}

func initCommandPages() {

	// first level page
	// page main options
	registerCommandPage(&CommandPage{PageID: 1, ParentPageID: -1, Cmds: make([]*Command, 0)})

	// seconde level page
	// page server connection options
	registerCommandPage(&CommandPage{PageID: 2, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page role options
	registerCommandPage(&CommandPage{PageID: 3, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page hero options
	registerCommandPage(&CommandPage{PageID: 4, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page item options
	registerCommandPage(&CommandPage{PageID: 5, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page equip options
	registerCommandPage(&CommandPage{PageID: 6, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page token options
	registerCommandPage(&CommandPage{PageID: 7, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// page blade options
	registerCommandPage(&CommandPage{PageID: 8, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// third level options
}

func initCommands() {
	// first level page
	// 0服务器连接管理
	registerCommand(&Command{Text: "服务器连接管理", PageID: 1, GotoPageID: 2, Cb: nil})

	// 1角色管理
	registerCommand(&Command{Text: "角色管理", PageID: 1, GotoPageID: 3, Cb: nil})

	// 2英雄管理
	registerCommand(&Command{Text: "英雄管理", PageID: 1, GotoPageID: 4, Cb: nil})

	// 3物品管理
	registerCommand(&Command{Text: "物品管理", PageID: 1, GotoPageID: 5, Cb: nil})

	// 4装备管理
	registerCommand(&Command{Text: "装备管理", PageID: 1, GotoPageID: 6, Cb: nil})

	// 5代币管理
	registerCommand(&Command{Text: "代币管理", PageID: 1, GotoPageID: 7, Cb: nil})

	// 6异刃管理
	registerCommand(&Command{Text: "异刃管理", PageID: 1, GotoPageID: 8, Cb: nil})

	// 9退出
	registerCommand(&Command{Text: "退出", PageID: 1, GotoPageID: -1, Cb: CmdQuit})

	///////////////////////////////////////////////
	// 服务器连接管理
	///////////////////////////////////////////////
	// 返回上页
	registerCommand(&Command{Text: "返回上页", PageID: 2, GotoPageID: 1, Cb: nil})

	// 1登录
	registerCommand(&Command{Text: "登录", PageID: 2, GotoPageID: -1, InputText: "请输入登录客户端ID和名字，以逗号分隔", DefaultInput: "1,dudu", Cb: CmdClientLogon})

	// 2发送心跳
	registerCommand(&Command{Text: "发送心跳", PageID: 2, GotoPageID: -1, Cb: CmdSendHeartBeat})

	// 3断开连接
	registerCommand(&Command{Text: "断开连接", PageID: 2, GotoPageID: -1, Cb: CmdClientDisconnect})

	///////////////////////////////////////////////
	// 角色管理
	///////////////////////////////////////////////
	// 返回上页
	registerCommand(&Command{Text: "返回上页", PageID: 3, GotoPageID: 1, Cb: nil})

	// 1查询账号下所有角色
	registerCommand(&Command{Text: "查询账号下所有角色", PageID: 3, GotoPageID: -1, Cb: CmdQueryPlayerInfos})

	// 2创建角色
	registerCommand(&Command{Text: "创建角色", PageID: 3, GotoPageID: -1, InputText: "请输入角色名字", DefaultInput: "加百列", Cb: CmdCreatePlayer})

	// 3选择角色
	registerCommand(&Command{Text: "选择角色", PageID: 3, GotoPageID: -1, InputText: "请输入角色ID", DefaultInput: "1", Cb: CmdSelectPlayer})

	// 4改变经验
	registerCommand(&Command{Text: "改变经验", PageID: 3, GotoPageID: -1, InputText: "请输入要改变的经验值:", DefaultInput: "120", Cb: CmdChangeExp})

	// 5改变等级
	registerCommand(&Command{Text: "改变等级", PageID: 3, GotoPageID: -1, InputText: "请输入要改变的等级:", DefaultInput: "10", Cb: CmdChangeLevel})

	///////////////////////////////////////////////
	// 英雄管理
	///////////////////////////////////////////////
	// 返回上页
	registerCommand(&Command{Text: "返回上页", PageID: 4, GotoPageID: 1, Cb: nil})

	// 1查询英雄信息
	registerCommand(&Command{Text: "查询英雄信息", PageID: 4, GotoPageID: -1, Cb: CmdQueryHeros})

	// 2添加英雄
	registerCommand(&Command{Text: "添加英雄", PageID: 4, GotoPageID: -1, InputText: "请输入要添加的英雄TypeID:", DefaultInput: "1", Cb: CmdAddHero})

	// 3删除英雄
	registerCommand(&Command{Text: "删除英雄", PageID: 4, GotoPageID: -1, InputText: "请输入要删除的英雄ID:", DefaultInput: "1", Cb: CmdDelHero})

	// 4增加经验
	registerCommand(&Command{Text: "增加经验", PageID: 4, GotoPageID: -1, InputText: "请输入英雄id和经验，用逗号分隔:", DefaultInput: "1,110", Cb: CmdHeroAddExp})

	// 5增加等级
	registerCommand(&Command{Text: "增加等级", PageID: 4, GotoPageID: -1, InputText: "请输入英雄id和等级，用逗号分隔:", DefaultInput: "1,3", Cb: CmdHeroAddLevel})

	///////////////////////////////////////////////
	// 物品管理
	///////////////////////////////////////////////
	// 返回上页
	registerCommand(&Command{Text: "返回上页", PageID: 5, GotoPageID: 1, Cb: nil})

	// 1查询物品信息
	registerCommand(&Command{Text: "查询物品信息", PageID: 5, GotoPageID: -1, Cb: CmdQueryItems})

	// 2添加物品
	registerCommand(&Command{Text: "添加物品", PageID: 5, GotoPageID: -1, InputText: "请输入要添加的物品TypeID:", DefaultInput: "1", Cb: CmdAddItem})

	// 3删除物品
	registerCommand(&Command{Text: "删除物品", PageID: 5, GotoPageID: -1, InputText: "请输入要删除的物品ID:", DefaultInput: "1", Cb: CmdDelItem})

	///////////////////////////////////////////////
	// 装备管理
	///////////////////////////////////////////////
	// 返回上页
	registerCommand(&Command{Text: "返回上页", PageID: 6, GotoPageID: 1, Cb: nil})

	// 1查询英雄装备
	registerCommand(&Command{Text: "查询英雄装备", PageID: 6, GotoPageID: -1, InputText: "请输入英雄ID:", DefaultInput: "1", Cb: CmdQueryHeroEquips})

	// 2穿装备
	registerCommand(&Command{Text: "穿装备", PageID: 6, GotoPageID: -1, InputText: "请输入英雄ID和物品ID:", DefaultInput: "1,1", Cb: CmdHeroPutonEquip})

	// 3脱装备
	registerCommand(&Command{Text: "脱装备", PageID: 6, GotoPageID: -1, InputText: "请输入英雄ID和物品ID:", DefaultInput: "1,1", Cb: CmdHeroTakeoffEquip})

	///////////////////////////////////////////////
	// 代币管理
	///////////////////////////////////////////////
	// 返回上页
	registerCommand(&Command{Text: "返回上页", PageID: 7, GotoPageID: 1, Cb: nil})

	// 1查询代币信息
	registerCommand(&Command{Text: "查询代币信息", PageID: 7, GotoPageID: -1, Cb: CmdQueryTokens})

	// 2变更代币数量
	registerCommand(&Command{Text: "变更代币数量", PageID: 7, GotoPageID: -1, InputText: "请输入要变更的代币类型和数量，用逗号分隔:", DefaultInput: "0,1000", Cb: CmdAddToken})

	///////////////////////////////////////////////
	// 异刃管理
	///////////////////////////////////////////////
	// 返回上页
	registerCommand(&Command{Text: "返回上页", PageID: 8, GotoPageID: 1, Cb: nil})

	// 1查询天赋信息
	registerCommand(&Command{Text: "查询天赋信息", PageID: 8, GotoPageID: -1, Cb: CmdQueryTalents})

	// 2增加天赋
	registerCommand(&Command{Text: "增加天赋", PageID: 8, GotoPageID: -1, InputText: "请输入要增加的天赋id:", DefaultInput: "1", Cb: CmdAddTalent})
}
