package snek

type CommandGroup []Command

func (c *CommandGroup) Add(commands ...Command) {
	*c = append(*c, commands...)
}

type Commands map[string]*CommandGroup

func (c *Commands) Group(name string) *CommandGroup {
	if *c == nil {
		*c = make(Commands)
	}

	val, ok := (*c)[name]
	if ok {
		return val
	}
	(*c)[name] = &CommandGroup{}
	return (*c)[name]
}
