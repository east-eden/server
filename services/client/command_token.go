package client

func (cmd *Commander) initTokenCommands() {
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Token, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Token, GotoPageID: Cmd_Page_Main, Cb: nil})
}
