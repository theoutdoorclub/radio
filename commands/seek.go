package commands

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"

	"github.com/theoutdoorclub/radio/helpers"
	"github.com/theoutdoorclub/radio/radio"
	"github.com/theoutdoorclub/radio/shared"
)

func init() {
	Registry = append(Registry, discord.SlashCommandCreate{
		Name:        "seek",
		Description: "Jumps to timestamp",
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
		},
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "timestamp",
				Description: "What timestop to jump to",
				Required:    true,
			},
		},
	})

	Router.SlashCommand("/seek", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		e.DeferCreateMessage(true)

		it := e.Ctx.Value(shared.RadioKey).(*radio.Radio)
		timestamp, valid := data.OptString("timestamp")

		if err := helpers.VerifyOpt(e, valid, "timestamp"); err != nil {
			return err
		}

		player := it.Lavalink.Client.Player(*e.GuildID())
		_, ok := it.Queues[*e.GuildID()]

		if !ok || player.Track() == nil {
			return helpers.NoPlayerActiveRespond(e)
		}

		seekPosition, err := time.ParseDuration(timestamp)
		if err != nil {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Invalid timestamp").
				SetEphemeral(true).
				Build(),
			)

			return err
		}

		seekPosition = time.Duration(seekPosition.Milliseconds())
		trackLength := time.Duration(player.Track().Info.Length)

		shared.Logger.Debug().Any("track_length", trackLength).Any("seek", seekPosition).Msg("Sought")

		if seekPosition > trackLength {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContentf("Timestamp too long! Track duration is only `%v`", trackLength*time.Millisecond).
				SetEphemeral(true).
				Build(),
			)

			return nil
		}

		if err := player.Update(context.Background(), lavalink.WithPosition(lavalink.Duration(seekPosition))); err != nil {
			helpers.GenericErrorRespond(e)
			return err
		}

		e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
			SetContent("hey it worked").
			SetEphemeral(true).
			Build(),
		)

		return nil
	})
}
