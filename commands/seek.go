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
		Name:        "seek",
		Description: "Jumps to ",
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
		},
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionInt{
				Name:        "count",
				Description: "How many tracks to skip",
				Required:    false,
				MinValue:    json.Ptr(0),
				Defa
			},
		},
	})

	Router.SlashCommand("/skip", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		e.DeferCreateMessage(true)

		it := e.Ctx.Value(shared.RadioKey).(*radio.Radio)
		count, valid := data.OptInt("count")

		if !valid {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().SetContent("Invalid count supplied").
				SetEphemeral(true).
				Build(),
			)

			return nil
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

		// i think stop player -> pop queue -> emit addedtoqueuesignal
		// hm yes how do we stop player :running
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

		q.PopTo(1)
		it.AddedToQueueSignal.Emit(context.Background(), *e.GuildID())

		e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
			SetContent("hey it worked").
			SetEphemeral(true).
			Build(),
		)

		return nil
	})
}
