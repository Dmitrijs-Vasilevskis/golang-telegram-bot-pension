package command

type Command struct {
	Name        string
	Icon        string
	Key         string
	Description string
	IsFeature   bool
	System      bool
}

type CommandManager struct {
	commands map[string]*Command
}

func New() *CommandManager {
	cm := &CommandManager{
		commands: make(map[string]*Command),
	}

	cm.initCommands()

	return cm
}

func (cm *CommandManager) initCommands() {
	cm.commands["ask"] = &Command{
		Name:        "Ask",
		Key:         "ask",
		Description: "Ask the bot any question and get an AI-powered answer.",
		IsFeature:   false,
		System:      false,
	}

	cm.commands["look"] = &Command{
		Name:        "Look",
		Key:         "look",
		Description: "Analyze messages or media content in the chat.",
		IsFeature:   false,
		System:      false,
	}

	cm.commands["status"] = &Command{
		Name:        "Status",
		Key:         "status",
		Description: "Check current bot status",
		IsFeature:   false,
		System:      true,
	}

	cm.commands["start"] = &Command{
		Name:      "Start",
		Key:       "start",
		IsFeature: false,
		System:    true,
	}

	cm.commands["factcheck"] = &Command{
		Name:        "Facksheck",
		Key:         "factcheck",
		Description: "Verify whether a statement is true or misleading.",
		IsFeature:   false,
		System:      false,
	}

	cm.commands["summary"] = &Command{
		Name: "Summary",
		Key:  "summary",
		Description: "Quickly generate a summary of recent chat activity.\n\n" +
			"Use `/summary <count>` (e.g. `/summary 200`) to choose how many messages to include.\n\n" +
			"• Maximum: 400 messages\n" +
			"• Messages are anonymized and securely processed",
		Icon:      "💾",
		IsFeature: true,
		System:    false,
	}

	cm.commands["duplicate_dm"] = &Command{
		Name:        "Duplicate DM",
		Key:         "duplicate_dm",
		Description: "Sends a copy of messages you send to the bot in private (DM) to the selected group chat.",
		Icon:        "💾",
		IsFeature:   true,
		System:      false,
	}
}

func (cm *CommandManager) GetAll() []Command {
	commands := make([]Command, 0, len(cm.commands))

	for _, cmd := range cm.commands {
		commands = append(commands, *cmd)
	}

	return commands
}

func (cm *CommandManager) GetByKey(name string) (*Command, bool) {
	cmd, exists := cm.commands[name]

	return cmd, exists
}

func (cm *CommandManager) GetFeatureCommands() []Command {
	return cm.filter(func(c *Command) bool {
		return c.IsFeature
	})
}

func (cm *CommandManager) GetSystemCommands() []Command {
	return cm.filter(func(c *Command) bool {
		return c.System
	})
}

func (cm *CommandManager) GetRegularCommands() []Command {
	return cm.filter(func(c *Command) bool {
		return !c.IsFeature
	})
}

func (cm *CommandManager) filter(fn func(*Command) bool) []Command {
	var commands []Command

	for _, cmd := range cm.commands {
		if fn(cmd) {
			commands = append(commands, *cmd)
		}
	}

	return commands
}
