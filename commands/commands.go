package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Registry = []discord.ApplicationCommandCreate{}
var Router = handler.New()
