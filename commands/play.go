package commands

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"

	"github.com/theoutdoorclub/radio/radio"
	"github.com/theoutdoorclub/radio/shared"
)

func init() {
	Registry = append(Registry, discord.SlashCommandCreate{
		Name:        "play",
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

	Router.SlashCommand("/play", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
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
		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancelFn()

		loadResult, err := it.Lavalink.Client.BestNode().LoadTracks(ctx, identifier)

		if err != nil {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Something went wrong while loading").
				SetEphemeral(true).
				Build(),
			)

			return err
		}

		queue := it.QueueManager.GetOrCreate(*e.GuildID())

		// wtf this is some black magic shit
		switch d := loadResult.Data.(type) {
		case lavalink.Track:
			queue.Insert(e.User(), d)

		case lavalink.Playlist:
			for _, track := range d.Tracks {
				queue.Insert(e.User(), track)
			}

		case lavalink.Search:

		case lavalink.Empty:
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Nothing was found for this identifier").
				SetEphemeral(true).
				Build(),
			)

			return nil

		case lavalink.Exception:
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Something went wrong while loading").
				SetEphemeral(true).
				Build(),
			)

			return err
		}

		player := it.Lavalink.Client.Player(*e.GuildID())
		if player.Track() != nil {
			// already playing a track, don't override it
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Queued").
				SetEphemeral(true).
				Build(),
			)

			return nil
		}

		// isn't playing a track already, play the first one in queue
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

		err = it.Client.UpdateVoiceState(context.TODO(), *e.GuildID(), voiceState.ChannelID, false, true)
		if err != nil {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Failed to join voice channel plz check permissions").
				SetEphemeral(true).
				Build(),
			)

			return err
		}

		// play
		toPlay := queue.QueuedTracks[0]
		err = player.Update(context.Background(), lavalink.WithTrack(toPlay.Track))

		if err != nil {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Player refused to play").
				SetEphemeral(true).
				Build(),
			)

			return err
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
