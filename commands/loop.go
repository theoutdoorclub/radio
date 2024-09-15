package commands

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"

	"github.com/theoutdoorclub/radio/radio"
	"github.com/theoutdoorclub/radio/shared"
)

func init() {
	Registry = append(Registry, discord.SlashCommandCreate{
		Name:        "loop",
		Description: "Play music. To search, use /search instead",
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
		},
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "identifier",
				Description: "Identifier to play",
				Required:    true,
			},
		},
	})

	Router.SlashCommand("/loop", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		// TYPE CASTING FTW
		it := e.Ctx.Value(shared.RadioKey).(radio.Radio)
		identifier, valid := data.OptString("identifier")

		if !valid {
			return e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("No identifier was provided").
				SetEphemeral(true).
				Build(),
			)
		}

		if err := e.DeferCreateMessage(true); err != nil {
			return err
		}

		// context that automatically cancels after timeout
		// TODO: playlist support
		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancelFn()



		

		var toPlay *lavalink.Track

		it.Lavalink.Client.BestNode().LoadTracksHandler(ctx, identifier, disgolink.NewResultHandler(
			func(track lavalink.Track) {
				// Loaded a single track
				toPlay = &track
			},
			func(playlist lavalink.Playlist) {
				// Loaded a playlist
			},
			func(tracks []lavalink.Track) {
				// Loaded a search result
			},
			func() {
				// nothing matching the query found
			},
			func(err error) {
				// something went wrong while loading the track
			},
		))

		// join vc
		voiceState, ok := it.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)

		if !ok {
			_, err := e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("You are not in a voice channel").
				SetEphemeral(true).
				Build(),
			)
			if err != nil {
				return err
			}
		}

		// play
		err := it.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), voiceState.ChannelID, false, true)
		if err != nil {
			_, err := e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Failed to join voice channel plz check permissions").
				SetEphemeral(true).
				Build(),
			)
			if err != nil {
				return err
			}
		}

		player := it.Lavalink.Client.Player(*e.GuildID())
		err = player.Update(context.Background(), lavalink.WithTrack(*toPlay))

		if err != nil {
			_, err := e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Player refused to play").
				SetEphemeral(true).
				Build(),
			)
			if err != nil {
				return err
			}
		}

		_, err = e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
			SetContent("hey it worked").
			SetEphemeral(true).
			Build(),
		)
		if err != nil {
			return err
		}

		return nil
	})
}
