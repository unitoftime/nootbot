package live

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jasonlvhit/gocron"
	"github.com/unitoftime/nootbot/pkg/httputils"
	"log"
	"time"
)

const GameStatusApiUrl = "https://alpha.mythfall.com/api/status"

var LiveBanners = []Banner{
	{
		ChannelId: "// TODO: SET CHANNEL ID //",
		Func: func(b *BannerSystem) string {
			return fmt.Sprintf("🤺 Players: %d", b.cachedStatus.NumPlayers)
		},
	},
	{
		ChannelId: "// TODO: SET CHANNEL ID //",
		Func: func(b *BannerSystem) string {
			lowestSphere := findSphereWithLowestStability(b.cachedStatus.Spheres)
			return fmt.Sprintf("🀄 Stability: %d/10000", lowestSphere.Stability)
		},
	},
}

type GameStatus struct {
	NumPlayers uint     `json:"NumPlayers"`
	Spheres    []Sphere `json:"Spheres"`
}

type Sphere struct {
	Stability uint `json:"Stability"`
}

type Banner struct {
	ChannelId string
	Func      func(b *BannerSystem) string
}

type BannerSystem struct {
	scheduler *gocron.Scheduler
	discord   *discordgo.Session

	refreshTimeMinutes uint64
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
	err := gocron.Every(b.refreshTimeMinutes).Minutes().Do(b.updateBanners)
	if err != nil {
		log.Println(err)
	}

	log.Println("Live banner system listening")

	<-gocron.Start()
}

func NewBannerSystem(session *discordgo.Session, banners []Banner, refreshTimeMinutes uint64) *BannerSystem {
	return &BannerSystem{
		scheduler: gocron.NewScheduler(),
		discord:   session,

		refreshTimeMinutes: refreshTimeMinutes,
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
