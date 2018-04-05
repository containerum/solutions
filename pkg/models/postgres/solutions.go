package postgres

import (
	"context"

	"fmt"

	"strings"

	"encoding/json"

	stypes "git.containerum.net/ch/json-types/solutions"
	jsoniter "github.com/json-iterator/go"
)

func (db *pgDB) SaveAvailableSolutionsList(ctx context.Context, solutions stypes.AvailableSolutionsList) error {
	db.log.Infoln("Saving solutions list")

	solutionsarr := []string{}
	for _, s := range solutions.Solutions {
		images, _ := json.Marshal(s.Images)

		solutionsarr = append(solutionsarr, fmt.Sprintf("('%v', '%v', '%v', '%v', '%v')", s.Name, s.Limits.CPU, s.Limits.RAM, string(images), s.URL))
	}

	rows, err := db.qLog.QueryxContext(ctx, "DELETE FROM available_solutions; INSERT INTO available_solutions (name, cpu, ram, images, url) "+
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

func (db *pgDB) GetAvailableSolutionsList(ctx context.Context) (*stypes.AvailableSolutionsList, error) {
	db.log.Infoln("Get solutions list")
	var ret stypes.AvailableSolutionsList

	rows, err := db.qLog.QueryxContext(ctx, "SELECT name, cpu, ram, images, url FROM available_solutions")
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

func (db *pgDB) GetAvailableSolution(ctx context.Context, name string) (*stypes.AvailableSolution, error) {
	db.log.Infoln("Get solution ", name)
	rows, err := db.qLog.QueryxContext(ctx, "SELECT name, cpu, ram, images, url FROM available_solutions WHERE name = $1", name)
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
	if err := jsoniter.UnmarshalFromString(images, &solution.Images); err != nil {
		return nil, err
	}

	return &solution, err
}
