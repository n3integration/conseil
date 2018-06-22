package actions

import (
	"testing"
)

func TestGetCommands(t *testing.T) {
	cmds := GetCommands()
	if len(cmds) == 0 {
		t.Fatal("no registered commands")
	}
}

func TestParseTemplates(t *testing.T) {
	base := parseTemplates()
	templates := base.Templates()
	if len(templates) == 0 {
		t.Fatal("no templates parsed")
	}
}
