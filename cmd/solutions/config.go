package main

import (
	"errors"

	"git.containerum.net/ch/solutions/pkg/models"
	"git.containerum.net/ch/solutions/pkg/models/postgres"
	"git.containerum.net/ch/solutions/pkg/server"
	"git.containerum.net/ch/solutions/pkg/server/impl"
	"github.com/urfave/cli"
)

const (
	solutionsFlag    = "solutions"
	debugFlag        = "debug"
	textlogFlag      = "textlog"
	dbFlag           = "db"
	dbURLFlag        = "db_url"
	dbMigrationsFlag = "db_migrations"
	csvURLFlag       = "csv_url"
	kubeURLFlag      = "kube_url"
	resourceURLFlag  = "resource_url"
	converterURLFlag = "converter_url"
)

var flags = []cli.Flag{
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
		EnvVar: "CH_SOLUTIONS_CSV_URL",
		Name:   csvURLFlag,
		Value:  "https://raw.githubusercontent.com/containerum/solution-list/master/containerum-solutions.csv",
		Usage:  "Solutions list CSV file URL",
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
	cli.StringFlag{
		EnvVar: "CH_SOLUTIONS_CONVERTER_URL",
		Name:   converterURLFlag,
		Value:  "http://model-converter:6543",
		Usage:  "Model converter service URL",
	},
}

func getSolutionsSrv(c *cli.Context, services server.Services) (server.SolutionsService, error) {
	switch c.String(solutionsFlag) {
	case "impl":
		return impl.NewSolutionsImpl(services), nil
	default:
		return nil, errors.New("invalid solutions impl")
	}
}

func getDB(c *cli.Context) (models.DB, error) {
	switch c.String(dbFlag) {
	case "postgres":
		return postgres.DBConnect(c.String(dbURLFlag), c.String(dbMigrationsFlag))
	default:
		return nil, errors.New("invalid db")
	}
}
