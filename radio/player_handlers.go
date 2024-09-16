package radio

import (
	"context"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"

	"github.com/theoutdoorclub/radio/radio/queue"
	"github.com/theoutdoorclub/radio/shared"
)

func (it *Radio) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	if !event.Reason.MayStartNext() {
		shared.Logger.Error().Str("reason", string(event.Reason)).Msg("Cannot start next")
		return
	}

	q := it.QueueManager.GetOrCreate(event.GuildID())
	justFinishedPlayingTrack := q.QueuedTracks[0]

	var (
		nextTrack queue.QueuedTrack
		ok        bool
	)

	switch q.RepeatType {
	case queue.RepeatTypeNormal:
		nextTrack, ok = q.Next()

	case queue.RepeatTypeTrack:
		nextTrack = justFinishedPlayingTrack
		ok = true

	case queue.RepeatTypeQueue:
		nextTrack, ok = q.Next()

		q.Insert(justFinishedPlayingTrack.Queuer, justFinishedPlayingTrack.Track)
	}

	if !ok {
		shared.Logger.Error().Msg("Failed to fetch next track from queue")
		return
	}

	if err := player.Update(context.Background(), lavalink.WithTrack(nextTrack.Track)); err != nil {
		shared.Logger.Err(err).Msg("Failed to play next track")
	}
}
