package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"

	"github.com/theoutdoorclub/radio/commands"
	"github.com/theoutdoorclub/radio/radio"
	"github.com/theoutdoorclub/radio/shared"
)

var (
	shouldSyncCommands *bool
)

func init() {
	shouldSyncCommands = flag.Bool("sync-commands", false, "Whether to sync the commands")
	flag.Parse()
}

func main() {
	shared.Logger.Info().Msg("hhm yes")
	it := radio.Radio{}

	conf, err := radio.ParseConfig()
	if err == nil {
		it.Config = conf
	} else {
		shared.Logger.Fatal().Err(err).Msg("Failed to parse config")
	}

	token := ""
	if os.Getenv("TOKEN") != "" {
		token = os.Getenv("TOKEN")
	} else {
		token = conf.Credentials.Token
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

	// sync commands to discord
	if *shouldSyncCommands {
		if err = handler.SyncCommands(client, commands.Registry, []snowflake.ID{}); err != nil {
			shared.Logger.Err(err).Msg("Failed to sync commands")
			return
		}
	}

	defer client.Close(context.Background())

	// connect to the gateway
	if err = client.OpenGateway(context.Background()); err != nil {
		shared.Logger.Fatal().Err(err).Msg("Failed to connect to gateway")
	}

	// lavalink stuff
	lavalinkClient := disgolink.New(client.ApplicationID())

	client.AddEventListeners(
		// onVoiceStateUpdate
		bot.NewListenerFunc(func(event *events.GuildVoiceStateUpdate) {
			// filter all non bot voice state updates out
			if event.VoiceState.UserID != client.ApplicationID() {
				return
			}
			lavalinkClient.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
		}),

		// onVoiceServerUpdate
		bot.NewListenerFunc(func(event *events.VoiceServerUpdate) {
			lavalinkClient.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
		}),
	)

	// connect to lavalink nodes as defined in config
	for _, nodeCfg := range conf.Nodes {
		node, err := lavalinkClient.AddNode(context.Background(), disgolink.NodeConfig{
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

	commands.Router.DefaultContext(func() context.Context {
		// this context is made available to every handler so we can access the Radio object
		// from anywhere
		ctx := context.WithValue(context.Background(), shared.RadioKey, it)
		return ctx
	})

	client.AddEventListeners(commands.Router)

	shared.Logger.Info().Msg("its up yo")

	// vvv TEST TEST TEST TEST TEST vvv
	/*
		var toPlay *lavalink.Track
		node.LoadTracksHandler(context.Background(), "https://www.youtube.com/watch?v=XDd9Yb0JvjE", disgolink.NewResultHandler(
			func(track lavalink.Track) {
				// Loaded a single track
				toPlay = &track
			},
			func(playlist lavalink.Playlist) {
				// HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE HERE
				// Loaded a playlist
			},
			func(tracks []lavalink.Track) {
				// Loaded a search result
			},
			func() {
				// nothing matching the query found
			},
			func(err error) {
				// something went wrong while loading the track
				shared.Logger.Err(err).Msg("")
			},
		))

		if toPlay == nil {
			panic("WTF")
		}
		// WHY ISNT IT WORKING
		client.UpdateVoiceState(context.Background(), snowflake.MustParse("741967925690368070"), json.Ptr(snowflake.MustParse("998214208690864148")), false, true)

		player := lavalinkClient.Player(snowflake.MustParse("741967925690368070"))
		player.Update(context.Background(), lavalink.WithTrack(*toPlay))
	*/
	// ^^^ TEST TEST TEST TEST TEST ^^^

	//
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

// used to silence unused complaints :ogre:
func noop(t any) {}
