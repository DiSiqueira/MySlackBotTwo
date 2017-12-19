package bot

import (
	"fmt"
	"log"
	"sync"

	"github.com/disiqueira/MySlackBotTwo/pkg/config"
	"io/ioutil"
	"strings"
	"plugin"
	"os"
)

// Cmd holds the parsed user's input for easier handling of commands
type Cmd struct {
	Raw         string       // Raw is full string passed to the command
	Channel     string       // Channel where the command was called
	ChannelData *ChannelData // More info about the channel, including network
	User        *User        // User who sent the message
	Message     string       // Full string without the prefix
	MessageData *Message     // Message with extra flags
	Command     string       // Command is the first argument passed to the bot
	RawArgs     string       // Raw arguments after the command
	Args        []string     // Arguments as array
}

// ChannelData holds the improved channel info, which includes protocol and server
type ChannelData struct {
	Protocol  string // What protocol the message was sent on (irc, slack, telegram)
	Server    string // The server hostname the message was sent on
	Channel   string // The channel name the message appeared in
	IsPrivate bool   // Whether the channel is a group or private chat
}

// URI gives back an URI-fied string containing protocol, server and channel.
func (c *ChannelData) URI() string {
	return fmt.Sprintf("%s://%s/%s", c.Protocol, c.Server, c.Channel)
}

// Message holds the message info - for IRC and Slack networks, this can include whether the message was an action.
type Message struct {
	Text     string // The actual content of this Message
	IsAction bool   // True if this was a '/me does something' message
}

// PassiveCmd holds the information which will be passed to passive commands when receiving a message
type PassiveCmd struct {
	Raw         string       // Raw message sent to the channel
	MessageData *Message     // Message with extra
	Channel     string       // Channel which the message was sent to
	ChannelData *ChannelData // Channel and network info
	User        *User        // User who sent this message
}

// PeriodicConfig holds a cron specification for periodically notifying the configured channels
type PeriodicConfig struct {
	CronSpec string                               // CronSpec that schedules some function
	Channels []string                             // A list of channels to notify
	CmdFunc  func(channel string) (string, error) // func to be executed at the period specified on CronSpec
}

// User holds user id, nick and real name
type User struct {
	ID       string
	Nick     string
	RealName string
	IsBot    bool
}

type customPlugin struct {
	Cmd         string
	CmdFunc     activePluginFunc
	Description string
	ExampleArgs string
	MinArgs		int
}

func (cp *customPlugin) SetMinArgs(min int) {
	cp.MinArgs = min
}

// CmdResult is the result message of V2 commands
type CmdResult struct {
	Channel string // The channel where the bot should send the message
	Message string // The message to be sent
}

// CmdResultV3 is the result message of V3 commands
type PluginResult struct {
	Channel string
	Message chan string
	Done    chan bool
}

const (
	commandNotAvailable   = "Command %v not available."
	noCommandsAvailable   = "No commands available."
	errorExecutingCommand = "Error executing %s: %s"
	seeUsage = "Invalid args, see usage with: !help %s."
)

type (
	passiveCmdFunc func(cmd *PassiveCmd) (string, error)

	activePluginFunc interface {
		Init()
		Command() string
		Description() string
		ExampleArgs() string
		Execute([]string, string) (string, string, error)
	}
)

var (
	plugins 		 = make(map[string]*customPlugin)
	passiveCommands  = make(map[string]passiveCmdFunc)
	periodicCommands = make(map[string]PeriodicConfig)
	cfgs config.Specification
)

func RegisterPlugin(cmd activePluginFunc) {
	plugins[cmd.Command()] = &customPlugin{
		Cmd:         cmd.Command(),
		CmdFunc:     cmd,
		Description: cmd.Description(),
		ExampleArgs: cmd.ExampleArgs(),
	}
}

// RegisterPassiveCommand adds a new passive command to the bot.
// The command should be registered in the Init() func of your package
// Passive commands receives all the text posted to a channel without any parsing
// command: String used to identify the command, for internal use only (ex: logs)
// cmdFunc: Function which will be executed. It will received the raw message, channel and nick
func RegisterPassiveCommand(command string, cmdFunc func(cmd *PassiveCmd) (string, error)) {
	passiveCommands[command] = cmdFunc
}

// RegisterPeriodicCommand adds a command that is run periodically.
// The command should be registered in the Init() func of your package
// config: PeriodicConfig which specify CronSpec and a channel list
// cmdFunc: A no-arg function which gets triggered periodically
func RegisterPeriodicCommand(command string, config PeriodicConfig) {
	periodicCommands[command] = config
}

func EnablePeriodicCommand(config PeriodicConfig) {
	StartPeriodicCommand(std, config)
}

//// TODO Last time solution this needs to be rethink
func RegisterConfigs(newConfig config.Specification) {
	cfgs = newConfig
}

func LoadPlugins() {
	files, err := ioutil.ReadDir("./plugins")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileDetails := strings.Split(file.Name(),".")
		if fileDetails[1] != "so" {
			continue
		}
		fmt.Printf("Loading: %s \n",file.Name())
		plug, err := plugin.Open(fmt.Sprintf("./plugins/%s",file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		symCustomPlugin, err := plug.Lookup("CustomPlugin")
		if err != nil {
			log.Fatal(err)
		}

		var customPlugin activePluginFunc
		customPlugin, ok := symCustomPlugin.(activePluginFunc)
		if !ok {
			fmt.Println("unexpected type from module symbol")
			os.Exit(1)
		}
		RegisterPlugin(customPlugin)
	}
}

func Configs() config.Specification {
	return cfgs
}

// Disable allows disabling commands that were registered.
// It is usefull when running multiple bot instances to disabled some plugins like url which
// is already present on some protocols.
func (b *Bot) Disable(cmds []string) {
	b.disabledCmds = append(b.disabledCmds, cmds...)
}

func (b *Bot) executePassiveCommands(cmd *PassiveCmd) {
	var wg sync.WaitGroup
	mutex := &sync.Mutex{}

	for k, v := range passiveCommands {
		if b.isDisabled(k) {
			continue
		}

		cmdFunc := v
		wg.Add(1)

		go func() {
			defer wg.Done()

			result, err := cmdFunc(cmd)
			if err != nil {
				log.Println(err)
			} else {
				mutex.Lock()
				b.handlers.Response(cmd.Channel, result, cmd.User)
				mutex.Unlock()
			}
		}()
	}
	wg.Wait()
}

func (b *Bot) isDisabled(cmd string) bool {
	for _, c := range b.disabledCmds {
		if c == cmd {
			return true
		}
	}
	return false
}

func (b *Bot) handlePlugin(c *Cmd) {
	cmd := plugins[c.Command]

	if cmd == nil {
		log.Printf("Command not found %v \n", c.Command)
		return
	}

	if cmd.MinArgs > len(c.Args) {
		b.handlers.Response(c.Channel, fmt.Sprintf(seeUsage,c.Command), c.User)
		return
	}

	message, channel, err := cmd.CmdFunc.Execute(c.Args, c.Channel)
	if err != nil {
		log.Println(err)
		return
	}
	if channel == "" {
		channel = c.Channel
	}
	if message != "" {
		b.handlers.Response(channel, message, c.User)
	}
}

func (b *Bot) checkCmdError(err error, c *Cmd) {
	if err != nil {
		errorMsg := fmt.Sprintf(errorExecutingCommand, c.Command, err.Error())
		log.Printf(errorMsg)
		b.handlers.Response(c.Channel, errorMsg, c.User)
	}
}
