go 1.15

replace (
	discordbot/bot v0.0.0 => ./bot
	discordbot/config v0.0.0 => ./config
	discordbot/table v0.0.0 => ./table
)

require (
	discordbot/bot v0.0.0
	discordbot/table v0.0.0
	discordbot/config v0.0.0
	github.com/bwmarrin/discordgo v0.22.0
	maunium.net/go/mautrix v0.7.13
)

module local/discordbot
