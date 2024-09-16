package commands

import (
	"context"
	"errors"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"

	"github.com/theoutdoorclub/radio/radio"
	"github.com/theoutdoorclub/radio/radio/queue"
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
		it := e.Ctx.Value(shared.RadioKey).(*radio.Radio)
		identifier, valid := data.OptString("identifier")

		if !valid {
			return e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("No identifier was provided").
				SetEphemeral(true).
				Build(),
			)
		}

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

		// context that automatically cancels after timeout
		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		loadResult, err := it.Lavalink.Client.BestNode().LoadTracks(ctx, identifier)

		defer cancelFn()

		if err != nil {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Something went wrong while loading").
				SetEphemeral(true).
				Build(),
			)

			return err
		}

		q, ok := it.Queues[*e.GuildID()]
		if !ok {
			q = &queue.Queue{
				QueuedTracks:  []queue.QueuedTrack{},
				RepeatType:    queue.RepeatTypeNormal,
				OriginChannel: e.Channel().ID(),
			}

			it.Queues[*e.GuildID()] = q
		}

		// wtf this is some black magic shit
		switch d := loadResult.Data.(type) {
		case lavalink.Track:
			q.Insert(e.User(), e.Channel().ID(), d)

		case lavalink.Playlist:
			for _, track := range d.Tracks {
				q.Insert(e.User(), e.Channel().ID(), track)
			}

		case lavalink.Search:
			// TODO: implement

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

			return errors.New(d.Message)
		}

		it.Lavalink.Client.Player(*e.GuildID())
		it.AddedToQueueSignal.Emit(context.Background(), *e.GuildID())

		e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
			SetContent("hey it worked").
			SetEphemeral(true).
			Build(),
		)

		return nil
	})
}
