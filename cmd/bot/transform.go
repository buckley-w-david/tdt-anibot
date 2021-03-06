package main

import (
	"fmt"
	"strings"

	"github.com/buckley-w-david/anibot/pkg/anilist"
	"github.com/bwmarrin/discordgo"
)

var (
	MissingToken     string
	DirectorReaction string
	CreatorReaction  string
	StudioReactions  []string
)

func init() {
	MissingToken = "No token provided. Please run: anibot -t <bot token>"

	DirectorReaction = "👉"
	CreatorReaction = "👈"
	StudioReactions = []string{"1⃣", "2⃣", "3⃣", "4⃣", "5⃣", "6⃣", "7⃣", "8⃣", "9⃣", "🔟"}
}

// Embed transforms an anilist.MediaResposne struct into a discordgo.MessageEmbed.
func Embed(media anilist.Media) (discordgo.MessageEmbed, error) {
	coverImage := discordgo.MessageEmbedThumbnail{
		URL: media.CoverImage.Medium,
	}

	mediaType := discordgo.MessageEmbedField{
		Name:   "Media Type",
		Value:  fmt.Sprintf("%s %s %s", media.MediaType, media.Format, media.Source),
		Inline: false,
	}
	fields := []*discordgo.MessageEmbedField{&mediaType}

	// TODO: Account for more than 2 studios w.r.t reactions
	var inline func(int) bool
	if len(media.Studios.Edges)&1 == 0 {
		inline = func(int) bool { return true }
	} else {
		inline = func(i int) bool {
			return i < len(media.Studios.Edges)-1
		}
	}
	studios := make([]*discordgo.MessageEmbedField, len(media.Studios.Edges))
	for i, studio := range media.Studios.Edges {
		value := fmt.Sprintf("[%s](%s)", studio.Studio.Name, studio.Studio.SiteURL)
		studios[i] = &discordgo.MessageEmbedField{
			Name:   "Studio",
			Value:  value,
			Inline: inline(i),
		}
	}
	fields = append(fields, studios...)

	director, err := media.Director()
	if err == nil {
		value := fmt.Sprintf("[%s %s](%s)", director.Name.First, director.Name.Last, director.SiteURL)
		director := discordgo.MessageEmbedField{
			Name:   "Director " + DirectorReaction,
			Value:  value,
			Inline: true,
		}
		fields = append(fields, &director)
	}

	creator, err := media.Creator()
	if err == nil {
		value := fmt.Sprintf("[%s %s](%s)", creator.Name.First, creator.Name.Last, creator.SiteURL)

		creator := discordgo.MessageEmbedField{
			Name:   "Original Creator " + CreatorReaction,
			Value:  value,
			Inline: true,
		}
		fields = append(fields, &creator)
	}

	return discordgo.MessageEmbed{
		URL:         media.SiteURL,
		Title:       media.Title.Romaji,
		Description: strings.Replace(media.Description, "<br>", "\n", -1),
		Color:       0x00ff00,
		Thumbnail:   &coverImage,
		Fields:      fields,
	}, nil
}
