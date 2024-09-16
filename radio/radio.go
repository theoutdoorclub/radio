package radio

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v3/disgolink"

	"github.com/theoutdoorclub/radio/radio/queue"
)

// Radio holds data and fields and shits so we don't have to
// do the fucked up thing of asserting from contexts
type Radio struct {
	Client bot.Client
	Config Config

	QueueManager queue.QueueManager

	Lavalink struct {
		Client disgolink.Client
		Nodes  []disgolink.Node
	}
}

func (it *Radio) SetupListeners() {
	it.Lavalink.Client.AddListeners(disgolink.NewListenerFunc(it.onTrackEnd))
}
