package command

import "github.com/go-telegram/bot"

type Command struct {
	Name        string
	Icon        string
	Key         string
	Description string
	IsFeature   bool
	System      bool
	Handler     bot.HandlerFunc
}

type CommandManager struct {
	commands map[string]*Command
}

func New() *CommandManager {
	cm := &CommandManager{
		commands: make(map[string]*Command),
	}

	cm.initSystemCommands()

	return cm
}

func (cm *CommandManager) initSystemCommands() {
	cm.Register(&Command{
		Name:        "Start",
		Key:         "start",
		Description: "Initialize the bot in the dm chat",
		IsFeature:   false,
		System:      true,
	})
}

func (cm *CommandManager) Register(cmd *Command) {
	cm.commands[cmd.Key] = cmd
}

func (cm *CommandManager) RegisterAll(cmds []*Command) {
	for _, cmd := range cmds {
		cm.Register(cmd)
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
