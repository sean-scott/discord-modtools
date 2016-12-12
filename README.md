# discord-modtools
A simple moderation bot for Discord, written using the discordgo library:

https://github.com/bwmarrin/discordgo

Requires an 'Admin' role for some commands. This can be changed if you want to edit it.

# Setup

You should have the [Go](https://golang.org/dl/) language and the above Go library installed.

Create your bot through https://discordapp.com/developers/applications/me. [Here's](https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token) a nice tutorial on how to do it.

Get this repo by running `go get github.com/sean-scott/discord-modtools` and then `cd` into the project directory. You'll want to update the info.json files with your bot's token and client id. Once that's done, install the project via `go install`. Lastly run `$GOPATH/bin/discord-modtools` and presto!

# Features

`@botname purge[ number]` removes the last 'number' of messages up to 100. Default is 5.
`@botname hello` greets you.
