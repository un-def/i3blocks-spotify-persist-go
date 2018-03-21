package main

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

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
