package postgres

import (
	"context"

	"fmt"

	"strings"

	"git.containerum.net/ch/solutions/pkg/sErrors"
	stypes "github.com/containerum/kube-client/pkg/model"
	"github.com/json-iterator/go"
)

func (pgdb *pgDB) SaveAvailableSolutionsList(ctx context.Context, solutions stypes.AvailableSolutionsList) error {
	pgdb.log.Infoln("Saving solutions list")

	solutionsstr := make([]string, 0)
	for _, s := range solutions.Solutions {
		images, _ := jsoniter.Marshal(s.Images)

		solutionsstr = append(solutionsstr, fmt.Sprintf("('%v', '%v', '%v', '%v', '%v', 'false', 'true')", s.Name, s.Limits.CPU, s.Limits.RAM, string(images), s.URL))
	}

	rows, err := pgdb.qLog.QueryxContext(ctx, fmt.Sprintf(
		`DELETE FROM available_solutions WHERE local!='true'; 
				INSERT INTO available_solutions (name, cpu, ram, images, url, local, active) VALUES %s 
				ON CONFLICT DO NOTHING;`, strings.Join(solutionsstr, ",")))
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return rows.Err()
	}
	return err
}

func (pgdb *pgDB) CreateAvailableSolution(ctx context.Context, solution stypes.AvailableSolution) error {
	pgdb.log.Infoln("Saving solution")

	images, _ := jsoniter.Marshal(solution.Images)

	rows, err := pgdb.qLog.QueryxContext(ctx,
		`INSERT INTO available_solutions (name, cpu, ram, images, url, local, active) VALUES ($1, $2, $3, $4, $5, $6, $7) 
				ON CONFLICT DO NOTHING;`, solution.Name, solution.Limits.CPU, solution.Limits.RAM, string(images), solution.URL, "true", "true")
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return rows.Err()
	}
	return err
}

func (pgdb *pgDB) UpdateAvailableSolution(ctx context.Context, solution stypes.AvailableSolution) error {
	pgdb.log.Infoln("Updating solution")

	images, _ := jsoniter.Marshal(solution.Images)

	res, err := pgdb.eLog.ExecContext(ctx,
		`UPDATE available_solutions SET (cpu, ram, images, url, local) = ($2, $3, $4, $5, $6) 
				WHERE name = $1`, solution.Name, solution.Limits.CPU, solution.Limits.RAM, string(images), solution.URL, "true")
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		return sErrors.ErrSolutionNotExist()
	}
	return err
}

func (pgdb *pgDB) ActivateAvailableSolution(ctx context.Context, solution string) error {
	pgdb.log.Infoln("Activating solution")

	res, err := pgdb.eLog.ExecContext(ctx,
		`UPDATE available_solutions SET active = 'true' 
				WHERE name = $1 AND active = 'false' AND local = 'true'`, solution)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		return sErrors.ErrSolutionNotExist()
	}
	return err
}

func (pgdb *pgDB) DeactivateAvailableSolution(ctx context.Context, solution string) error {
	pgdb.log.Infoln("Activating solution")

	res, err := pgdb.eLog.ExecContext(ctx,
		`UPDATE available_solutions SET active = 'false' 
				WHERE name = $1 AND active = 'true' AND local = 'true'`, solution)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		return sErrors.ErrSolutionNotExist()
	}
	return err
}

func (pgdb *pgDB) DeleteAvailableSolution(ctx context.Context, solution string) error {
	pgdb.log.Infoln("Updating solution")

	res, err := pgdb.eLog.ExecContext(ctx,
		`DELETE FROM available_solutions WHERE name = $1 AND local = 'true'`, solution)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		return sErrors.ErrSolutionNotExist()
	}
	return err
}

func (pgdb *pgDB) GetAvailableSolutionsList(ctx context.Context, isAdmin bool) (*stypes.AvailableSolutionsList, error) {
	pgdb.log.Infoln("Get solutions list")
	var ret stypes.AvailableSolutionsList

	query := "SELECT name, cpu, ram, images, url, active FROM available_solutions"

	if !isAdmin {
		query = query + " WHERE active = 'true'"
	}

	rows, err := pgdb.qLog.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		solution := stypes.AvailableSolution{Limits: &stypes.SolutionLimits{}}
		var images string
		err := rows.Scan(&solution.Name, &solution.Limits.CPU, &solution.Limits.RAM, &images, &solution.URL, &solution.Active)
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

	solution := stypes.AvailableSolution{Limits: &stypes.SolutionLimits{}}
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
