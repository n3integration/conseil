package actions

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/n3integration/conseil"
)

var (
	update    bool
	templates *template.Template
)

func init() {
	flag.BoolVar(&update, "update", false, "update .golden files")
	flag.Parse()

	templates = parseTemplates()
}

func TestParseWebAppTemplates(t *testing.T) {
	apps := listApps()
	templateList := templates.Templates()

	if len(templateList) < len(apps) {
		t.Fatalf("insufficient number of templates parsed. expected at least %d; actual %d", len(apps), len(templateList))
	}

	for _, app := range apps {
		found := false
		for _, tpl := range templateList {
			if strings.Contains(tpl.Name(), app) {
				found = true
			}
		}
		if !found {
			t.Errorf("failed to locate a parsed tpl file for %s", app)
		}
	}
}

func TestCreateWebApp(t *testing.T) {
	stageTest(t, func(t *testing.T) {
		tests := []struct {
			Framework string
			Host      string
			Port      int
			Error     bool
		}{
			{"echo", "localhost", 8080, false},
			{"gin", "localhost", 8080, false},
			{"grpc", "localhost", 9000, false},
			{"iris", "localhost", 8080, false},
			{"ozzo", "localhost", 8080, false},
			{"eggio", "localhost", 8080, true},
		}

		for _, test := range tests {
			framework = test.Framework
			host = test.Host
			port = test.Port
			err := createWebApp(templates)

			if test.Error {
				if err == nil {
					t.Error("expected test to generate an error")
				}
				break
			}

			if err != nil {
				t.Errorf("failed to create %s web application: %s", framework, err)
			}

			actual, err := ioutil.ReadFile(filepath.Join(wd, "app.go"))
			golden := filepath.Join("testdata", test.Framework+".golden")
			if update {
				ioutil.WriteFile(golden, actual, 0644)
			}

			expected, _ := ioutil.ReadFile(golden)
			if !bytes.Equal(actual, expected) {
				t.Fatalf("generated %s application contents did not match: \n%s", test.Framework, actual)
			}
		}
	})
}

func TestStageMigrations(t *testing.T) {
	stageTest(t, func(t *testing.T) {
		if err := stageMigrations(templates); err != nil {
			t.Errorf("failed to stage migrations: %s", err)
		}

		if actual := conseil.FileCount(wd, ".*\\.sql"); actual != 2 {
			t.Errorf("expected 2 migration files; actual %d", actual)
		}
	})
}

func TestSetupDb(t *testing.T) {
	stageTest(t, func(t *testing.T) {
		tests := []struct {
			Driver string
			Error  bool
		}{
			{"postgres", false},
			{"sqlite3", false},
			{"oracle", true},
		}

		for _, test := range tests {
			driver = test.Driver
			err := setupDb(templates)

			if test.Error {
				if err == nil {
					t.Error("expected test to generate an error")
				}
				break
			}

			if err != nil {
				t.Errorf("failed to setup database file: %s", err)
			}

			actual, err := ioutil.ReadFile(filepath.Join(wd, "sql", "sql.go"))
			golden := filepath.Join("testdata", test.Driver+".golden")
			if update {
				ioutil.WriteFile(golden, actual, 0644)
			}

			expected, _ := ioutil.ReadFile(golden)
			if !bytes.Equal(actual, expected) {
				t.Fatalf("generated %s application contents did not match: \n%s", test.Driver, actual)
			}
		}
	})
}

func TestDepInit(t *testing.T) {
	if _, err := exec.LookPath("dep"); err != nil {
		t.Log("dep not found, skipping")
		t.Skip()
	}

	// TODO
	// stageTest(t, func(t *testing.T) {
	// 	if err := os.Chdir(wd); err != nil {
	// 		t.Fatalf("err: %s", err)
	// 	}
	//
	// 	output, err := depInit()
	// 	if err != nil {
	// 		t.Errorf("failed to initialize dep: %s\n%s", err, output)
	// 	}
	// })
}

func TestGitInit(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Log("git not found, skipping")
		t.Skip()
	}

	stageTest(t, func(t *testing.T) {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("err: %s", err)
		}

		_, err := gitInit(templates)
		if err != nil {
			t.Errorf("failed to initialize git: %s", err)
		}

		if actual := conseil.FileCount(wd, ".*\\.gitignore"); actual != 1 {
			t.Errorf("expected 1 .gitignore file; actual %d", actual)
		}
	})
}

func stageTest(t *testing.T, fn func(*testing.T)) {
	var cleanup func()
	wd, cleanup = conseil.StageTestDir(t)
	defer cleanup()
	fn(t)
}
