package actions

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/n3integration/conseil"
	"github.com/pkg/errors"

	"gopkg.in/urfave/cli.v1"
)

var (
	wd        string
	driver    string
	framework string
	host      string
	port      int

	dep        bool
	git        bool
	migrations bool
)

func init() {
	register(cli.Command{
		Name:    "new",
		Aliases: []string{"n"},
		Usage:   "bootstrap a new application",
		Action:  appAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "framework",
				Value:       "gin",
				Usage:       fmt.Sprintf("app framework [i.e. %v]", strings.Join(listApps(), ", ")),
				Destination: &framework,
			},
			cli.StringFlag{
				Name:        "host",
				Value:       "127.0.0.1",
				Usage:       "ip address to bind",
				Destination: &host,
			},
			cli.IntFlag{
				Name:        "port",
				Value:       8080,
				Usage:       "local port to bind",
				Destination: &port,
			},
			cli.BoolFlag{
				Name:        "migrations",
				Destination: &migrations,
				Usage:       "whether or not to include support for database migrations",
			},
			cli.StringFlag{
				Name:        "driver",
				Value:       "postgres",
				Usage:       "database driver",
				Destination: &driver,
			},
			cli.BoolFlag{
				Destination: &dep,
				Name:        "dep",
				Usage:       "whether or not to initialize dependency management through dep",
			},
			cli.BoolFlag{
				Destination: &git,
				Name:        "git",
				Usage:       "whether or not to initialize git repo",
			},
		},
	})
}

// Context provides the application context options
type Context struct {
	App    string
	Host   string
	Port   int
	Driver string
	Conn   string
	Import string
}

func appAction(_ *cli.Context) error {
	if wd == "" {
		wd = "."
	}

	templates := parseTemplates()
	if err := createWebApp(templates); err != nil {
		return err
	}

	if migrations {
		if err := stageMigrations(templates); err != nil {
			return err
		}

		if err := setupDb(templates); err != nil {
			return err
		}
	}

	if dep {
		if out, err := depInit(); err != nil {
			return err
		} else {
			log.Println(out)
		}
	}

	if git {
		if out, err := gitInit(templates); err != nil {
			return err
		} else {
			log.Println(out)
		}
	}

	return nil
}

func createWebApp(templates *template.Template) error {
	t := templates.Lookup(fmt.Sprintf("templates/app/%s.tpl", framework))
	if t == nil {
		return errors.Errorf("unable to find a '%s' app framework template", framework)
	}

	app, err := os.Create(filepath.Join(wd, "app.go"))
	if err != nil {
		return err
	}

	log.Println("creating app...")
	context := &Context{
		Host: host,
		Port: port,
	}

	if err := t.Execute(app, context); err != nil {
		return err
	}

	if framework == "grpc" {
		path := filepath.Join(wd, "proto")
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}

		_, err := os.Create(filepath.Join(path, "rpc.proto"))
		return err
	}
	return nil
}

func stageMigrations(templates *template.Template) error {
	path := filepath.Join(wd, "sql/migrations")
	log.Println("staging migrations...")
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	up, _ := os.Create(filepath.Join(path, "1.up.sql"))
	if err := templates.Lookup("templates/sql/1.up.tpl").Execute(up, nil); err != nil {
		return err
	}

	down, _ := os.Create(filepath.Join(path, "1.down.sql"))
	return templates.Lookup("templates/sql/1.down.tpl").Execute(down, nil)
}

func setupDb(templates *template.Template) error {
	dbConn, err := conn(driver)
	if err != nil {
		return err
	}

	path := filepath.Join(wd, "sql")
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	sql, _ := os.Create(filepath.Join(path, "sql.go"))
	context := &Context{
		Driver: driver,
		Conn:   dbConn,
		Import: imp(driver),
	}
	return templates.Lookup("templates/sql/sql.tpl").Execute(sql, context)
}

func depInit() (string, error) {
	cmd := exec.Command("dep", "init")
	log.Println("initializing dependencies...")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Errorf("unable to initialize dep: %s", err)
	}

	return fmt.Sprintf("%s\n", bytes.TrimSpace(output)), nil
}

func gitInit(templates *template.Template) (string, error) {
	t := templates.Lookup("templates/gitignore.tpl")
	ign, _ := os.Create(filepath.Join(wd, ".gitignore"))
	wd, _ := os.Getwd()

	context := &Context{
		App: filepath.Base(wd),
	}

	if err := t.Execute(ign, context); err != nil {
		return "", err
	}

	cmd := exec.Command("git", "init")
	log.Println("initializing repo...")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Errorf("unable to initialize git: %s", err)
	}

	return fmt.Sprintf("%s\n", bytes.TrimSpace(output)), nil
}

func conn(driver string) (string, error) {
	wd, _ := os.Getwd()
	switch driver {
	case "postgres":
		return fmt.Sprintf("postgres://localhost:5432/%s", filepath.Base(wd)), nil
	case "sqlite3":
		return fmt.Sprintf("file:%s.sqlite", filepath.Base(wd)), nil
	}
	return "", fmt.Errorf("%s is not a supported database driver", driver)
}

func imp(driver string) string {
	switch driver {
	case "postgres":
		return "github.com/lib/pq"
	case "sqlite3":
		return "github.com/mattn/go-sqlite3"
	}
	return ""
}

func listApps() []string {
	appList := make([]string, 0)
	for _, app := range conseil.AssetNames() {
		if strings.HasPrefix(app, "templates/app/") {
			base := filepath.Base(app)
			appList = append(appList, strings.Replace(base, ".tpl", "", 1))
		}
	}
	return appList
}
