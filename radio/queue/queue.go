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
	Queuer discord.User
	Track  lavalink.Track
}

type Queue struct {
	QueuedTracks []QueuedTrack
	RepeatType   RepeatType
}

func (queue *Queue) Next() (QueuedTrack, bool) {
	// empty queue or has reached the end of the queue already
	if len(queue.QueuedTracks) == 0 {
		return QueuedTrack{}, false
	}

	track := queue.QueuedTracks[0]
	queue.QueuedTracks = queue.QueuedTracks[1:]

	return track, true
}

func (queue *Queue) Insert(queuer discord.User, track lavalink.Track) {
	queue.QueuedTracks = append(queue.QueuedTracks, QueuedTrack{
		Queuer: queuer,
		Track:  track,
	})
}

type QueueManager struct {
	Queues map[snowflake.ID]*Queue
}

func (mgr *QueueManager) Create(serverId snowflake.ID) *Queue {
	mgr.Queues[serverId] = &Queue{
		QueuedTracks: []QueuedTrack{},
		RepeatType:   RepeatTypeNormal,
	}

	return mgr.Queues[serverId]
}

func (mgr *QueueManager) GetOrCreate(serverId snowflake.ID) *Queue {
	queue, ok := mgr.Queues[serverId]

	if !ok {
		queue = mgr.Create(serverId)
	}

	return queue
}
