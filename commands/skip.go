package commands

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"

	"github.com/theoutdoorclub/radio/radio"
	"github.com/theoutdoorclub/radio/shared"
)

func init() {
	Registry = append(Registry, discord.SlashCommandCreate{
		Name:        "skip",
		Description: "Skips the music duh",
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
		},
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionInt{
				Name:        "count",
				Description: "How many tracks to skip",
				Required:    false,
				MinValue:    json.Ptr(1),
			},
		},
	})

	Router.SlashCommand("/skip", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		e.DeferCreateMessage(true)

		it := e.Ctx.Value(shared.RadioKey).(*radio.Radio)
		count, supplied := data.OptInt("count")

		if !supplied {
			count = 1
		}

		q, ok := it.Queues[*e.GuildID()]
		if !ok {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("No player is active").
				SetEphemeral(true).
				Build(),
			)

			return nil
		}

		if count < 1 || count > len(q.QueuedTracks) {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Cannot skip that many tracks. Queue doesn't have that many.").
				SetEphemeral(true).
				Build(),
			)

			return nil
		}

		player := it.Lavalink.Client.Player(*e.GuildID())

		// setting it to a Null track should pass all the checks in OnAddedToQueue i think
		if err := player.Update(context.Background(), lavalink.WithNullTrack()); err != nil {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Something went wrong").
				SetEphemeral(true).
				Build(),
			)

			return err
		}

		q.PopTo(count)
		it.AddedToQueueSignal.Emit(context.Background(), *e.GuildID())

		e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
			SetContent("hey it worked").
			SetEphemeral(true).
			Build(),
		)

		return nil
	})
}
