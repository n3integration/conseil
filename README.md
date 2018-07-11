# conseil [ ![Codeship Status for n3integration/conseil](https://app.codeship.com/projects/a59591e0-66f2-0136-ff01-5e954b520d34/status?branch=master)](https://app.codeship.com/projects/297512)

Inspired by the CapitalGo "Rapid Application Development" talk. This
project is a work in progress. If you're interested in
following the project progress, click the :star: above to be notified
when updates are available.

### Installation

```sh
go get -u github.com/n3integration/conseil

```

## Usage

### Bootstrap a New Application

When bootstrapping a new Go application, there are multiple steps
that are required for each project. Use the `new` command to build
out the application scaffolding for a new application including:
git repo, dependency management, etc.

```sh
NAME:
   conseil new - bootstrap a new application

USAGE:
   conseil new [command options] [arguments...]

OPTIONS:
   --framework value  app framework [i.e. echo, gin, iris, ozzo, grpc] (default: "gin")
   --host value       ip address to bind (default: "127.0.0.1")
   --port value       local port to bind (default: 8080)
   --migrations       whether or not to include support for database migrations
   --driver value     database driver (default: "postgres")
   --dep              whether or not to initialize dependency management through dep
   --git              whether or not to initialize git repo
```

Once executed, the following project structure is setup:

```sh
.
|-- .gitignore            (*requires --git)
|-- Gopkg.lock            (*requires --dep)
|-- Gopkg.toml            (*requires --dep)
|-- app.go
`-- sql                   (*requires --migrations)
    |-- migrations
    |   |-- 1.down.sql
    |   `-- 1.up.sql
    `-- sql.go

3 directories, 7 files

```

The `app.go` file contains a basic application for the framework specified,
which includes a single stubbed `/health` endpoint.

#### Database Migrations

If your application requires database migrations, enable the `migrations`
option. This will setup a `sql/sql.go` file that initializes the database driver.
It will also create a `sql/migrations` folder that contains skeleton `up` and
`down` migration templates. Otherwise, the `driver` option is ignored.

