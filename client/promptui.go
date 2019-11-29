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
)

type PromptUI struct {
	ctx       context.Context
	cancel    context.CancelFunc
	se        *promptui.Select
	po        *promptui.Prompt
	tcpClient *TcpClient
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
