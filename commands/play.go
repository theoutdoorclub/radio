package commands

import (
	"context"
	"errors"
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
		e.DeferCreateMessage(false)

		it := e.Ctx.Value(shared.RadioKey).(*radio.Radio)
		identifier, valid := data.OptString("identifier")

		if err := helpers.VerifyOpt(e, valid, "identifier"); err != nil {
			return err
		}

		if err := helpers.JoinVC(e); err != nil {
			return err
		}

		// context that automatically cancels after timeout
		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		loadResult, err := it.Lavalink.Client.BestNode().LoadTracks(ctx, identifier)

		defer cancelFn()

		if err != nil {
			helpers.LoadingErrorRespond(e)
			return err
		}

		// do not question the deref thats what autocomplete does
		q := helpers.GetQueue(it, e.Channel().ID(), *e.GuildID())

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
				Build(),
			)

			return nil

		case lavalink.Exception:
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Something went wrong while loading").
				Build(),
			)

			return errors.New(d.Message)
		}

		it.Lavalink.Client.Player(*e.GuildID())
		it.AddedToQueueSignal.Emit(context.Background(), *e.GuildID())

		e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
			SetAllowedMentions(&discord.AllowedMentions{RepliedUser: true}).
			SetContent("m").
			Build(),
		)

		return nil
	})
}
