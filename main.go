package main

import "github.com/Shemistan/healths_service/internal"

func main() {
	var bot internal.Bot
	bot.Init("saves.bot")
	bot.Receive()
	defer bot.Stop()
}
