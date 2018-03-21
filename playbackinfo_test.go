package main

import (
	"fmt"
	"testing"

	. "github.com/franela/goblin"
)

func TestCompilePlaybackInfoTemplate(t *testing.T) {
	g := Goblin(t)
	g.Describe("compilePlaybackInfoTemplate", func() {
		cases := map[string]string{
			``:                       ``,
			`test`:                   `{{"test"}}`,
			`"quotes"`:               `{{"\"quotes\""}}`,
			`sla/\shes`:              `{{"sla/\\shes"}}`,
			`%u%n%k%n%o%w%n`:         `{{"u"}}{{"n"}}{{"k"}}{{"n"}}{{"o"}}{{"w"}}{{"n"}}`,
			`percent %% percent`:     `{{"percent "}}%{{" percent"}}`,
			`%foo %Abar %T [%S%I] %`: `{{"foo "}}{{.Artist}}{{"bar "}}{{.Title}}{{" ["}}{{.Status}}{{.StatusIcon}}{{"] "}}`,
			`%%%I %A — %T%%`:         `%{{.StatusIcon}}{{" "}}{{.Artist}}{{" — "}}{{.Title}}%`,
		}
		for k, v := range cases {
			passed := k
			expected := v
			g.It(fmt.Sprintf("Should compile `%s` to `%s`", passed, expected), func() {
				compiled := compilePlaybackInfoTemplate(passed)
				g.Assert(compiled.Root.String()).Equal(expected + "\n")
			})
		}

	})
}
