package main

import (
	"errors"
	"fmt"

	"git.containerum.net/ch/solutions/pkg/db"
	"git.containerum.net/ch/solutions/pkg/db/postgres"
	"git.containerum.net/ch/solutions/pkg/server"
	"git.containerum.net/ch/solutions/pkg/server/impl"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	portFlag         = "port"
	solutionsFlag    = "solutions"
	debugFlag        = "debug"
	textlogFlag      = "textlog"
	dbFlag           = "db"
	dbPGLoginFlag    = "db_pg_login"
	dbPGPasswordFlag = "db_pg_password"
	dbPGAddrFlag     = "db_pg_addr"
	dbPGNameFlag     = "db_pg_dbname"
	dbPGNoSSLFlag    = "db_pg_nossl"
	dbMigrationsFlag = "db_migrations"
	kubeURLFlag      = "kube_url"
	resourceURLFlag  = "resource_url"
	corsFlag         = "cors"
)

var flags = []cli.Flag{
	cli.StringFlag{
		EnvVar: "PORT",
		Name:   portFlag,
		Value:  "6767",
		Usage:  "port for solutions server",
	},
	cli.StringFlag{
		EnvVar: "SOLUTIONS",
		Name:   solutionsFlag,
		Value:  "impl",
		Usage:  "Solutions impl",
	},
	cli.BoolFlag{
		EnvVar: "DEBUG",
		Name:   debugFlag,
		Usage:  "Start the server in Debug mode",
	},
	cli.BoolFlag{
		EnvVar: "TEXTLOG",
		Name:   textlogFlag,
		Usage:  "Display output log in text format",
	},
	cli.StringFlag{
		EnvVar: "DB",
		Name:   dbFlag,
		Value:  "postgres",
		Usage:  "DB for project",
	},
	cli.StringFlag{
		EnvVar: "PG_LOGIN",
		Name:   dbPGLoginFlag,
		Usage:  "DB Login (PostgreSQL)",
	},
	cli.StringFlag{
		EnvVar: "PG_PASSWORD",
		Name:   dbPGPasswordFlag,
		Usage:  "DB Password (PostgreSQL)",
	},
	cli.StringFlag{
		EnvVar: "PG_ADDR",
		Name:   dbPGAddrFlag,
		Usage:  "DB Address (PostgreSQL)",
	},
	cli.StringFlag{
		EnvVar: "PG_DBNAME",
		Name:   dbPGNameFlag,
		Usage:  "DB name (PostgreSQL)",
	},
	cli.BoolFlag{
		EnvVar: "PG_NOSSL",
		Name:   dbPGNoSSLFlag,
		Usage:  "DB disable ssl (PostgreSQL)",
	},
	cli.StringFlag{
		EnvVar: "MIGRATIONS_PATH",
		Name:   dbMigrationsFlag,
		Value:  "../../pkg/migrations/",
		Usage:  "Location of DB migrations",
	},
	cli.StringFlag{
		EnvVar: "KUBE_API_URL",
		Name:   kubeURLFlag,
		Value:  "http://kube-api:1214",
		Usage:  "Kube-API service URL",
	},
	cli.StringFlag{
		EnvVar: "RESOURCE_URL",
		Name:   resourceURLFlag,
		Value:  "http://resource-service:1213",
		Usage:  "Resource service URL",
	},
	cli.BoolFlag{
		EnvVar: "CORS",
		Name:   "cors",
		Usage:  "enable CORS",
	},
}

func setupLogs(c *cli.Context) {
	if c.Bool("debug") {
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.InfoLevel)
	}

	if c.Bool("textlog") {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func getSolutionsSrv(c *cli.Context, services server.Services) (server.SolutionsService, error) {
	switch c.String(solutionsFlag) {
	case "impl":
		return impl.NewSolutionsImpl(services), nil
	default:
		return nil, errors.New("invalid solutions impl")
	}
}

func getDB(c *cli.Context) (db.DB, error) {
	switch c.String(dbFlag) {
	case "postgres":
		url := fmt.Sprintf("postgres://%v:%v@%v/%v", c.String(dbPGLoginFlag), c.String(dbPGPasswordFlag), c.String(dbPGAddrFlag), c.String(dbPGNameFlag))
		if c.Bool(dbPGNoSSLFlag) {
			url = url + "?sslmode=disable"
		}
		return postgres.DBConnect(url, c.String(dbMigrationsFlag))
	default:
		return nil, errors.New("invalid db")
	}
}
