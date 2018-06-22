package actions

import (
	"log"
	"strings"
	"sync"
	"text/template"

	"gopkg.in/urfave/cli.v1"

	"github.com/n3integration/goji"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[goji] ")
}

func Commands() []cli.Command {
	return registry.actions
}

var registry = struct {
	actions []cli.Command
	sync.Mutex
}{
	actions: make([]cli.Command, 0),
}

func register(command cli.Command) {
	registry.Lock()
	defer registry.Unlock()

	registry.actions = append(registry.actions, command)
}

func parseTemplates() *template.Template {
	templates := template.New("t")
	for _, f := range goji.AssetNames() {
		if strings.HasSuffix(f, ".tpl") {
			templates = templates.New(f)
			templates.Parse(string(goji.MustAsset(f)))
		}
	}
	return templates
}
