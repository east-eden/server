package client

import (
	"fmt"
	"reflect"
	"strconv"
	"syscall"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
)

type Command struct {
	Number       int
	Text         string
	PageID       int
	GotoPageID   int
	Cb           func(*TcpClient, []string)
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

			fv.Set(reflect.ValueOf(inputValue))
		}

		if ft.Kind() == reflect.String {
			fv.Set(reflect.ValueOf(result[n]))
		}
	}

	return nil
}

func CmdQuit(c *TcpClient, result []string) {
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

func CmdClientLogon(c *TcpClient, result []string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_ClientLogon",
		Body: &pbClient.MC_ClientLogon{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdClientLogon command failed:", err)
		return
	}

	logon, ok := msg.Body.(*pbClient.MC_ClientLogon)
	if !ok {
		logger.Info("cannot assert to yokai_client.MC_ClientLogon")
		return
	}

	c.Connect(logon.ClientId, logon.ClientName)
}

func CmdCreatePlayer(c *TcpClient, result []string) {
	if !c.connected {
		logger.Warn("未连接到服务器")
		return
	}

	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_CreatePlayer",
		Body: &pbClient.MC_CreatePlayer{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdCreatePlayer command failed:", err)
		return
	}

	c.SendMessage(msg)
}

func CmdSendHeartBeat(c *TcpClient, result []string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_HeartBeat",
		Body: &pbClient.MC_HeartBeat{},
	}

	c.SendMessage(msg)
}

func CmdClientDisconnect(c *TcpClient, result []string) {
	c.Disconnect()
}

func CmdPlayerInfo(c *TcpClient, result []string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_QueryPlayerInfo",
		Body: &pbClient.MC_QueryPlayerInfo{},
	}

	c.SendMessage(msg)
}

func CmdChangeExp(c *TcpClient, result []string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_ChangeExp",
		Body: &pbClient.MC_ChangeExp{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdChangeExp command failed:", err)
		return
	}

	c.SendMessage(msg)
}

func CmdChangeLevel(c *TcpClient, result []string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_ChangeLevel",
		Body: &pbClient.MC_ChangeLevel{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		fmt.Println("CmdChangeLevel command failed:", err)
		return
	}

	c.SendMessage(msg)
}

func registerCommand(c *Command) {
	cmdPage, ok := CmdPages[c.PageID]
	if !ok {
		fmt.Println("register command failed:", c)
		return
	}

	cmdPage.Cmds = append(cmdPage.Cmds, c)
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

	// page blade options
	registerCommandPage(&CommandPage{PageID: 7, ParentPageID: 1, Cmds: make([]*Command, 0)})

	// third level options
}

func initCommands() {
	// first level page
	// 0服务器连接
	registerCommand(&Command{Number: 0, Text: "服务器连接", PageID: 1, GotoPageID: 2, Cb: nil})

	// 1角色管理
	registerCommand(&Command{Number: 1, Text: "角色管理", PageID: 1, GotoPageID: 3, Cb: nil})

	// 2英雄管理
	registerCommand(&Command{Number: 2, Text: "英雄管理", PageID: 1, GotoPageID: 4, Cb: nil})

	// 3物品管理
	registerCommand(&Command{Number: 3, Text: "物品管理", PageID: 1, GotoPageID: 5, Cb: nil})

	// 4装备管理
	registerCommand(&Command{Number: 4, Text: "装备管理", PageID: 1, GotoPageID: 6, Cb: nil})

	// 5异刃管理
	registerCommand(&Command{Number: 5, Text: "异刃管理", PageID: 1, GotoPageID: 7, Cb: nil})

	// 9退出
	registerCommand(&Command{Number: 9, Text: "退出", PageID: 1, GotoPageID: -1, Cb: CmdQuit})

	// second level page
	// 返回上页
	registerCommand(&Command{Number: 0, Text: "返回上页", PageID: 2, GotoPageID: 1, Cb: nil})

	// 1登录
	registerCommand(&Command{Number: 1, Text: "登录", PageID: 2, GotoPageID: -1, InputText: "请输入登录客户端ID和名字，以逗号分隔", DefaultInput: "1,dudu", Cb: CmdClientLogon})

	// 2创建角色
	registerCommand(&Command{Number: 2, Text: "创建角色", PageID: 2, GotoPageID: -1, InputText: "请输入角色名字", DefaultInput: "加百列", Cb: CmdCreatePlayer})

	// 3发送心跳
	registerCommand(&Command{Number: 3, Text: "发送心跳", PageID: 2, GotoPageID: -1, Cb: CmdSendHeartBeat})

	// 4断开连接
	registerCommand(&Command{Number: 4, Text: "断开连接", PageID: 2, GotoPageID: -1, Cb: CmdClientDisconnect})

	// 返回上页
	registerCommand(&Command{Number: 0, Text: "返回上页", PageID: 3, GotoPageID: 1, Cb: nil})

	// 1查询角色信息
	registerCommand(&Command{Number: 1, Text: "查询角色信息", PageID: 3, GotoPageID: -1, Cb: CmdPlayerInfo})

	// 2改变经验
	registerCommand(&Command{Number: 2, Text: "改变经验", PageID: 3, GotoPageID: -1, InputText: "请输入要改变的经验值:", DefaultInput: "120", Cb: CmdChangeExp})

	// 3改变等级
	registerCommand(&Command{Number: 3, Text: "改变等级", PageID: 3, GotoPageID: -1, InputText: "请输入要改变的等级:", DefaultInput: "10", Cb: CmdChangeLevel})
}
