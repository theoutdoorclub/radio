package radio

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/maniartech/signals"

	"github.com/theoutdoorclub/radio/radio/queue"
)

// Radio holds data and fields and shits so we don't have to
// do the fucked up thing of asserting from contexts
type Radio struct {
	Client bot.Client
	Config Config

	// Map of GuildID -> queue
	Queues map[snowflake.ID]*queue.Queue

	Lavalink struct {
		Client disgolink.Client
		Nodes  []disgolink.Node
	}

	AddedToQueueSignal signals.Signal[snowflake.ID]
}
