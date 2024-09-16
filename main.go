package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/maniartech/signals"

	"github.com/theoutdoorclub/radio/commands"
	"github.com/theoutdoorclub/radio/handlers"
	"github.com/theoutdoorclub/radio/radio"
	"github.com/theoutdoorclub/radio/radio/queue"
	"github.com/theoutdoorclub/radio/shared"
)

var (
	shouldSyncCommands *bool
)

func init() {
	shouldSyncCommands = flag.Bool("sync-commands", false, "Whether to sync the commands")
	flag.Parse()
}

func parseConfig(it *radio.Radio) {
	conf, err := radio.ParseConfig()
	if err == nil {
		it.Config = conf
	} else {
		shared.Logger.Fatal().Err(err).Msg("Failed to parse config")
	}
}

func createClient(it *radio.Radio) {
	token := ""
	if os.Getenv("TOKEN") != "" {
		token = os.Getenv("TOKEN")
	} else {
		token = it.Config.Credentials.Token
	}

	client, err := disgo.New(token,
		// set gateway options
		bot.WithGatewayConfigOpts(
			// set enabled intents
			gateway.WithIntents(
				gateway.IntentsNonPrivileged,
			),
		),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagsAll)),
	)
	if err != nil {
		shared.Logger.Fatal().Err(err).Msg("Failed to create client")
	}

	it.Client = client
}

func syncCommands(it *radio.Radio) {
	if *shouldSyncCommands {
		if err := handler.SyncCommands(it.Client, commands.Registry, []snowflake.ID{}); err != nil {
			shared.Logger.Err(err).Msg("Failed to sync commands")
		}
	}
}

func setupLavalink(it *radio.Radio) {
	lavalinkClient := disgolink.New(it.Client.ApplicationID())

	it.Client.AddEventListeners(
		handlers.OnGuildVoiceStateUpdate(it),
		handlers.OnGuildVoiceServerUpdate(it),
	)

	lavalinkClient.AddListeners(
		handlers.OnTrackEnded(it),
	)

	// connect to lavalink nodes as defined in config
	for _, nodeCfg := range it.Config.Nodes {
		ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFn()

		node, err := lavalinkClient.AddNode(ctx, disgolink.NodeConfig{
			Name:      nodeCfg.Name,
			Address:   nodeCfg.Address,
			Password:  nodeCfg.Password,
			Secure:    nodeCfg.Secure,
			SessionID: "",
		})

		if err != nil {
			shared.Logger.Err(err).Msg("Failed to connect to node " + nodeCfg.Name)
			continue
		}

		it.Lavalink.Nodes = append(it.Lavalink.Nodes, node)
	}

	it.Lavalink.Client = lavalinkClient
}

func main() {
	shared.Logger.Info().Msg("hhm yes")
	it := &radio.Radio{
		Queues:             map[snowflake.ID]*queue.Queue{},
		AddedToQueueSignal: signals.New[snowflake.ID](),
	}

	parseConfig(it)
	createClient(it)
	syncCommands(it)
	setupLavalink(it)

	defer it.Client.Close(context.Background())

	// connect to the gateway
	if err := it.Client.OpenGateway(context.Background()); err != nil {
		shared.Logger.Fatal().Err(err).Msg("Failed to connect to gateway")
	}

	commands.Router.DefaultContext(func() context.Context {
		// this context is made available to every handler so we can access the Radio object
		// from anywhere
		ctx := context.WithValue(context.Background(), shared.RadioKey, it)
		return ctx
	})

	it.Client.AddEventListeners(commands.Router)

	// other event listeners/handlers here
	handlers.OnAddedToQueue(it)

	shared.Logger.Info().Msg("its up yo")

	//
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
