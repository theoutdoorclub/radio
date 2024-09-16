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
		Name:        "pause",
		Description: "The World.",
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
		},
	})

	Router.SlashCommand("/pause", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		it := e.Ctx.Value(shared.RadioKey).(*radio.Radio)
		e.DeferCreateMessage(true)

		// join vc
		voiceState, ok := it.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)

		if !ok {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("You are not in a voice channel").
				SetEphemeral(true).
				Build(),
			)

			return nil
		}

		err := it.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), voiceState.ChannelID, false, true)
		if err != nil {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Failed to join voice channel plz check permissions").
				SetEphemeral(true).
				Build(),
			)

			return err
		}

		/* if err := player.Update(context.Background(), lavalink.WithPosition(lavalink.Millisecond*lavalink.Duration(timestamp))); err != nil {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Something went wrong").
				SetEphemeral(true).
				Build(),
			)

			return err
		} */

		e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
			SetContent("hey it worked").
			SetEphemeral(true).
			Build(),
		)

		return nil
	})
}
