package main

import (
	"errors"

	"git.containerum.net/ch/solutions/pkg/db"
	"git.containerum.net/ch/solutions/pkg/db/postgres"
	"git.containerum.net/ch/solutions/pkg/server"
	"git.containerum.net/ch/solutions/pkg/server/impl"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

const (
	portFlag         = "port"
	solutionsFlag    = "solutions"
	debugFlag        = "debug"
	textlogFlag      = "textlog"
	dbFlag           = "db"
	dbURLFlag        = "db_url"
	dbMigrationsFlag = "db_migrations"
	kubeURLFlag      = "kube_url"
	resourceURLFlag  = "resource_url"
	corsFlag         = "cors"
)

var flags = []cli.Flag{
	cli.StringFlag{
		EnvVar: "CH_SOLUTIONS_PORT",
		Name:   portFlag,
		Value:  "6767",
		Usage:  "port for solutions server",
	},
	cli.StringFlag{
		EnvVar: "CH_SOLUTIONS",
		Name:   solutionsFlag,
		Value:  "impl",
		Usage:  "Solutions impl",
	},
	cli.BoolFlag{
		EnvVar: "CH_SOLUTIONS_DEBUG",
		Name:   debugFlag,
		Usage:  "Start the server in Debug mode",
	},
	cli.BoolFlag{
		EnvVar: "CH_SOLUTIONS_TEXTLOG",
		Name:   textlogFlag,
		Usage:  "Display output log in text format",
	},
	cli.StringFlag{
		EnvVar: "CH_SOLUTIONS_DB",
		Name:   dbFlag,
		Value:  "postgres",
		Usage:  "DB for project",
	},
	cli.StringFlag{
		EnvVar: "CH_SOLUTIONS_DB_URL",
		Name:   dbURLFlag,
		Usage:  "DB URL",
	},
	cli.StringFlag{
		EnvVar: "CH_SOLUTIONS_MIGRATIONS_PATH",
		Name:   dbMigrationsFlag,
		Value:  "../../pkg/migrations/",
		Usage:  "Location of DB migrations",
	},
	cli.StringFlag{
		EnvVar: "CH_SOLUTIONS_KUBE_API_URL",
		Name:   kubeURLFlag,
		Value:  "http://kube-api:1214",
		Usage:  "Kube-API service URL",
	},
	cli.StringFlag{
		EnvVar: "CH_SOLUTIONS_RESOURCE_URL",
		Name:   resourceURLFlag,
		Value:  "http://resource-service:1213",
		Usage:  "Resource service URL",
	},
	cli.BoolFlag{
		EnvVar: "CH_SOLUTIONS_CORS",
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
		return postgres.DBConnect(c.String(dbURLFlag), c.String(dbMigrationsFlag))
	default:
		return nil, errors.New("invalid db")
	}
}
