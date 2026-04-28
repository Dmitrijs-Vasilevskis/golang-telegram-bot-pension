package menu

func (mm *MenuManager) registerCallbackRoutes() {
	mm.router.Register("nav:main", mm.handleMain)

	mm.router.Register("nav:chats", mm.handleChats)

	mm.router.Register("nav:chat:select", mm.handleSelectChat)
	mm.router.Register("nav:chat:selected", mm.handleSelectedChat)

	mm.router.Register("nav:settings", mm.handleNavSettings)

	mm.router.Register("nav:commands", mm.handleCommands)
	mm.router.Register("nav:commands:all", mm.handleAllCommands)
	mm.router.Register("nav:commands:select", mm.handleSelectCommand)

	mm.router.Register("nav:features", mm.handleFeatures)
	mm.router.Register("nav:features:select", mm.handleSelectFeature)

	mm.router.Register("cfg:command", mm.handleToggleCommand)
	mm.router.Register("cfg:feature", mm.handleToggleFeature)

	mm.router.Register("nav:feature:duplicate_dm", mm.handleDuplicateMode)
}

func (mm *MenuManager) registerMessageRoutes() {
	mm.messageRouter.Register(
		StateDuplicateMessage,
		mm.handleDuplicateMessage,
	)
}
