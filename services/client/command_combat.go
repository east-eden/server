package client

func (cmd *Commander) initCombatCommands() {
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Combat, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Combat, GotoPageID: Cmd_Page_Main, Cb: nil})

}
