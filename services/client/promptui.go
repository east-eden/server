package client

import (
	"fmt"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type PromptUI struct {
	se *promptui.Select
	po *promptui.Prompt
	c  *Client
}

func NewPromptUI(ctx *cli.Context, c *Client) *PromptUI {

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

	ui.se.Items = c.cmder.pages[1].Cmds

	return ui
}

func (p *PromptUI) Run(ctx *cli.Context) error {
	enable := ctx.Bool("prompt_ui")
	if !enable {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("prompt ui context done...")
			return nil
		default:
			time.Sleep(time.Millisecond * 200)
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
				return fmt.Errorf("PromptUI.Run failed: %w", err)
			}

			splitArgs = strings.Split(result, ",")
		}

		if cmd.Cb != nil {
			waitReturnMsg, msgNames := cmd.Cb(ctx, splitArgs)
			if waitReturnMsg {
				chTimeOut := time.After(time.Second * 5)
				select {
				case name := <-p.c.transport.ReturnMsgName():
					names := strings.Split(msgNames, ",")
					hit := false
					for _, n := range names {
						if n == name {
							hit = true
							break
						}
					}

					if hit {
						continue
					}

				case <-chTimeOut:
					continue
				}
			}
		}

	}
}
