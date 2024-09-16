package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

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

		repeatType, ok := data.OptString("type")
		if !ok {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("Invalid repeat type").
				SetEphemeral(true).
				Build(),
			)

			return nil
		}

		it := e.Ctx.Value(shared.RadioKey).(radio.Radio)
		player := it.Lavalink.Client.ExistingPlayer(*e.GuildID())
		q := it.QueueManager.Queues[*e.GuildID()]

		if player == nil || q == nil {
			e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetContent("No player is active").
				SetEphemeral(true).
				Build(),
			)

			return nil
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
