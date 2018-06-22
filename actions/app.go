package actions

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

var (
	driver    string
	framework string
	host      string
	port      int

	dep bool
	git bool
)

type Context struct {
	App    string
	Host   string
	Port   int
	Driver string
	Conn   string
	Import string
}

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
				Usage:       "web framework",
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

func appAction(c *cli.Context) error {
	templates := parseTemplates()
	if err := createWebApp(templates); err != nil {
		log.Fatal(err)
	}

	if err := stageMigrations(templates); err != nil {
		log.Fatal(err)
	}

	if err := setupDb(templates); err != nil {
		log.Fatal(err)
	}

	if dep {
		if out, err := depInit(); err != nil {
			log.Fatal(err)
		} else {
			log.Println(out)
		}
	}

	if git {
		if out, err := gitInit(templates); err != nil {
			log.Fatal(err)
		} else {
			log.Println(out)
		}
	}

	return nil
}

func createWebApp(templates *template.Template) error {
	t := templates.Lookup(fmt.Sprintf("templates/rest/%s.tpl", framework))
	if t == nil {
		return errors.Errorf("unable to find a '%s' web framework template", framework)
	}
	app, err := os.Create("app.go")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("creating app...")
	context := &Context{
		Host: host,
		Port: port,
	}
	return t.Execute(app, context)
}

func stageMigrations(templates *template.Template) error {
	path := "sql/migrations"
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
	sql, _ := os.Create(filepath.Join("sql", "sql.go"))
	context := &Context{
		Driver: driver,
		Conn:   conn(driver),
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
	ign, _ := os.Create(".gitignore")
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

func conn(driver string) string {
	wd, _ := os.Getwd()
	switch driver {
	case "postgres":
		return fmt.Sprintf("postgres://localhost:5432/%s", filepath.Base(wd))
	}
	return fmt.Sprintf("file://%s", filepath.Base(wd))
}

func imp(driver string) string {
	switch driver {
	case "postgres":
		return "github.com/lib/pq"
	}
	return "github.com/mattn/go-sqlite3"
}
