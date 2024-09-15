package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	Registry = append(Registry, discord.SlashCommandCreate{
		Name:        "ping",
		Description: "Send back a pong",
	})

	Router.SlashCommand("/ping", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		return e.CreateMessage(discord.
			NewMessageCreateBuilder().
			SetContent("Pong").
			SetEphemeral(true).
			Build(),
		)
	})
}
