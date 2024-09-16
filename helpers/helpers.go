package helpers

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"

	"github.com/theoutdoorclub/radio/radio"
	"github.com/theoutdoorclub/radio/radio/queue"
	"github.com/theoutdoorclub/radio/shared"
)

func JoinVC(e *handler.CommandEvent) error {
	it := e.Ctx.Value(shared.RadioKey).(*radio.Radio)
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

	return nil
}

func VerifyOpt(e *handler.CommandEvent, valid bool, name string) error {
	if !valid {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContentf("No %s was provided", name).
			SetEphemeral(true).
			Build(),
		)
	}

	return nil
}

func GetQueue(it *radio.Radio, origin snowflake.ID, guildID snowflake.ID) *queue.Queue {
	q, ok := it.Queues[guildID]
	if !ok {
		q = &queue.Queue{
			QueuedTracks:  []queue.QueuedTrack{},
			RepeatType:    queue.RepeatTypeNormal,
			OriginChannel: origin,
		}

		it.Queues[guildID] = q
	}

	return q
}

func LoadingErrorRespond(e *handler.CommandEvent) {
	e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
		SetContent("Something went wrong while loading").
		SetEphemeral(true).
		Build(),
	)
}

func GenericErrorRespond(e *handler.CommandEvent) {
	e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
		SetContent("Something went wrong").
		SetEphemeral(true).
		Build(),
	)
}

// This returns error so it can return nil, so the callers can then just
// return the call itself
func NoPlayerActiveRespond(e *handler.CommandEvent) error {
	e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
		SetContent("No player is active").
		SetEphemeral(true).
		Build(),
	)

	return nil
}
