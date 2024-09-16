package queue

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type RepeatType string

const (
	RepeatTypeNormal = "normal"
	RepeatTypeQueue  = "repeat_queue"
	RepeatTypeTrack  = "repeat_track"
)

type QueuedTrack struct {
	Queuer        discord.User
	Track         lavalink.Track
	OriginChannel snowflake.ID
}

type Queue struct {
	QueuedTracks  []QueuedTrack
	RepeatType    RepeatType
	OriginChannel snowflake.ID
	IsPlaying     bool
}

func (queue *Queue) PopTo(i int) (QueuedTrack, bool) {
	if !(i > 0 && i < len(queue.QueuedTracks)) {
		return QueuedTrack{}, false
	}

	queue.QueuedTracks = queue.QueuedTracks[i:]
	return queue.QueuedTracks[0], true
}

func (queue *Queue) Insert(queuer discord.User, origin snowflake.ID, track lavalink.Track) {
	queue.QueuedTracks = append(queue.QueuedTracks, QueuedTrack{
		Queuer:        queuer,
		Track:         track,
		OriginChannel: origin,
	})
}
