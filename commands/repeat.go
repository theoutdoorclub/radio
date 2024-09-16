package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

	"github.com/theoutdoorclub/radio/helpers"
	"github.com/theoutdoorclub/radio/radio"
	"github.com/theoutdoorclub/radio/radio/queue"
	"github.com/theoutdoorclub/radio/shared"
)

func init() {
	Registry = append(Registry, discord.SlashCommandCreate{
		Name:        "repeat",
		Description: "Set the queue's repeat type",
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
		},
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "type",
				Description: "Repeat type",
				Required:    true,
				Choices: []discord.ApplicationCommandOptionChoiceString{
					{
						Name:  "normal",
						Value: queue.RepeatTypeNormal,
					},
					{
						Name:  "queue",
						Value: queue.RepeatTypeQueue,
					},
					{
						Name:  "track",
						Value: queue.RepeatTypeTrack,
					},
				},
			},
		},
	})

	Router.SlashCommand("/repeat", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		e.DeferCreateMessage(true)

		it := e.Ctx.Value(shared.RadioKey).(*radio.Radio)
		repeatType, ok := data.OptString("type")

		if err := helpers.VerifyOpt(e, ok, "repeat type"); err != nil {
			return err
		}

		player := it.Lavalink.Client.ExistingPlayer(*e.GuildID())
		q := it.Queues[*e.GuildID()]

		if player == nil || q == nil {
			return helpers.NoPlayerActiveRespond(e)
		}

		q.RepeatType = queue.RepeatType(repeatType)

		e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
			SetContent("it worked i think").
			SetEphemeral(true).
			Build(),
		)

		return nil
	})
}
