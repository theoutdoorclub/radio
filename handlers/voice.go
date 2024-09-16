package handlers

import (
	"context"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"

	"github.com/theoutdoorclub/radio/radio"
)

func OnGuildVoiceStateUpdate(it *radio.Radio) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.GuildVoiceStateUpdate) {
		// filter all non bot voice state updates out
		if event.VoiceState.UserID != it.Client.ApplicationID() {
			return
		}

		it.Lavalink.Client.OnVoiceStateUpdate(
			context.Background(),
			event.VoiceState.GuildID,
			event.VoiceState.ChannelID,
			event.VoiceState.SessionID,
		)

		if event.VoiceState.ChannelID == nil {
			// channelID being nil means it left vc
			delete(it.Queues, event.VoiceState.GuildID)
		}
	})
}

func OnGuildVoiceServerUpdate(it *radio.Radio) bot.EventListener {
	return bot.NewListenerFunc(func(event *events.VoiceServerUpdate) {
		it.Lavalink.Client.OnVoiceServerUpdate(
			context.Background(),
			event.GuildID,
			event.Token,
			*event.Endpoint,
		)
	})
}
