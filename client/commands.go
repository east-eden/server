package client

import (
	"fmt"

	"github.com/yokaiio/yokai_server/internal/transport"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
)

type Command struct {
	Number       int
	Text         string
	PageID       int
	GotoPageID   int
	Cb           func(*TcpClient, string)
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

func CmdSendHeartBeat(c *TcpClient, result string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_HeartBeat",
		Body: &pbClient.MC_HeartBeat{},
	}

	c.SendMessage(msg)
}

func CmdClientDisconnect(c *TcpClient, result string) {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_ClientDisconnect",
		Body: &pbClient.MC_ClientDisconnect{},
	}

	c.SendMessage(msg)
}

func CmdChangeExp(c *TcpClient, result string) {

}

func CmdChangeLevel(c *TcpClient, result string) {

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
	registerCommand(&Command{Number: 0, Text: "服务器连接", PageID: 1, GotoPageID: 2, CB: nil})

	// 1角色管理
	registerCommand(&Command{Number: 1, Text: "角色管理", PageID: 1, GotoPageID: 3, CB: nil})

	// 2英雄管理
	registerCommand(&Command{Number: 2, Text: "英雄管理", PageID: 1, GotoPageID: 4, CB: nil})

	// 3物品管理
	registerCommand(&Command{Number: 3, Text: "物品管理", PageID: 1, GotoPageID: 5, CB: nil})

	// 4装备管理
	registerCommand(&Command{Number: 4, Text: "装备管理", PageID: 1, GotoPageID: 6, CB: nil})

	// 5异刃管理
	registerCommand(&Command{Number: 5, Text: "异刃管理", PageID: 1, GotoPageID: 7, CB: nil})

	// second level page
	// 返回上页
	registerCommand(&Command{Number: 0, Text: "返回上页", PageID: 2, GotoPageID: 1, CB: nil})

	// 1登录
	registerCommand(&Command{Number: 1, Text: "登录", PageID: 2, GotoPageID: -1, InputText: "请输入登录客户端ID和名字，以逗号分隔", DefaultInput: "1,dudu", CB: nil})

	// 2发送心跳
	registerCommand(&Command{Number: 2, Text: "发送心跳", PageID: 2, GotoPageID: -1, nil, CB: CmdSendHeartBeat})

	// 3断开连接
	registerCommand(&Command{Number: 3, Text: "断开连接", PageID: 2, GotoPageID: -1, CB: CmdClientDisconnect})

	// 返回上页
	registerCommand(&Command{Number: 0, Text: "返回上页", PageID: 3, GotoPageID: 1, CB: nil})

	// 1改变经验
	registerCommand(&Command{Number: 1, Text: "改变经验", PageID: 3, GotoPageID: -1, InputText: "请输入要改变的经验值:", DefaultInput: "120", CB: CmdChangeExp})

	// 2改变等级
	registerCommand(&Command{Number: 2, Text: "改变等级", PageID: 3, GotoPageID: -1, InputText: "请输入要改变的等级:", DefaultInput: "10", CB: CmdChangeLevel})
}
