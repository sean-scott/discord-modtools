package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"

  "github.com/bwmarrin/discordgo"
)

type TokenJson struct {
  Token string
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func getToken() string {
  dat, err := ioutil.ReadFile("info.json")
  check(err)

  var j TokenJson
  jErr := json.Unmarshal(dat, &j)
  check(jErr)

  return j.Token
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
  ch_id := m.ChannelID
  //content := m.Content
  //mentions := m.Mentions
  //user := m.Author

  fmt.Println(ch_id)

}

func main() {
  dg, err := discordgo.New(getToken())
  check(err)

  // Register messageCreate as a callback for the messageCreate events.
  dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	dg.Open()

  // Simple way to keep program running until any key press.
  var input string
  fmt.Scanln(&input)
  dg.Close()
  return
}
