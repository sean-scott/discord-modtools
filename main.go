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
  Token string    `json:"token"`
  ClientID string `json:"client_id"`
}
var bot Bot

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func setupBot() {
  dat, err := ioutil.ReadFile("info.json")
  check(err)

  jErr := json.Unmarshal(dat, &bot)
  check(jErr)
}

func purgeMessages(s *discordgo.Session, ch_id string, count int) {
  var ids []string

  messages, mErr := s.ChannelMessages(ch_id, count, "", "")
  if mErr != nil {
    fmt.Println("ERROR")
  } else {
    for _, m := range messages {
      ids = append(ids, m.ID)
    }
  }

  err := s.ChannelMessagesBulkDelete(ch_id, ids)
  if err != nil {
    fmt.Println(err)
  }
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

  var purge bool

  ch_id := m.ChannelID
  content := m.Content
  mentions := m.Mentions
  //user := m.Author

  purge_usage := "Why don't you have a seat and I'll tell you" +
    "how to use `purge`...\n" +
    "`purge` will delete the last 5 messages.\n" +
    "`purge x` will delete the last `x` messages." +
    " `x` must be an integer between 1 and 100.\n"

  // commands are called via mentioning the bot
  // check against the client ID to determine if it was mentioned
  if len(mentions) != 0 && mentions[0].ID == bot.ClientID {
    // remove the mention and get the remaining message content
    content = strings.TrimSpace(strings.TrimPrefix(content,
      "<@" + bot.ClientID + ">"))

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
          message, bErr := s.ChannelMessageSend(ch_id, purge_usage)
          fmt.Println(message)
          fmt.Println(bErr)
        } else {
          if i > 100 {
            // something that says it will only delete 100...
            purgeMessages(s, ch_id, 100)
          } else if i < 1 {
            // no deletion, usage message
          } else {
            purgeMessages(s, ch_id, i)
          }
        }
      } else {
        purgeMessages(s, ch_id, 5)
      }
    }
  }
}

func main() {
  setupBot()

  dg, err := discordgo.New(bot.Token)
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
