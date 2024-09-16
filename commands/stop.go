package commands

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/theoutdoorclub/radio/radio"
	"github.com/theoutdoorclub/radio/shared"
)

func init() {
	Registry = append(Registry, discord.SlashCommandCreate{
		Name:        "stop",
		Description: "Stop the player",
	})

	Router.SlashCommand("/stop", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		e.DeferCreateMessage(true)

		it := e.Ctx.Value(shared.RadioKey).(*radio.Radio)
		player := it.Lavalink.Client.ExistingPlayer(*e.GuildID())

		if player == nil {
			return nil
		}

		player.Destroy(context.Background())

		// TODO: investigate if this actually does delete the thing or not
		// cause pointers are fucked :sure:
		delete(it.Queues, *e.GuildID())

		it.Client.UpdateVoiceState(context.Background(), *e.GuildID(), nil, false, false)

		e.CreateFollowupMessage(discord.
			NewMessageCreateBuilder().
			SetContent("Stopped").
			SetEphemeral(true).
			Build(),
		)

		return nil
	})
}
