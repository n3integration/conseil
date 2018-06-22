package actions

import (
	"strings"
	"sync"
	"text/template"

	"gopkg.in/urfave/cli.v1"

	"github.com/n3integration/conseil"
)

var registry = struct {
	actions []cli.Command
	mu      sync.Mutex
}{
	actions: make([]cli.Command, 0),
}

func register(command cli.Command) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.actions = append(registry.actions, command)
}

func GetCommands() []cli.Command {
	return registry.actions
}

func parseTemplates() *template.Template {
	templates := template.New("t")
	for _, f := range conseil.AssetNames() {
		if strings.HasSuffix(f, ".tpl") {
			templates = templates.New(f)
			templates.Parse(string(conseil.MustAsset(f)))
		}
	}
	return templates
}
