package postgres

import (
	"context"

	"fmt"

	"strings"

	stypes "git.containerum.net/ch/solutions/pkg/models"
	jsoniter "github.com/json-iterator/go"
)

func (pgdb *pgDB) SaveAvailableSolutionsList(ctx context.Context, solutions stypes.AvailableSolutionsList) error {
	pgdb.log.Infoln("Saving solutions list")

	solutionsarr := []string{}
	for _, s := range solutions.Solutions {
		images, _ := jsoniter.Marshal(s.Images)

		solutionsarr = append(solutionsarr, fmt.Sprintf("('%v', '%v', '%v', '%v', '%v')", s.Name, s.Limits.CPU, s.Limits.RAM, string(images), s.URL))
	}

	rows, err := pgdb.qLog.QueryxContext(ctx, "DELETE FROM available_solutions; INSERT INTO available_solutions (name, cpu, ram, images, url) "+
		"VALUES "+strings.Join(solutionsarr, ","))
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return rows.Err()
	}
	return err
}

func (pgdb *pgDB) GetAvailableSolutionsList(ctx context.Context) (*stypes.AvailableSolutionsList, error) {
	pgdb.log.Infoln("Get solutions list")
	var ret stypes.AvailableSolutionsList

	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT name, cpu, ram, images, url FROM available_solutions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		solution := stypes.AvailableSolution{Limits: &stypes.Limits{}}
		var images string
		err := rows.Scan(&solution.Name, &solution.Limits.CPU, &solution.Limits.RAM, &images, &solution.URL)
		if err != nil {
			return nil, err
		}
		if err := jsoniter.UnmarshalFromString(images, &solution.Images); err != nil {
			return nil, err
		}

		ret.Solutions = append(ret.Solutions, solution)
	}

	return &ret, rows.Err()
}

func (pgdb *pgDB) GetAvailableSolution(ctx context.Context, name string) (*stypes.AvailableSolution, error) {
	pgdb.log.Infoln("Get solution ", name)
	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT name, cpu, ram, images, url FROM available_solutions WHERE name = $1", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, rows.Err()
	}

	solution := stypes.AvailableSolution{Limits: &stypes.Limits{}}
	var images string
	err = rows.Scan(&solution.Name, &solution.Limits.CPU, &solution.Limits.RAM, &images, &solution.URL)
	if err != nil {
		return nil, err
	}
	if err = jsoniter.UnmarshalFromString(images, &solution.Images); err != nil {
		return nil, err
	}

	return &solution, err
}
