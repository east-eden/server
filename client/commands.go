package client

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
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

func CmdAccountLogon(c *TcpClient, result []string) bool {
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
		logger.Warn("json marshal failed when call CmdAccountLogon:", err)
		return false
	}

	resp, err := httpPost(c, header, body)
	if err != nil {
		logger.Warn("http post failed when call CmdAccountLogon:", err)
		return false
	}

	var gameInfo struct {
		UserID     int64  `json:"userId"`
		UserName   string `json:"userName"`
		AccountID  int64  `json:"accountId"`
		GameID     string `json:"gameId"`
		PublicAddr string `json:"publicAddr"`
		Section    string `json:"section"`
	}

	if err := json.Unmarshal(resp, &gameInfo); err != nil {
		logger.Warn("json unmarshal failed when call CmdAccountLogon:", err)
		return false
	}

	logger.Info("metadata unmarshaled result:", gameInfo)

	if len(gameInfo.PublicAddr) == 0 {
		logger.Warn("invalid game_addr")
		return false
	}

	c.SetTcpAddress(gameInfo.PublicAddr)
	c.SetUserInfo(gameInfo.UserID, gameInfo.AccountID, gameInfo.UserName)
	c.Connect()
	return true
}

func CmdCreatePlayer(c *TcpClient, result []string) bool {
	if !c.connected {
		logger.Warn("未连接到服务器")
		return false
	}

	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_CreatePlayer",
		Body: &pbGame.C2M_CreatePlayer{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdCreatePlayer command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return true
}

func CmdExpirePlayer(c *TcpClient, result []string) bool {
	if !c.connected {
		logger.Warn("未连接到服务器")
		return false
	}

	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_ExpirePlayer",
		Body: &pbGame.MC_ExpirePlayer{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdExpirePlayer command failed:", err)
		return false
	}

	c.SendMessage(msg)
	return false
}

func CmdSendHeartBeat(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_account.C2M_HeartBeat",
		Body: &pbAccount.C2M_HeartBeat{},
	}

	c.SendMessage(msg)

	return false
}

func CmdAccountDisconnect(c *TcpClient, result []string) bool {
	c.Disconnect()
	return false
}

func CmdQueryPlayerInfo(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_QueryPlayerInfo",
		Body: &pbGame.C2M_QueryPlayerInfo{},
	}

	c.SendMessage(msg)
	return true
}

func CmdChangeExp(c *TcpClient, result []string) bool {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_ChangeExp",
		Body: &pbGame.C2M_ChangeExp{},
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
	registerCommand(&Command{Text: "登录", PageID: 2, GotoPageID: -1, InputText: "请输入登录user ID和名字，以逗号分隔", DefaultInput: "100001,dudu", Cb: CmdAccountLogon})

	// 2发送心跳
	registerCommand(&Command{Text: "发送心跳", PageID: 2, GotoPageID: -1, Cb: CmdSendHeartBeat})

	// 3断开连接
	registerCommand(&Command{Text: "断开连接", PageID: 2, GotoPageID: -1, Cb: CmdAccountDisconnect})

	///////////////////////////////////////////////
	// 角色管理
	///////////////////////////////////////////////
	// 返回上页
	registerCommand(&Command{Text: "返回上页", PageID: 3, GotoPageID: 1, Cb: nil})

	// 1查询账号下所有角色
	registerCommand(&Command{Text: "查询账号下所有角色", PageID: 3, GotoPageID: -1, Cb: CmdQueryPlayerInfo})

	// 2创建角色
	registerCommand(&Command{Text: "创建角色", PageID: 3, GotoPageID: -1, InputText: "请输入角色名字", DefaultInput: "加百列", Cb: CmdCreatePlayer})

	// 3角色缓存失效
	registerCommand(&Command{Text: "角色缓存失效", PageID: 3, GotoPageID: -1, Cb: CmdExpirePlayer})

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
	//registerCommand(&Command{Text: "增加经验", PageID: 4, GotoPageID: -1, InputText: "请输入英雄id和经验，用逗号分隔:", DefaultInput: "1,110", Cb: CmdHeroAddExp})

	// 5增加等级
	//registerCommand(&Command{Text: "增加等级", PageID: 4, GotoPageID: -1, InputText: "请输入英雄id和等级，用逗号分隔:", DefaultInput: "1,3", Cb: CmdHeroAddLevel})

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
	registerCommand(&Command{Text: "脱装备", PageID: 6, GotoPageID: -1, InputText: "请输入英雄ID和装备位置索引:", DefaultInput: "1,0", Cb: CmdHeroTakeoffEquip})

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
	registerCommand(&Command{Text: "查询天赋信息", PageID: 8, GotoPageID: -1, InputText: "请输入异刃ID:", DefaultInput: "1", Cb: CmdQueryTalents})

	// 2增加天赋
	registerCommand(&Command{Text: "增加天赋", PageID: 8, GotoPageID: -1, InputText: "请输入异刃ID和天赋ID:", DefaultInput: "1,1", Cb: CmdAddTalent})

	//expression, err := govaluate.NewEvaluableExpression("atk*2 + 10")
	//parameters := make(map[string]interface{}, 8)
	//parameters["foo"] = -1
	//parameters["atk"] = 3
	//result, err := expression.Evaluate(parameters)
	//if err != nil {
	//logger.Info("expression result:", result)
	//}
}
