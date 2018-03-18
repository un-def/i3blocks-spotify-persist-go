package main

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
)

const busName = "org.mpris.MediaPlayer2.spotify"
const objectPath = "/org/mpris/MediaPlayer2"
const playerInterface = "org.mpris.MediaPlayer2.Player"

var statusToIcon = map[string]string{
	"Playing": "",
	"Paused":  "",
	"Stopped": "",
}

// TrackInfo ...
type TrackInfo struct {
	artist []string
	title  string
}

// Spotify ...
type Spotify struct {
	*dbus.Object
}

// Get ...
func (spotify *Spotify) get(propName string) interface{} {
	variant, err := spotify.GetProperty(playerInterface + "." + propName)
	if err != nil {
		panic(err)
	}
	return variant.Value()
}

// GetPlaybackStatus ...
func (spotify *Spotify) GetPlaybackStatus() string {
	return spotify.get("PlaybackStatus").(string)
}

// GetTrackInfo ...
func (spotify *Spotify) GetTrackInfo() TrackInfo {
	metadata := spotify.get("Metadata").(map[string]dbus.Variant)
	artist := metadata["xesam:artist"].Value().([]string)
	title := metadata["xesam:title"].Value().(string)
	return TrackInfo{
		artist: artist,
		title:  title,
	}
}

// ShowInfo ...
func (spotify *Spotify) ShowInfo(playbackStatus string, trackInfo TrackInfo) {
	fmt.Printf(
		"%s %s — %s\n",
		statusToIcon[playbackStatus],
		strings.Join(trackInfo.artist, ", "),
		trackInfo.title,
	)
}

// ShowInitialInfo ...
func (spotify *Spotify) ShowInitialInfo() {
	plackbackStatus := spotify.GetPlaybackStatus()
	trackInfo := spotify.GetTrackInfo()
	spotify.ShowInfo(plackbackStatus, trackInfo)
}

// GetSpotify ...
func GetSpotify(conn *dbus.Conn) Spotify {
	object := conn.Object(busName, objectPath)
	return Spotify{object.(*dbus.Object)}
}

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	spotify := GetSpotify(conn)
	spotify.ShowInitialInfo()
}
