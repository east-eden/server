package client

import (
	"context"
	"fmt"
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
	Number  int
	Text    string
	Message *transport.Message
}

func (c *Command) String() string {
	return fmt.Sprintf("number:%d, text:%s, Message:%s", c.Number, c.Text, c.Message.Name)
}

func NewPromptUI(ctx *cli.Context, client *TcpClient) *PromptUI {

	ui := &PromptUI{
		se: &promptui.Select{
			Label: "方向键选择操作",
		},
		po:        &promptui.Prompt{},
		tcpClient: client,
	}

	ui.ctx, ui.cancel = context.WithCancel(ctx)

	ui.initCommands()

	return ui
}

func (p *PromptUI) initCommands() {
	items := make([]*Command, 0)

	// 发送心跳
	items = append(items, &Command{
		Number: 0,
		Text:   "发送心跳消息",
		Message: &transport.Message{
			Type: transport.BodyProtobuf,
			Name: "yokai_client.MC_HeartBeat",
			Body: &pbClient.MC_HeartBeat{},
		},
	})

	// 发送登录消息
	items = append(items, &Command{
		Number: 1,
		Text:   "发送登录消息",
		Message: &transport.Message{
			Type: transport.BodyProtobuf,
			Name: "yokai_client.MC_ClientLogon",
			Body: &pbClient.MC_ClientLogon{
				ClientId:   1,
				ClientName: "dudu",
			},
		},
	})

	p.se.Items = items
}

func (p *PromptUI) Run() error {
	for {
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
			fmt.Printf("Prompt failed %v\n", err)
			return err
		}

		items, ok := p.se.Items.([]*Command)
		if !ok {
			fmt.Println("items invalid")
			return nil
		}

		p.tcpClient.SendMessage(items[index].Message)
	}

	return nil
}
