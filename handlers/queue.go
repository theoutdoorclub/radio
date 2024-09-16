package handlers

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"

	"github.com/theoutdoorclub/radio/radio"
	"github.com/theoutdoorclub/radio/radio/queue"
	"github.com/theoutdoorclub/radio/shared"
)

func OnAddedToQueue(it *radio.Radio) {
	it.AddedToQueueSignal.AddListener(func(ctx context.Context, guildId snowflake.ID) {
		shared.Logger.Debug().Msg(guildId.String())

		player := it.Lavalink.Client.ExistingPlayer(guildId)
		queue, ok := it.Queues[guildId]

		// doesn't already have a player, does nothing
		if player == nil || !ok {
			return
		}

		// player is paused, does nothing
		if player.Paused() {
			return
		}

		// player is already playing, does nothing
		if player.Track() != nil {
			return
		}

		if err := player.Update(context.Background(), lavalink.WithTrack(queue.QueuedTracks[0].Track)); err != nil {
			shared.Logger.Err(err).Msg("Failed to play track")
		}
	})
}

func OnTrackEnded(it *radio.Radio) disgolink.EventListener {
	return disgolink.NewListenerFunc(func(player disgolink.Player, event lavalink.TrackEndEvent) {
		q := it.Queues[event.GuildID()]
		endedTrack := q.QueuedTracks[0]

		shared.Logger.Debug().Any("player", player).Any("event", event).Msg("Track ended")

		it.Client.Rest().CreateMessage(
			endedTrack.OriginChannel,
			discord.NewMessageCreateBuilder().
				SetAllowedMentions(&discord.AllowedMentions{Parse: []discord.AllowedMentionType{}}).
				SetContentf("Track `%s` ended. %s",
					endedTrack.Track.Info.Title,
					*endedTrack.Track.Info.URI,
				).
				Build(),
		)

		// track probably died so don't start next
		if !event.Reason.MayStartNext() {
			return
		}

		// player is paused, does nothing
		if player.Paused() {
			return
		}

		switch q.RepeatType {
		case queue.RepeatTypeNormal:
			newQueue := q.QueuedTracks[1:]
			q.QueuedTracks = newQueue

			if len(newQueue) <= 0 {
				shared.Logger.Debug().Any("player", player).Msg("Queued ended")

				it.Client.Rest().CreateMessage(
					q.OriginChannel,
					discord.NewMessageCreateBuilder().
						SetContent("QUEUE ENDED").
						Build(),
				)

				return
			}

		case queue.RepeatTypeTrack:
			// do nothing lol

		case queue.RepeatTypeQueue:
			// pop the queue and append the current track back to the end of the queue
			newQueue := q.QueuedTracks[1:]
			newQueue = append(newQueue, endedTrack)

			q.QueuedTracks = newQueue
		}

		it.AddedToQueueSignal.Emit(context.Background(), event.GuildID())
	})
}
