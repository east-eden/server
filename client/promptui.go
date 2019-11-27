package client

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
)

type PromptUI struct {
	ctx       context.Context
	cancel    context.CancelFunc
	se        *promptui.Select
	po        *promptui.Prompt
	tcpClient *TcpClient
}

type Command struct {
	Number       int
	Text         string
	PageID       int
	GotoPageID   int
	Message      *transport.Message
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
	registerCommand(&Command{Number: 0, Text: "服务器连接", PageID: 1, GotoPageID: 2, Message: nil})

	// 1角色管理
	registerCommand(&Command{Number: 1, Text: "角色管理", PageID: 1, GotoPageID: 3, Message: nil})

	// 2英雄管理
	registerCommand(&Command{Number: 2, Text: "英雄管理", PageID: 1, GotoPageID: 4, Message: nil})

	// 3物品管理
	registerCommand(&Command{Number: 3, Text: "物品管理", PageID: 1, GotoPageID: 5,
		Message: nil})

	// 4装备管理
	registerCommand(&Command{Number: 4, Text: "装备管理", PageID: 1, GotoPageID: 6, Message: nil})

	// 5异刃管理
	registerCommand(&Command{Number: 5, Text: "异刃管理", PageID: 1, GotoPageID: 7, Message: nil})

	// second level page
	// 返回上页
	registerCommand(&Command{Number: 0, Text: "返回上页", PageID: 2, GotoPageID: 1, Message: nil})

	// 1发送登录
	registerCommand(&Command{Number: 1, Text: "发送登录", PageID: 2, GotoPageID: -1, InputText: "请输入登录客户端ID和名字，以逗号分隔", DefaultInput: "1,dudu", Message: &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_ClientLogon",
		Body: &pbClient.MC_ClientLogon{},
	}})

	// 2发送心跳
	registerCommand(&Command{Number: 2, Text: "发送心跳", PageID: 2, GotoPageID: -1, Message: &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_HeartBeat",
		Body: &pbClient.MC_HeartBeat{},
	}})

	// 返回上页
	registerCommand(&Command{Number: 0, Text: "返回上页", PageID: 3, GotoPageID: 1, Message: nil})

	// 1改变经验
	registerCommand(&Command{Number: 1, Text: "改变经验", PageID: 3, GotoPageID: -1, InputText: "请输入要改变的经验值:", DefaultInput: "120", Message: nil})

	// 2改变等级
	registerCommand(&Command{Number: 2, Text: "改变等级", PageID: 3, GotoPageID: -1, InputText: "请输入要改变的等级:", DefaultInput: "10", Message: nil})
}

func NewPromptUI(ctx *cli.Context, client *TcpClient) *PromptUI {

	initCommandPages()
	initCommands()

	ui := &PromptUI{
		se: &promptui.Select{
			Label: "方向键选择操作",
			Size:  10,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ .Text }}?",
				Active:   "  {{ .Number | cyan }} {{ .Text | cyan }}",
				Inactive: " {{ .Number | white }} {{ .Text | white }}",
				Selected: "  {{ .Number | red | cyan }} {{ .Text | red | cyan }}",
			},
		},
		po:        &promptui.Prompt{},
		tcpClient: client,
	}

	ui.ctx, ui.cancel = context.WithCancel(ctx)

	ui.se.Items = CmdPages[1].Cmds

	return ui
}

func (p *PromptUI) Run() error {
	for {
		time.Sleep(time.Millisecond * 500)

		select {
		case <-p.ctx.Done():
			logger.Info("prompt ui context done...")
			return nil
		default:
			if !p.tcpClient.connected {
				time.Sleep(time.Second)
				continue
			}
		}

		index, _, err := p.se.Run()

		if err != nil {
			fmt.Printf("Select failed %v\n", err)
			return err
		}

		items, ok := p.se.Items.([]*Command)
		if !ok {
			fmt.Println("items invalid")
			return nil
		}

		cmd := items[index]

		// jump to next page
		if nextPage := cmd.GotoPageID; nextPage != -1 {
			p.se.Items = CmdPages[nextPage].Cmds
			continue
		}

		// wait input
		var splitArgs []string
		if len(cmd.InputText) > 0 {
			p.po.Label = cmd.InputText
			p.po.Default = cmd.DefaultInput

			result, err := p.po.Run()
			if err != nil {
				fmt.Println("prompt run error:", err)
				continue
			}

			splitArgs = strings.Split(result, ",")
		}

		// trans input into cmd.Message
		if cmd.Message != nil {
			tp := reflect.TypeOf(cmd.Message.Body).Elem()
			value := reflect.ValueOf(cmd.Message.Body).Elem()

			// proto.Message struct has 3 invalid field
			if value.NumField()-3 != len(splitArgs) {
				fmt.Println("输入数据无效")
				continue
			}

			// reflect into proto.Message
			success := true
			for n := 0; n < len(splitArgs); n++ {
				ft := tp.Field(n).Type
				fv := value.Field(n)

				if ft.Kind() >= reflect.Int && ft.Kind() <= reflect.Uint64 {
					inputValue, err := strconv.ParseInt(splitArgs[n], 10, ft.Bits())
					if err != nil {
						fmt.Printf("input value<%s> cannot assert to type<%s>\r\n", splitArgs[n], ft.Name())
						success = false
						break
					}

					fv.Set(reflect.ValueOf(inputValue))
				}

				if ft.Kind() == reflect.String {
					fv.Set(reflect.ValueOf(splitArgs[n]))
				}
			}

			if success {
				p.tcpClient.SendMessage(cmd.Message)
			}
		}
	}

	return nil
}
