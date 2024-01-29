package live

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/unitoftime/nootbot/pkg/httputils"
)

const GameStatusApiUrl = "https://alpha.mythfall.com:7779/api/status"

var LiveBanners = []Banner{
	{
		ChannelId: "1201270649365213215",
		Func: func(b *BannerSystem) string {
			return fmt.Sprintf("🤺 Players: %d", b.cachedStatus.NumPlayers)
		},
	},
	{
		ChannelId: "1201270754419953765",
		Func: func(b *BannerSystem) string {
			lowestSphere := findSphereWithLowestStability(b.cachedStatus.Spheres)
			return fmt.Sprintf("🀄 Stability: %d", lowestSphere.Stability)
		},
	},
}

type GameStatus struct {
	NumPlayers uint     `json:"NumPlayers"`
	Spheres    []Sphere `json:"Spheres"`
}

type Sphere struct {
	Stability uint64 `json:"Stability"`
}

type Banner struct {
	ChannelId string
	Func      func(b *BannerSystem) string
}

type BannerSystem struct {
	discord   *discordgo.Session

	refreshTime time.Duration
	banners            []Banner

	// Basic cache, so we can reuse instead of dispatching requests per banner
	cachedStatus *GameStatus
}

func (b *BannerSystem) updateBanners() {
	var gameStatus GameStatus
	err := httputils.GetJson(GameStatusApiUrl, &gameStatus)
	if err != nil {
		log.Println("Unable to fetch updated data for live banners: ", err)
		return
	}
	b.cachedStatus = &gameStatus

	if b.cachedStatus != nil {
		for _, banner := range LiveBanners {
			newName := banner.Func(b)

			// Goroutine because of horrible design by the discordgo library
			go editChannelName(b.discord, banner.ChannelId, newName)
		}
		log.Println("Updated banners at", time.Now())
	} else {
		log.Println("No data was available for banners, skipping")
	}
}

func editChannelName(discord *discordgo.Session, channelId string, name string) {
	_, err := discord.ChannelEdit(channelId, name)
	if err != nil {
		log.Println("Unable to edit channel", channelId, ":", err)
	}
}

func (b *BannerSystem) Listen() {
	// Only at most twice every 10 minutes because of discord limits
	log.Println("Live banner system listening")
	for {
		b.updateBanners()
		time.Sleep(b.refreshTime)
	}
}

func NewBannerSystem(session *discordgo.Session, banners []Banner, refreshTime time.Duration) *BannerSystem {
	return &BannerSystem{
		discord:   session,

		refreshTime: refreshTime,
		banners:            banners,

		cachedStatus: nil,
	}
}

// Messy placement, but we can move it somewhere later
func findSphereWithLowestStability(spheres []Sphere) Sphere {
	minSphere := Sphere{Stability: 10_000}

	for _, sphere := range spheres {
		if sphere.Stability < minSphere.Stability {
			minSphere = sphere
		}
	}
	return minSphere
}
