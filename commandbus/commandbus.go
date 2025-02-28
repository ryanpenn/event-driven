package commandbus

import (
	"context"
	"sync"
)

// CommandName is the name of a command
type CommandName string

// Command is the interface for a command
type Command interface {
	Name() CommandName
	SetNext(Command) Command
	Execute(context.Context) error
}

// BaseCommand is the base implementation of a command
type BaseCommand struct {
	next Command
}

// SetNext sets the next command in the chain
func (b *BaseCommand) SetNext(next Command) Command {
	b.next = next
	return b.next
}

// Execute executes the command and passes the execution to the next command in the chain
func (b *BaseCommand) Execute(ctx context.Context) error {
	// execute the command logic here
	// optionally call the next command in the chain
	if b.next != nil {
		return b.next.Execute(ctx)
	}
	return nil
}

// CommandBus is the interface for a command bus
type CommandBus interface {
	// Set sets a command in the bus
	Set(CommandName, Command)
	// Get gets a command from the bus by name
	Get(CommandName) Command
	// Has checks if a command is in the bus
	Has(CommandName) bool
	// Unset unsets a command from the bus
	Unset(CommandName)
}

type commandbus struct {
	lock     *sync.RWMutex
	commands map[CommandName]Command
}

func NewCommandBus() CommandBus {
	return &commandbus{
		lock:     &sync.RWMutex{},
		commands: make(map[CommandName]Command),
	}
}

func (c *commandbus) Set(name CommandName, command Command) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.commands[name] = command
}

func (c *commandbus) Get(name CommandName) Command {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.commands[name]
}

func (c *commandbus) Has(name CommandName) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.commands[name]
	return ok
}

func (c *commandbus) Unset(name CommandName) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.commands, name)
}
