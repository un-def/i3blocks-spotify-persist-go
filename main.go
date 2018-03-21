package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/godbus/dbus"
)

const busName = "org.mpris.MediaPlayer2.spotify"
const objectPath = "/org/mpris/MediaPlayer2"
const playerInterface = "org.mpris.MediaPlayer2.Player"

var statusToIcon = map[string]string{
	"Playing": "\uf04b",
	"Paused":  "\uf04c",
	"Stopped": "\uf04d",
}

var playbackInfoFormatPlaceholders = map[string]string{
	"S": "{{.Status}}",
	"I": "{{.StatusIcon}}",
	"A": "{{.Artist}}",
	"T": "{{.Title}}",
	"%": "%",
}

var playbackInfoFormatRegexp = regexp.MustCompile(`%[SIAT%]|[^%]+`)

var playbackInfoFormat = `\%% %I %A — "%T" {{ %S }} %%/`

// PlaybackInfo ...
type PlaybackInfo struct {
	Status string
	Artist string
	Title  string
}

// StatusIcon ...
func (playbackInfo PlaybackInfo) StatusIcon() string {
	return statusToIcon[playbackInfo.Status]
}

// Spotify ...
type Spotify struct {
	*dbus.Object
	playbackInfoTemplate *template.Template
}

// Get ...
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

func compilePlaybackInfoTemplate(playbackInfoFormat string) *template.Template {
	var templateText string
	for _, substr := range playbackInfoFormatRegexp.FindAllString(playbackInfoFormat, -1) {
		if strings.HasPrefix(substr, "%") {
			substr = playbackInfoFormatPlaceholders[substr[1:]]
		} else {
			substr = "{{" + fmt.Sprintf("%q", substr) + "}}"
		}
		templateText += substr
	}
	templateText += "\n"
	return template.Must(template.New("PlaybackInfo").Parse(templateText))
}

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	spotify := GetSpotify(conn, playbackInfoFormat)
	spotify.ShowInitialPlaybackInfo()
}
