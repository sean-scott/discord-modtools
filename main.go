package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Token    string `json:"token"`
	ClientID string `json:"client_id"`
}

var bot Bot

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func purgeMessages(s *discordgo.Session, ch_id string, count int) {
	var ids []string

	messages, mErr := s.ChannelMessages(ch_id, count, "", "", "")
	if mErr == nil {
		for _, m := range messages {
			ids = append(ids, m.ID)
		}
	}

	err := s.ChannelMessagesBulkDelete(ch_id, ids)
	if err != nil {
		//fmt.Println(err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	var purge bool

	ch_id := m.ChannelID
	content := m.Content
	mentions := m.Mentions
	user := m.Author

	// get the guild ID of the sent message
	msg_channel, _ := s.Channel(ch_id)
	g_id := msg_channel.GuildID

	purge_usage := "Why don't you have a seat, " + user.Username +
		", and I'll tell you how to use 'purge'...\n" +
		"'purge' will delete the last 5 messages.\n" +
		"'purge x' will delete the last 'x' messages." +
		" 'x' must be an integer between 1 and 100."

	// commands are called via mentioning the bot
	// check against the client ID to determine if it was mentioned
	if len(mentions) != 0 && mentions[0].ID == bot.ClientID {
		var roles []string
		state := s.State

		// get the roles for the user
		member, mErr := s.GuildMember(g_id, user.ID)
		if mErr == nil {
			for _, r := range member.Roles {
				role, rErr := state.Role(g_id, r)
				if rErr == nil {
					roles = append(roles, role.Name)
				}
			}
		}

		// remove the mention and get the remaining message content
		content = strings.TrimSpace(strings.TrimPrefix(content,
			"<@"+bot.ClientID+">"))

		// non-Admin commands here
		if strings.EqualFold(content, "hello") {
			s.ChannelMessageSend(ch_id, "Why don't you have a seat...")
		} else {
			// make sure they have "Admin"
			var isAdmin bool

			for _, r := range roles {
				if r == "Admin" {
					isAdmin = true
				}
			}

			// These commands only execute if user is "Admin" role
			if isAdmin {

				// check if command was called
				purge = strings.HasPrefix(content, "purge")

				// purge the specified value of messages
				if purge {
					// expected input:
					// 'purge' to default to delete 5 last messages
					// 'purge x' where x is an integer between 1-100

					value := strings.TrimSpace(strings.TrimPrefix(content, "purge"))

					if len(value) > 0 {
						i, err := strconv.Atoi(value)
						if err != nil {
							s.ChannelMessageSend(ch_id, purge_usage)
						} else {
							if i > 100 {
								s.ChannelMessageSend(ch_id, "I'm going to only delete 100")
								purgeMessages(s, ch_id, 100)
							} else if i < 1 {
								s.ChannelMessageSend(ch_id, "That's not valid")
							} else {
								purgeMessages(s, ch_id, i)
							}
						}
					} else {
						purgeMessages(s, ch_id, 5)
					}
				}
			} else {
				s.ChannelMessageSend(ch_id, "<@"+user.ID+"> you aren't "+
					"authorized. Must have the `Admin` role!")
			}
		}
	}
}

func setupBot() {
	dat, err := ioutil.ReadFile("info.json")
	check(err)

	jErr := json.Unmarshal(dat, &bot)
	check(jErr)
}

func main() {
	setupBot()

	dg, err := discordgo.New("Bot " + bot.Token)
	check(err)

	// register messageCreate as a callback for the messageCreate events
	dg.AddHandler(messageCreate)

	// open the websocket and begin listening
	dg.Open()

	// simple way to keep program running until any key press
	var input string
	fmt.Scanln(&input)
	dg.Close()
	return
}
