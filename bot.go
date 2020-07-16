package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Execute     func(s *discordgo.Session, m *discordgo.MessageCreate)
	Description string
}

const COMMAND_PREFIX = "!"

var (
	Token    = flag.String("token", "", "API token (required)")
	Name     = flag.String("name", "c00kie-bot", "Name of the bot")
	Commands = map[string]Command{
		"hello": {Execute: helloCommand, Description: "Say hello"},
	}
)

func main() {
	flag.Parse()

	if *Token == "" {
		fmt.Printf("API token missing\n")
		return
	}

	Commands["list"] = Command{
		Execute:     listCommand,
		Description: "List all commands",
	}

	fmt.Printf("Starting %s...\n", *Name)

	d, err := discordgo.New("Bot " + *Token)

	if err != nil {
		fmt.Printf("Error creating new discord session: %s\n", err)
		return
	}

	err = d.Open()
	if err != nil {
		fmt.Printf("Error opening connection to discord: %s\n", err)
		return
	}

	defer d.Close()

	d.AddHandler(ready)
	d.AddHandler(messageCreate)

	fmt.Println("Press Ctrl-C to exit")

	sigIntChannel := make(chan os.Signal)
	signal.Notify(sigIntChannel, syscall.SIGINT, syscall.SIGTERM)
	<-sigIntChannel
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "Looking 4 n00bs")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	/* Ignore the bot itself */
	if m.Author.ID == s.State.User.ID {
		fmt.Println("Message from myself")
		return
	}

	if !strings.HasPrefix(m.Content, COMMAND_PREFIX) {
		fmt.Println("Not a command")
		return
	}

	commandKey := m.Content[1:]

	command, ok := Commands[commandKey]
	if !ok {
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("Unknown command !%s", commandKey),
		)

		if err != nil {
			fmt.Printf("Error sending back message: %s\n", err)
		}
	}

	command.Execute(s, m)
}

func helloCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf("Hey %s, what's up?", m.Author.Username),
	)

	if err != nil {
		fmt.Printf("Error sending back message: %s\n", err)
	}
}

func listCommand(s *discordgo.Session, m *discordgo.MessageCreate) {

	msg := " === Available commands ===\n"
	for key := range Commands {
		msg += fmt.Sprintf(
			"`!%5s`\t%s\n",
			key,
			Commands[key].Description,
		)
	}

	_, err := s.ChannelMessageSend(m.ChannelID, msg)

	if err != nil {
		fmt.Printf("Error sending back message: %s\n", err)
	}
}
