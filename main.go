package main

import (
	"os"
	"strings"
	"text/template"

	"github.com/godbus/dbus"
)

const busName = "org.mpris.MediaPlayer2.spotify"
const objectPath = "/org/mpris/MediaPlayer2"
const playerInterface = "org.mpris.MediaPlayer2.Player"

// Spotify ...
type Spotify struct {
	*dbus.Object
	playbackInfoTemplate *template.Template
}

func (spotify *Spotify) get(propName string) interface{} {
	variant, err := spotify.GetProperty(playerInterface + "." + propName)
	if err != nil {
		panic(err)
	}
	return variant.Value()
}

// GetPlaybackInfo ...
func (spotify *Spotify) GetPlaybackInfo() PlaybackInfo {
	status := spotify.get("PlaybackStatus").(string)
	metadata := spotify.get("Metadata").(map[string]dbus.Variant)
	artist := strings.Join(metadata["xesam:artist"].Value().([]string), ", ")
	title := metadata["xesam:title"].Value().(string)
	return PlaybackInfo{
		Status: status,
		Artist: artist,
		Title:  title,
	}
}

// ShowPlaybackInfo ...
func (spotify *Spotify) ShowPlaybackInfo(playbackInfo PlaybackInfo) {
	spotify.playbackInfoTemplate.Execute(os.Stdout, playbackInfo)
}

// ShowInitialPlaybackInfo ...
func (spotify *Spotify) ShowInitialPlaybackInfo() {
	info := spotify.GetPlaybackInfo()
	spotify.ShowPlaybackInfo(info)
}

// GetSpotify ...
func GetSpotify(conn *dbus.Conn, playbackInfoFormat string) Spotify {
	object := conn.Object(busName, objectPath)
	playbackInfoTemplate := compilePlaybackInfoTemplate(playbackInfoFormat)
	return Spotify{
		object.(*dbus.Object),
		playbackInfoTemplate,
	}
}

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	spotify := GetSpotify(conn, playbackInfoFormat)
	spotify.ShowInitialPlaybackInfo()
}
