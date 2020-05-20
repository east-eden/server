package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type PromptUI struct {
	ctx    context.Context
	cancel context.CancelFunc
	se     *promptui.Select
	po     *promptui.Prompt
	c      *Client
}

func NewPromptUI(c *Client, ctx *cli.Context) *PromptUI {

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
		po: &promptui.Prompt{},
		c:  c,
	}

	ui.ctx, ui.cancel = context.WithCancel(ctx)

	ui.se.Items = c.cmder.pages[1].Cmds

	return ui
}

func (p *PromptUI) Run() error {
	for {

		select {
		case <-p.ctx.Done():
			logger.Info("prompt ui context done...")
			return nil
		default:
			time.Sleep(time.Millisecond * 500)
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
			p.se.Items = p.c.cmder.pages[nextPage].Cmds
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

		if cmd.Cb != nil {
			needRecv := cmd.Cb(splitArgs)
			if needRecv {
				timeOut := time.NewTimer(time.Second * 5)
				select {
				case <-p.c.transport.WaitRecv():
					continue
				case <-timeOut.C:
					continue
				}
			}
		}

	}

	return nil
}
